package auto

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optrefresh"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto/optup"
)

type Stack interface {
	Config(context.Context) error
	Outputs(context.Context) (auto.OutputMap, error)
	Refresh(context.Context, ...optrefresh.Option) (auto.RefreshResult, error)
	Destroy(context.Context, ...optdestroy.Option) (auto.DestroyResult, error)
	Up(context.Context, ...optup.Option) (auto.UpResult, error)
}

func RunStack(ctx context.Context, stack Stack, destroy bool) {
	stack.Config(ctx)

	fmt.Println("Successfully set config")
	fmt.Println("Starting refresh")

	_, err := stack.Refresh(ctx)
	if err != nil {
		fmt.Printf("Failed to refresh stack: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Refresh succeeded!")

	if destroy {
		fmt.Println("Starting stack destroy")
		// wire up our destroy to stream progress to stdout
		stdoutStreamer := optdestroy.ProgressStreams(os.Stdout)
		// destroy our stack and exit early
		_, err := stack.Destroy(ctx, stdoutStreamer)
		if err != nil {
			fmt.Printf("Failed to destroy stack: %v", err)
		}
		fmt.Println("Stack successfully destroyed")
		os.Exit(0)
	}

	fmt.Println("Starting update")

	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	// run the update to deploy our fargate web service
	res, err := stack.Up(ctx, stdoutStreamer)
	if err != nil {
		fmt.Printf("Failed to update stack: %v\n\n", err)
		os.Exit(1)
	}

	fmt.Println("Update succeeded!")

	printOutputs(res.Outputs)
}

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

func PrintOutputs(ctx context.Context, stack Stack) {
	outs, err := stack.Outputs(ctx)
	if err != nil {
		fmt.Printf("PrintOutputs: %v\n", err)
	}
	printOutputs(outs)
}
