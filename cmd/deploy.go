package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"ccmaas/auto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	destroy bool
	outputs bool
	yaml    bool
	stack   auto.Stack
)

var deployCmd = &cobra.Command{
	Use:   "deploy [cfgFileName]",
	Short: "Provision CCI using config file (without .yaml suffix)",
	Long:  `Provision CCI project (installing ESXi nodes)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deploy(args[0])
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVarP(&destroy, "destory", "d", false, "Destory stack")
	deployCmd.Flags().BoolVarP(&outputs, "outputs", "o", false, "Outputs of stack")
	deployCmd.Flags().BoolVarP(&yaml, "yaml", "y", false, "Yaml output")
}

func deploy(cfgName string) {
	ctx := context.Background()
	cfg := readConfig(cfgName)

	// Initial pulumi stack
	switch cfg.Type {
	case auto.DeployExample:
		cfg.Stack = "dev"
		stack = auto.InitExampleStack(ctx, cfg)
	case auto.DeployEsxi:
		cfg.Stack = cfgName
		stack = auto.InitEsxiStack(ctx, cfg)
	}

	if outputs {
		fmt.Println()
		fmt.Println("Outputs")
		fmt.Println("-------")
		auto.PrintOutputs(ctx, stack)
	} else if yaml {
		yamlOutput := auto.GenYaml(ctx, cfg, stack)
		fmt.Println()
		fmt.Println("Yaml Outputs")
		fmt.Println("-------")
		fmt.Println(string(yamlOutput))
	} else {
		auto.RunStack(ctx, stack, destroy)
	}
}

func readConfig(cfgName string) auto.Config {
	cfgPath := path.Join("etc", fmt.Sprintf("%s.yaml", cfgName))
	cfg, err := auto.ReadConfig(cfgPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if cfg.Type == "" {
		cfg.Type = auto.DeployEsxi
	}
	cfg.Props.Password = viper.GetString("os_password")
	return cfg
}
