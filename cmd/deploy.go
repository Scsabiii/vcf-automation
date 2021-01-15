package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"cci-operator/auto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	destroy bool
	outputs bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy [config]",
	Short: "Provision CCI using config file (without .yaml suffix)",
	Long:  `Provision CCI project (installing ESXi nodes)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var stack auto.Stack
		ctx := context.Background()

		cfgName := args[0]
		cfg := readConfig(cfgName)

		// Initial pulumi stack
		switch cfg.Type {
		case auto.DeployExample:
			cfg.Stack = "dev"
			stack = auto.InitExampleStack(ctx, cfg)
		case auto.DeployEsxi:
			cfg.Stack = cfgName
			stack = auto.InitEsxiStack(ctx, cfg)
		case auto.DeployVCF:
			cfg.Stack = cfgName
			fmt.Println("Not implemented")
			os.Exit(0)
		}

		if outputs {
			auto.PrintOutputs(ctx, stack)
		} else {
			auto.RunStack(ctx, stack, destroy)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVarP(&destroy, "destory", "d", false, "Destory stack")
	deployCmd.Flags().BoolVarP(&outputs, "outputs", "o", false, "Outputs of stack")
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
