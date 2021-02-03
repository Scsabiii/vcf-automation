/******************************************************************************
*
*  Copyright 2021 SAP SE
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
*
******************************************************************************/

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
	Workspace() auto.Workspace
	Configure(context.Context, Config) error
	Outputs(context.Context) (auto.OutputMap, error)
	Refresh(context.Context, ...optrefresh.Option) (auto.RefreshResult, error)
	Destroy(context.Context, ...optdestroy.Option) (auto.DestroyResult, error)
	Up(context.Context, ...optup.Option) (auto.UpResult, error)
	GenYaml(context.Context, Config) ([]byte, error)
}

type YamlOutput struct {
	Nodes []NodeOutput `yaml:"nodes"`
}

type NodeOutput struct {
	ID string `yaml:"id"`
	IP string `yaml:"ip"`
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

func PrintStackOutputs(ctx context.Context, stack Stack) {
	outs, err := stack.Outputs(ctx)
	if err != nil {
		fmt.Printf("PrintOutputs: %v\n", err)
		os.Exit(1)
	}
	printOutputs(outs)
}
