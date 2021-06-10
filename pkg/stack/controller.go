package stack

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/sapcc/vcf-automation/pkg/stack/esxi"
	"github.com/sapcc/vcf-automation/pkg/stack/vcf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Controller struct {
	*Config
	configFilePath string
	projectPath    string
	projectRoot    string
	stack          Stack
	stackName      string
	configured     bool
	err            error
	mu             sync.Mutex
}

// NewControllerFromConfigFile reads stack config from configFile (full path of
// configuration file), and initialize controller from it.
func NewControllerFromConfigFile(projectRootDirectory, configFile string) (*Controller, error) {
	cfg, err := ReadConfig(configFile)
	if err != nil {
		return nil, err
	}
	projectName, stackName := getProjectStackName(cfg)
	l := Controller{
		projectRoot:    projectRootDirectory,
		projectPath:    path.Join(projectRootDirectory, projectName),
		configFilePath: configFile,
		Config:         cfg,
		stackName:      stackName,
	}
	err = l.Validate()
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func getProjectStackName(cfg *Config) (projectName, stackName string) {
	p := strings.Split(string(cfg.ProjectType), "/")
	projectName = p[0]
	if len(p) > 1 {
		stackName = p[1] + "-" + cfg.Stack
	} else {
		stackName = cfg.Stack
	}
	return
}

func (c *Controller) ReloadConfig() error {
	cfg, err := ReadConfig(c.configFilePath)
	if err != nil {
		return err
	}
	if cfg.ProjectType != c.ProjectType {
		return fmt.Errorf("project does not match")
	}
	if cfg.Stack != c.Stack {
		return fmt.Errorf("config does not match")
	}
	c.Config = cfg
	return nil
}

func (c *Controller) ConfigName() string {
	return path.Base(c.configFilePath)
}

func (c *Controller) ConfigFullName() string {
	return c.configFilePath
}

func (c *Controller) Validate() error {
	switch c.ProjectType {
	case ProjectEsxi, ProjectExample, ProjectVCFWorkload, ProjectVCFManagement:
		if f, err := os.Stat(c.projectPath); err != nil || !f.IsDir() {
			return fmt.Errorf("project directory does not exist: %s", c.projectPath)
		}
	default:
		return fmt.Errorf("project not supported: %s", c.ProjectType)
	}
	return nil
}

func (c *Controller) Run(updateCh <-chan bool, cancelCh <-chan bool) {
	logger := log.WithFields(log.Fields{
		"package": "stack",
		"project": c.ProjectType,
		"stack":   c.stackName,
	})
	tickerDuration := 15 * time.Minute
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

Forloop:
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
		case <-cancelCh:
			c.configured = false
			break Forloop
		case <-ticker.C:
		}
	}
}

// RuntimeError returns error thrown when refresh/update/destroy stack
func (c *Controller) RuntimeError() error {
	return c.stack.GetError()
}

func (l *Controller) InitStack(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch v := ProjectType(l.ProjectType); v {
	case ProjectExample:
		if s, err := InitExampleStack(ctx, l.stackName, l.projectPath); err != nil {
			return err
		} else {
			l.stack = s
		}
	case ProjectEsxi:
		s, err := esxi.InitEsxiStack(ctx, l.stackName, l.projectPath)
		if err != nil {
			return err
		}
		l.stack = s
	case ProjectVCFManagement, ProjectVCFWorkload:
		s, err := vcf.InitVCFStack(ctx, l.stackName, l.projectPath)
		if err != nil {
			return err
		}
		l.stack = s
	default:
		return fmt.Errorf("project %q: %v", l.ProjectType, ErrNotSupported)
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
		return fmt.Errorf("env variable AUTOMATION_OS_USERNAME not configured")
	}
	osPassword := viper.GetString("os_password")
	if osPassword == "" {
		return fmt.Errorf("env variable AUTOMATION_OS_PASSWORD not configured")
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

// configure stack props
func configureStackProps(ctx context.Context, s Stack, cfg *Config) error {
	switch v := ProjectType(cfg.ProjectType); v {
	case ProjectExample:
	case ProjectEsxi:
		stackProps := append([]StackProps{cfg.Props.StackProps}, cfg.Props.BaseStackProps...)
		props := make([]esxi.StackProps, len(stackProps))
		err := unmarshalStackProps(cfg.Props.StackProps, &props)
		if err != nil {
			return err
		}
		err = s.(*esxi.Stack).Configure(ctx, props...)
		if err != nil {
			return err
		}
	case ProjectVCFManagement, ProjectVCFWorkload:
		stackProps := append([]StackProps{cfg.Props.StackProps}, cfg.Props.BaseStackProps...)
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

func (c *Controller) PrintStackResources() {
	c.mu.Lock()
	defer c.mu.Unlock()
	printStackResources(c.stackName)
}
