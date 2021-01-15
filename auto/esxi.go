package auto

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type EsxiStack struct {
	*auto.Stack
	config Config
}

func InitEsxiStack(ctx context.Context, cfg Config) EsxiStack {
	workDir := filepath.Join(".", "projects", "esxi")

	s, err := auto.UpsertStackLocalSource(ctx, cfg.Stack, workDir)
	if err != nil {
		fmt.Printf("Failed to create or select stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created/Selected stack %q\n", cfg.Stack)

	return EsxiStack{&s, cfg}
}

// Config stack
func (s EsxiStack) Config(ctx context.Context) error {
	// get node from config
	if len(s.config.Nodes) != 1 {
		return fmt.Errorf("Only one node is allowed, got %d instead", len(s.config.Nodes))
	}
	node := s.config.Nodes[0]

	deployProps := s.config.Props
	osAuthURL := fmt.Sprintf("https://identity-3.%s.cloud.sap/v3", deployProps.Region)
	osProjectDomainName := deployProps.Domain
	osProjectName := deployProps.Project
	osUserName := deployProps.UserName
	osPassword := deployProps.Password

	// config openstack
	s.SetConfig(ctx, "openstack:region", auto.ConfigValue{Value: deployProps.Region})
	s.SetConfig(ctx, "openstack:authUrl", auto.ConfigValue{Value: osAuthURL})
	s.SetConfig(ctx, "openstack:projectDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:tenantName", auto.ConfigValue{Value: osProjectName})
	s.SetConfig(ctx, "openstack:userDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:userName", auto.ConfigValue{Value: osUserName})
	s.SetConfig(ctx, "openstack:password", auto.ConfigValue{Value: osPassword, Secret: true})
	s.SetConfig(ctx, "openstack:insecure", auto.ConfigValue{Value: "true"})

	// config instance
	s.SetConfig(ctx, "imageName", auto.ConfigValue{Value: node.ImageName})
	s.SetConfig(ctx, "flavorName", auto.ConfigValue{Value: node.FlavorName})
	s.SetConfig(ctx, "nodeUUID", auto.ConfigValue{Value: node.UUID})
	return nil
}
