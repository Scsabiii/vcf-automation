package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	auto "cci-operator/automation"
)

var destroy bool
var project string
var node string

// automationCmd represents the automation command
var automationCmd = &cobra.Command{
	Use:   "automation",
	Short: "cci automation",
	Long:  `cci automation command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("automation called")
		run()
	},
}

func init() {
	rootCmd.AddCommand(automationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// automationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// automationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	automationCmd.Flags().BoolVarP(&destroy, "destory", "d", false, "Destory stack")
	automationCmd.Flags().StringVarP(&project, "project", "p", "esxi", "Project name")
	automationCmd.Flags().StringVarP(&node, "node", "n", "", "Node")
}

func run() {
	var err error
	ctx := context.Background()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch project {
	case "example":
		stack := auto.InitExampleStack(ctx, "dev")
		auto.RunStack(ctx, stack, destroy)
	case "esxi":
		var config auto.EsxiConfig
		if node == "" {
			err = fmt.Errorf("Must specify node name when project is esxi")
			break
		}
		filePath := path.Join(cwd, "etc")
		filePath = fmt.Sprintf("%s/%s.yaml", filePath, node)
		config, err = auto.ReadEsxiConfig(filePath)
		fmt.Println(config)
		if err != nil {
			break
		}
		stack := auto.InitEsxiStack(ctx, config)
		auto.RunStack(ctx, stack, destroy)
	default:
		err = fmt.Errorf("Invalid project: %q", project)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
