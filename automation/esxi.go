package automation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type EsxiStack struct {
	*auto.Stack
	config EsxiConfig
}

func InitEsxiStack(ctx context.Context, config EsxiConfig) EsxiStack {
	stackName := config.Node.Name
	workDir := filepath.Join(".", "projects", "esxi")

	s, err := auto.UpsertStackLocalSource(ctx, stackName, workDir)
	if err != nil {
		fmt.Printf("Failed to create or select stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created/Selected stack %q\n", stackName)

	return EsxiStack{&s, config}
}

// Config stack
func (s EsxiStack) Config(ctx context.Context) error {
	cfg := s.config
	osRegion := cfg.DeployProps.Region
	osAuthURL := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", osRegion)
	osProjectDomainName := cfg.DeployProps.Domain
	osProjectName := cfg.DeployProps.Project
	osUserName := "d067954"
	osPassword := os.Getenv("OS_PASSWORD")

	// config openstack
	s.SetConfig(ctx, "openstack:region", auto.ConfigValue{Value: osRegion})
	s.SetConfig(ctx, "openstack:authUrl", auto.ConfigValue{Value: osAuthURL})
	s.SetConfig(ctx, "openstack:projectDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:tenantName", auto.ConfigValue{Value: osProjectName})
	s.SetConfig(ctx, "openstack:userDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:userName", auto.ConfigValue{Value: osUserName})
	s.SetConfig(ctx, "openstack:password", auto.ConfigValue{Value: osPassword, Secret: true})
	s.SetConfig(ctx, "openstack:insecure", auto.ConfigValue{Value: "true"})

	// config instance
	// s.SetConfig(ctx, "imageName", auto.ConfigValue{Value: cfg.DeployProps.ImageName})
	s.SetConfig(ctx, "imageName", auto.ConfigValue{Value: "ubuntu-18.04-amd64-vmware"})
	s.SetConfig(ctx, "flavorName", auto.ConfigValue{Value: cfg.DeployProps.FlavorName})
	s.SetConfig(ctx, "nodeUUID", auto.ConfigValue{Value: cfg.Node.UUID})
	return nil
}

// Output from stack
func (s EsxiStack) Output(res auto.UpResult) error {
	return nil
}
