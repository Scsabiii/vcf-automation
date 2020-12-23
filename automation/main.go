package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	destroy := false

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		usage()
		os.Exit(1)
	}

	if argsWithoutProg[0] == "server" {
		runSever()
		os.Exit(0)
	}

	if argsWithoutProg[0] == "destroy" {
		destroy = true
		argsWithoutProg = argsWithoutProg[1:]
	}

	if len(argsWithoutProg) == 0 {
		err := fmt.Errorf("project name not given")
		fmt.Printf("Error: %v\n", err)
		usage()
		os.Exit(1)
	}

	projectName := argsWithoutProg[0]
	argsWithoutProg = argsWithoutProg[1:]
	switch projectName {
	case "example":
		ctx := context.Background()
		stack := initExampleStack(ctx, "dev")
		runStack(ctx, stack, destroy)

	case "esxi":
		if len(argsWithoutProg) == 0 {
			err := fmt.Errorf("node name not given")
			fmt.Printf("Error: %v\n", err)
			usageEsxi()
			os.Exit(1)
		}
		nodeName := argsWithoutProg[0]
		ctx := context.Background()
		// cfg, err := GetEsxiConfig("./etc/node003-bb096.yaml")
		cfg, err := GetEsxiConfig(fmt.Sprintf("./etc/%s.yaml", nodeName))
		if err != nil {
			fmt.Print(err)
		}
		stack := initEsxiStack(ctx, cfg)
		runStack(ctx, stack, destroy)

	default:
		err := fmt.Errorf("project name %q not known", projectName)
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("./automation server")
	fmt.Println("./automation [destroy] projectName")
	fmt.Printf("'projectName' can be %q or %q\n", "example", "esxi")
}

func usageEsxi() {
	fmt.Println("./automation [destroy] esxi nodeName")
}
