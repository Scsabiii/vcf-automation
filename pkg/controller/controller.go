package controller

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"sync"
)

type Controller struct {
	*Config
	ConfigFile  string
	ProjectPath string
	stack       Stack
	mu          sync.Mutex
}

// NewController creates controller with given Config c, and writes the config to disk
func NewController(ppath, cpath string, c *Config) (*Controller, error) {
	l := Controller{
		ProjectPath: ppath,
		ConfigFile:  path.Join(cpath, c.FileName()),
		Config:      c}
	err := writeConfig(l.ConfigFile, c, false)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// NewControllerFromConfigFile reads configuration file (fname) in the
// configuration directory (cpath), and creates controller from it.
func NewControllerFromConfigFile(prjpath, cfgfilepath string) (*Controller, error) {
	c, err := readConfig(cfgfilepath)
	if err != nil {
		return nil, err
	}
	if !isValidProject(c.Project) {
		return nil, fmt.Errorf("project not suported: %q", c.Project)
	}
	l := Controller{
		Config:      c,
		ProjectPath: prjpath,
		ConfigFile:  cfgfilepath,
	}
	return &l, nil
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
	err = writeConfig(l.ConfigFile, nc, true)
	if err != nil {
		return err
	}
	l.Config = nc
	return nil
}

// RuntimeError returns error thrown when refresh/update/destroy stack
func (c *Controller) RuntimeError() error {
	return c.stack.Error()
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

	default:
		return fmt.Errorf("project %q: %v", l.Project, ErrNotSupported)
	}

	return nil
}

func (l *Controller) ConfigureStack(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.stack == nil {
		return fmt.Errorf("stack uninitialized")
	}
	if err := l.stack.Configure(ctx, l.Config); err != nil {
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

func (c *Controller) DestoryStack(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("INFO", "Starting stack destroy")
	if err := c.stack.Destroy(ctx); err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		return err
	}
	fmt.Println("Stack successfully destroyed")
	return nil
}

func (c *Controller) PrintStackResources() {
	c.mu.Lock()
	defer c.mu.Unlock()
	printStackResources(c.Stack)
}

func isValidProject(p ProjectType) bool {
	if p == DeployEsxi {
		return true
	} else if p == DeployExample {
		return true
	}
	return false
}

// func (c *Controller) State() ([]byte, error) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	if c.stack == nil {
// 		return nil, fmt.Errorf("stack uninitialized")
// 	}
// 	return json.Marshal(c.stack.State())
// }

// func (c *Controller) GetState(ctx context.Context) error {
// 	s := c.stack
// 	if s == nil {
// 		return ErrStackNotInitialized
// 	}
// 	return nil
// }

// func (c *Controller) PrintStackOutputs(ctx context.Context) {
// 	outs, err := c.stack.Outputs(ctx)
// 	if err != nil {
// 		fmt.Printf("PrintOutputs: %v\n", err)
// 		os.Exit(1)
// 	}
// 	printOutputs(outs)
// }
// func printOutputs(outs auto.OutputMap) {
// 	var value string
// 	for key, out := range outs {
// 		switch v := out.Value.(type) {
// 		case string:
// 			value = v
// 		case int:
// 			value = strconv.Itoa(v)
// 		case int64:
// 			value = fmt.Sprintf("%d", v)
// 		default:
// 			value = ""
// 		}
// 		fmt.Printf("%30s\t%s\n", key, value)
// 	}
// }
