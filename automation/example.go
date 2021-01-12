package automation

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type ExampleStack struct {
	*auto.Stack
}

func InitExampleStack(ctx context.Context, stackName string) ExampleStack {
	workDir := filepath.Join("projects", "example-go")

	s, err := auto.UpsertStackLocalSource(ctx, stackName, workDir)
	if err != nil {
		fmt.Printf("Failed to create or select stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created/Selected stack %q\n", stackName)
	return ExampleStack{&s}
}

// Config set stack configuration
func (s ExampleStack) Config(ctx context.Context) error {
	osRegion := "qa-de-1"
	osAuthURL := "https://identity-3.qa-de-1.cloud.sap/v3"
	osProjectDomainName := "monsoon3"
	osProjectName := "d067954"
	osUserName := "d067954"
	osPassword := os.Getenv("OS_PASSWORD")
	s.SetConfig(ctx, "openstack:region", auto.ConfigValue{Value: osRegion})
	s.SetConfig(ctx, "openstack:authUrl", auto.ConfigValue{Value: osAuthURL})
	s.SetConfig(ctx, "openstack:projectDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:tenantName", auto.ConfigValue{Value: osProjectName})
	s.SetConfig(ctx, "openstack:userDomainName", auto.ConfigValue{Value: osProjectDomainName})
	s.SetConfig(ctx, "openstack:userName", auto.ConfigValue{Value: osUserName})
	s.SetConfig(ctx, "openstack:password", auto.ConfigValue{Value: osPassword, Secret: true})
	s.SetConfig(ctx, "openstack:insecure", auto.ConfigValue{Value: "true"})
	return nil
}

//
func (e ExampleStack) Output(res auto.UpResult) error {
	// get the instance IP from the stack outputs
	instanceIP, ok := res.Outputs["instanceIP"].Value.(string)
	if !ok {
		fmt.Println("Failed to unmarshall output URL")
		os.Exit(1)
	}

	fmt.Println("Output:")
	fmt.Printf("InstanceIP: %s\n", instanceIP)
	return nil
}
