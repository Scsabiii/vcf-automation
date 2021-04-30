package stack

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Controller struct {
	*Config
	ConfigFilePath string
	ProjectPath    string
	configured     bool
	stack          Stack
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
func NewControllerFromConfigFile(prjpath, cfgfilepath string) (*Controller, error) {
	c, err := ReadConfig(cfgfilepath)
	if err != nil {
		return nil, err
	}
	err = validateConfig(c)
	if err != nil {
		return nil, err
	}
	l := Controller{
		ProjectPath:    prjpath,
		ConfigFilePath: cfgfilepath,
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
				c.err = c.InitStack(ctx)
				if c.err != nil {
					logger.WithError(c.err).Error("initialize stack failed")
					return
				}
			}
			if !c.configured {
				logger.Info("configure stack")
				c.err = c.ConfigureStack(ctx)
				if c.err != nil {
					logger.WithError(c.err).Error("configure stack failed")
					return
				}
				c.configured = true
			}
			logger.Info("refresh stack")
			c.err = c.RefreshStack(ctx)
			if c.err != nil {
				logger.WithError(c.err).Error("refresh stack failed")
				return
			}
			logger.Info("update stack")
			c.err = c.UpdateStack(ctx)
			if c.err != nil {
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
		s, err := InitEsxiStack(ctx, stackName, projectDir)
		if err != nil {
			return err
		}
		l.stack = s
	case DeployManagement:
		projectDir := filepath.Join(l.ProjectPath, "management")
		s, err := InitManagementStack(ctx, l.Config.Stack, projectDir)
		if err != nil {
			return err
		}
		l.stack = s

	default:
		return fmt.Errorf("project %q: %v", l.Project, ErrNotSupported)
	}

	return nil
}

// ConfigureStack applies c.Config to the stack's configuration file, by
// calling stack's Configure() function.
//
// Note: The stack configuration file is a yaml file located in the project
// directory, with name Pulumi.{stack_name}.yaml. Controller updates the file
// in each loop to make sure it is always consistent with the controller's
// Config.
//
// Note: SSH key pair files (id_rsa and id_rsa.pub) are in the .ssh/
// subdirectory of the directory where config file locates. If the config file
// path is /foo/bar/config.yaml, the ssh key files are /foo/bar/.ssh/id_rsa and
// /foo/bar/.ssh/id_rsa.pub.  The stack's Configure() function should return
// ErrKeypairNotSet if the key pair is not read yet (see inline comment below).
func (c *Controller) ConfigureStack(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stack == nil {
		return fmt.Errorf("stack uninitialized")
	}
	if err := c.stack.Configure(ctx, c.Config); err != nil {
		// Read the keypair from disk and run stack's Configure() function
		// again, if the keypair is not read yet.
		if errors.Is(err, ErrKeypairNotSet) {
			err = c.readKeypair(path.Join(path.Dir(c.ConfigFilePath), ".ssh"))
			if err != nil {
				return err
			}
			err = c.stack.Configure(ctx, c.Config)
			if err != nil {
				return err
			}
		} else {
			return err
		}
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
