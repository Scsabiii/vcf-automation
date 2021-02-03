package auto

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optup"
	"gopkg.in/yaml.v2"
)

type Controller struct {
	Config
	WorkDir string
	stack   Stack
}

func NewController(workDir, project, stack string) (c Controller, err error) {
	c = Controller{WorkDir: workDir}
	if project != "esxi" && project != "example" {
		err = fmt.Errorf("project must be one of %q and %q", "esxi", "example")
		return
	}
	if err = c.ReadConfig(project, stack); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.Config = Config{
				Project: DeployType(project),
				Stack:   stack,
			}
			err = c.WriteConfig()
		}
	}
	return
}

// ReadConfig reads stack configuration from ./etc directory
func (c Controller) ReadConfig(project, stack string) error {
	fname := fmt.Sprintf("%s-%s.yaml", project, stack)
	fpath := path.Join(c.WorkDir, "etc", fname)
	yamlBytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlBytes, &c.Config)
}

// WriteConfig writes config file in ./etc directory
func (c Controller) WriteConfig() error {
	fname := fmt.Sprintf("%s-%s.yaml", c.Project, c.Stack)
	fpath := path.Join(c.WorkDir, "etc", fname)
	return c.Config.Write(fpath)
}

func (c Controller) AddNode(n Node) error {
	err := c.Config.AddNode(n)
	if err != nil {
		return err
	}
	return c.WriteConfig()
}

func (c Controller) InitStack(ctx context.Context) (s Stack, err error) {
	switch DeployType(c.Project) {
	case DeployExample:
		projectDir := filepath.Join(c.WorkDir, "projects", "example-go")
		s, err = InitExampleStack(ctx, c.Stack, projectDir)
		if err != nil {
			return
		}
	case DeployEsxi:
		projectDir := filepath.Join(c.WorkDir, "projects", "esxi")
		s, err = InitEsxiStack(ctx, c.Stack, projectDir)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("project %q: %v", c.Project, ErrNotSupported)
	}
	c.stack = s
	return
}

func (c Controller) UpdateStack(ctx context.Context, s Stack) error {
	s.Configure(ctx, c.Config)

	fmt.Println("Successfully set config")
	fmt.Println("Starting refresh")

	_, err := s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		return err
	}

	fmt.Println("Refresh succeeded!")
	fmt.Println("Starting update")

	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// run the update to deploy our fargate web service
	res, err := s.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		return err
	}

	fmt.Println("Update succeeded!")

	printOutputs(res.Outputs)
	return nil
}

func (c Controller) DestoryStack(ctx context.Context, s Stack) error {
	s.Configure(ctx, c.Config)

	fmt.Println("Successfully set config")
	fmt.Println("Starting refresh")

	_, err := s.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		return err
	}

	fmt.Println("Refresh succeeded!")
	fmt.Println("Starting stack destroy")

	// wire up our destroy to stream progress to stdout
	stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)

	// destroy our stack and exit early
	_, err = s.Destroy(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to destroy stack: %v", err)
		return err
	}

	fmt.Println("Stack successfully destroyed")
	return nil
}

func (c Controller) GetState(ctx context.Context, s Stack) error {

	return nil
}
