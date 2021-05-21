package stack

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/sapcc/avocado-automation/pkg/stack/esxi"
	"github.com/sapcc/avocado-automation/pkg/stack/vcf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Controller struct {
	*Config
	ConfigFilePath string
	ProjectPath    string
	stack          Stack
	configured     bool
	err            error
	mu             sync.Mutex
}

// NewController
func NewController(prjPath, cfgPath string, c *Config) (*Controller, error) {
	err := validateConfig(c)
	if err != nil {
		return nil, err
	}
	return &Controller{
		ProjectPath:    prjPath,
		ConfigFilePath: path.Join(cfgPath, c.FileName()),
		Config:         c,
	}, nil
}

// NewControllerFromConfigFile reads configuration file (fname) in the
// configuration directory (cpath), and creates controller from it.
func NewControllerFromConfigFile(projectDir, configFile string) (*Controller, error) {
	c, err := ReadConfig(configFile)
	if err != nil {
		return nil, err
	}
	err = validateConfig(c)
	if err != nil {
		return nil, err
	}
	l := Controller{
		ProjectPath:    projectDir,
		ConfigFilePath: configFile,
		Config:         c,
	}
	return &l, nil
}

func (c *Controller) Run(updateCh <-chan bool) {
	logger := log.WithFields(log.Fields{
		"project": c.Project,
		"stack":   c.Stack,
	})
	tickerDuration := 15 * time.Minute
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		func() {
			ctx := context.Background()
			if c.stack == nil {
				logger.Info("initialize stack")
				if err := c.InitStack(ctx); err != nil {
					c.err = err
					logger.WithError(c.err).Error("initialize stack failed")
					return
				}
			}
			if !c.configured {
				logger.Info("configure stack")
				if err := c.ConfigureStack(ctx); err != nil {
					c.err = err
					logger.WithError(c.err).Error("configure stack failed")
					return
				}
				c.configured = true
			}
			logger.Info("refresh stack")
			if err := c.RefreshStack(ctx); err != nil {
				c.err = err
				logger.WithError(c.err).Error("refresh stack failed")
				return
			}
			logger.Info("update stack")
			if err := c.UpdateStack(ctx); err != nil {
				c.err = err
				logger.WithError(c.err).Error("update stack failed")
				return
			}
			c.err = nil
		}()

		if c.err == nil {
			logger.Info("stack resources:")
			c.PrintStackResources()
		}

		select {
		case <-updateCh:
			// force re-configuring stack since configuration might have
			// changed; reset timer so that next update will wait full
			// tickerDuration
			c.configured = false
			ticker.Reset(tickerDuration)
		case <-ticker.C:
		}
	}
}

// UpdateConfig updates Props.StackProps field of the controller's Config with
// the given Config s. The updated Config is written to the configuration file
// on disk.
func (l *Controller) UpdateConfig(s *Config) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Config.Project != s.Project {
		return fmt.Errorf("unmatched project")
	}
	if l.Config.Stack != s.Stack {
		return fmt.Errorf("unmatched stack")
	}
	nc, err := MergeStackPropsToConfig(l.Config, s.Props.StackProps)
	if err != nil {
		return err
	}
	err = WriteConfig(l.ConfigFilePath, nc)
	if err != nil {
		return err
	}
	l.Config = nc
	return nil
}

// RuntimeError returns error thrown when refresh/update/destroy stack
func (c *Controller) RuntimeError() error {
	return c.stack.GetError()
}

func (l *Controller) InitStack(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch ProjectType(l.Project) {
	case DeployExample:
		projectDir := filepath.Join(l.ProjectPath, "example-go")
		if s, err := InitExampleStack(ctx, l.Stack, projectDir); err != nil {
			return err
		} else {
			l.stack = s
		}
	case DeployEsxi:
		projectDir := filepath.Join(l.ProjectPath, "esxi")
		stackName := l.Stack
		s, err := esxi.InitEsxiStack(ctx, stackName, projectDir)
		if err != nil {
			return err
		}
		l.stack = s
	case DeployVCF:
		projectDir := filepath.Join(l.ProjectPath, "management")
		s, err := vcf.InitVCFStack(ctx, l.Stack, projectDir)
		if err != nil {
			return err
		}
		l.stack = s

	default:
		return fmt.Errorf("project %q: %v", l.Project, ErrNotSupported)
	}

	return nil
}

// ConfigureStack configure the stack with openstack properties (user, domain,
// project and etc.), ssh key pair and stack specific properties.
//
// NOTE All files in <ConfigFilePath> are persistent. Therefore, the SSH key
// pair files are saved in directory <ConfigFilePath>/.ssh.
func (c *Controller) ConfigureStack(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stack == nil {
		return fmt.Errorf("stack uninitialized")
	}
	err := configureOpenstackProps(ctx, c.stack, c.Config.Props.OpenstackProps)
	if err != nil {
		return err
	}
	err = c.readKeypair(path.Join(path.Dir(c.ConfigFilePath), ".ssh"))
	if err != nil {
		return err
	}
	err = configureKeypair(ctx, c.stack, c.Config.Props.Keypair)
	if err != nil {
		return err
	}
	err = configureStackProps(ctx, c.stack, c.Config)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) RefreshStack(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stack == nil {
		return fmt.Errorf("stack uninitialized")
	}
	if err := c.stack.Refresh(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Controller) UpdateStack(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stack == nil {
		return fmt.Errorf("stack uninitialized")
	}
	if res, err := c.stack.Update(ctx); err != nil {
		return err
	} else {
		printStackOutputs(res.Outputs)
	}
	return nil
}

func (c *Controller) GetError() error {
	return c.err
}

// config openstack
func configureOpenstackProps(ctx context.Context, s Stack, p OpenstackProps) error {
	if p.Region == "" {
		return fmt.Errorf("Config.Props.Openstack.Region not set")
	}
	if p.Domain == "" {
		return fmt.Errorf("Config.Props.Openstack.Domain not set")
	}
	if p.Tenant == "" {
		return fmt.Errorf("Config.Props.Openstack.Tenant not set")
	}
	osAuthURL := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", p.Region)
	osUsername := viper.GetString("os_username")
	if osUsername == "" {
		return fmt.Errorf("env variable CCMAAS_OS_USERNAME not configured")
	}
	osPassword := viper.GetString("os_password")
	if osPassword == "" {
		return fmt.Errorf("env variable CCMAAS_OS_PASSWORD not configured")
	}
	c := auto.ConfigMap{
		"openstack:authUrl":           configValue(osAuthURL),
		"openstack:region":            configValue(p.Region),
		"openstack:projectDomainName": configValue(p.Domain),
		"openstack:tenantName":        configValue(p.Tenant),
		"openstack:userDomainName":    configValue(p.Domain),
		"openstack:userName":          configValue(osUsername),
		"openstack:insecure":          configValue("true"),
		"openstack:password":          configSecret(osPassword),
	}
	return s.SetAllConfig(ctx, c)
}

// config key pair
func configureKeypair(ctx context.Context, s Stack, kp Keypair) error {
	if kp.publicKey == "" || kp.privateKey == "" {
		return ErrKeypairNotSet
	}
	err := s.SetConfig(ctx, "publicKey", configValue(kp.publicKey))
	if err != nil {
		return err
	}
	err = s.SetConfig(ctx, "privateKey", configSecret(kp.privateKey))
	if err != nil {
		return err
	}
	return nil
}

// configure stack props
func configureStackProps(ctx context.Context, s Stack, cfg *Config) error {
	switch ProjectType(cfg.Project) {
	case DeployExample:
	case DeployEsxi:
		stackProps := append([]StackProps{cfg.Props.StackProps}, cfg.Props.MoreStackProps...)
		props := make([]esxi.StackProps, len(stackProps))
		err := unmarshalStackProps(cfg.Props.StackProps, &props)
		if err != nil {
			return err
		}
		err = s.(*esxi.Stack).Configure(ctx, props...)
		if err != nil {
			return err
		}
	case DeployVCF:
		stackProps := append([]StackProps{cfg.Props.StackProps}, cfg.Props.MoreStackProps...)
		props := make([]vcf.StackProps, len(stackProps))
		err := unmarshalStackPropList(stackProps, &props)
		if err != nil {
			return err
		}
		err = s.(*vcf.Stack).Configure(ctx, props...)
		if err != nil {
			return err
		}
	}
	return nil
}
