package auto

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	pkgstack "github.com/pulumi/pulumi/pkg/v2/resource/stack"
	"github.com/pulumi/pulumi/sdk/v2/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	// pkgstack "github.com/pulumi/pulumi/pkg/v2/resource/stack"
)

type Controller struct {
	Config
	workdir string
	stack   Stack
	mu      sync.Mutex
}

func NewController(workdir, project, stack string) (c *Controller, err error) {
	c = &Controller{workdir: workdir, stack: nil}
	if project != "esxi" && project != "example" {
		err = fmt.Errorf("project must be one of %q and %q", "esxi", "example")
		return
	}
	if err = c.ReadConfig(project, stack); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("WARN", "config does not exist")
			fmt.Println("INFO", "create new config")
			c.Config = Config{
				Stack:   stack,
				Project: DeployType(project),
				Props:   DeployProps{Prefix: stack},
			}
			err = c.WriteConfig()
		}
	}
	return
}

// ReadConfig reads stack configuration from ./etc directory
func (c *Controller) ReadConfig(project, stack string) error {
	fname := fmt.Sprintf("%s-%s.yaml", project, stack)
	fpath := path.Join(c.workdir, "etc", fname)
	fmt.Println("INFO", "load config", fpath)
	return c.Config.Read(fpath)
}

// WriteConfig writes config file in ./etc directory
func (c *Controller) WriteConfig() error {
	fname := fmt.Sprintf("%s-%s.yaml", c.Project, c.Stack)
	fpath := path.Join(c.workdir, "etc", fname)
	fmt.Println("INFO", "write config", fpath)
	return c.Config.Write(fpath)
}

func (c *Controller) AddNode(n Node) error {
	err := c.Config.AddNode(n)
	if err != nil {
		return err
	}
	return c.WriteConfig()
}

func (c *Controller) InitStack(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch DeployType(c.Project) {
	case DeployExample:
		projectDir := filepath.Join(c.workdir, "projects", "example-go")
		if s, err := InitExampleStack(ctx, c.Stack, projectDir); err != nil {
			return err
		} else {
			c.stack = s
		}

	case DeployEsxi:
		fmt.Println("INFO", "initializing esxi stack")

		projectDir := filepath.Join(c.workdir, "projects", "esxi")
		stackName := c.Stack

		fmt.Println("INFO", "use project:", projectDir)
		fmt.Println("INFO", "use stack:", stackName)

		s, err := InitEsxiStack(ctx, stackName, projectDir)
		if err != nil {
			return err
		}
		c.stack = s

		fmt.Println("INFO", "start configuring stack")

		if err := s.Configure(ctx, c.Config); err != nil {
			return err
		}

		fmt.Println("INFO", "successfully set stack config")
		fmt.Println("INFO", "esxi stack initialized")

		//ToDo: load stack stae before refresh

		// 		fmt.Println("INFO", "starting refresh")
		// 		if err := c.stack.Refresh(ctx); err != nil {
		// 			fmt.Printf("Failed to refresh stack: %v\n", err)
		// 			return err
		// 		}
		// 		fmt.Println("INFO", "refresh succeeded!")
	default:
		return fmt.Errorf("project %q: %v", c.Project, ErrNotSupported)
	}

	return nil
}

func (c *Controller) Refresh(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.stack.Refresh(ctx); err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		return err
	}
	return nil
}

func (c *Controller) Update(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("INFO", "Starting update")
	if err := c.stack.Update(ctx); err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		return err
	}
	fmt.Println("INFO", "Update succeeded!")
	return nil
}

func (c *Controller) Destory(ctx context.Context) error {
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

func (c *Controller) State() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return json.Marshal(c.stack.State())
}

func readCheckpoint(stack string) (chk *apitype.CheckpointV3, err error) {
	backendURL := os.Getenv("PULUMI_BACKEND_URL")
	if backendURL == "" {
		return nil, ErrBackendURLNotSet
	}
	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}
	fname := fmt.Sprintf("%s.json", stack)
	fpath := path.Join(parsedURL.Path, ".pulumi", "stacks", fname)
	bytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	chk, err = pkgstack.UnmarshalVersionedCheckpointToLatestCheckpoint(bytes)
	if err != nil {
		return nil, err
	}
	return
}

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

func printOutputs(outs auto.OutputMap) {
	var value string
	for key, out := range outs {
		switch v := out.Value.(type) {
		case string:
			value = v
		case int:
			value = strconv.Itoa(v)
		case int64:
			value = fmt.Sprintf("%d", v)
		default:
			value = ""
		}
		fmt.Printf("%30s\t%s\n", key, value)
	}
}
