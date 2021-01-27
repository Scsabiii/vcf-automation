package auto

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optup"
)

type Controller struct {
	Config
	WorkDir     string
	ProjectName string
	StackName   string
}

func NewController(workDir, project, stack string) (*Controller, error) {
	if project != "esxi" && project != "example" {
		return nil, fmt.Errorf("project must be one of %q and %q", "esxi", "example")
	}
	c := Controller{WorkDir: workDir, ProjectName: project, StackName: stack}
	return &c, nil
}

func (c Controller) AddNode(n Node) error {
	err := c.LoadConfig()
	if err != nil {
		return err
	}
	err = c.Config.AddNode(n)
	if err != nil {
		return err
	}
	return c.SaveConfig()
}

func (c Controller) LoadConfig() error {
	fpath := path.Join(c.WorkDir, fmt.Sprintf("%s-%s.yaml", c.ProjectName, c.StackName))
	return ReadConfig(fpath, &c.Config)
}

// Save new configuration file; overwrite is not allowed
func (c Controller) SaveNewConfig() error {
	fpath := path.Join(c.WorkDir, fmt.Sprintf("%s-%s.yaml", c.ProjectName, c.StackName))
	if fileExists(fpath) {
		return fmt.Errorf("file %q exists", fpath)
	}
	return c.Config.Write(fpath)
}

// Save configuration file; overwrite is allowed
func (c Controller) SaveConfig() error {
	fpath := path.Join(c.WorkDir, fmt.Sprintf("%s-%s.yaml", c.ProjectName, c.StackName))
	return c.Config.Write(fpath)
}

func (c Controller) InitStack(ctx context.Context) (s Stack, err error) {
	switch c.Type {
	case DeployExample:
		projectDir := filepath.Join(c.WorkDir, "projects", "example-go")
		s, err = InitExampleStack(ctx, c.StackName, projectDir)
		if err != nil {
			return
		}
	case DeployEsxi:
		projectDir := filepath.Join(c.WorkDir, "projects", "esxi")
		s, err = InitEsxiStack(ctx, c.StackName, projectDir)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("type not supported: %s", c.Type)
	}
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

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
