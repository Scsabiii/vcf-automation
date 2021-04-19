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

package controller

import (
	"context"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type Stack interface {
	Workspace() auto.Workspace
	Configure(context.Context, *Config) error
	GenYaml(context.Context, *Config) ([]byte, error)

	Refresh(context.Context) error
	Update(context.Context) (auto.UpResult, error)
	// Destroy(context.Context) error
	GetState() interface{}
	GetError() error

	// GetState()

	// Outputs(context.Context) (auto.OutputMap, error)
	// Destroy(context.Context, ...optdestroy.Option) (auto.DestroyResult, error)
	// Preview(ctx context.Context, opts ...optpreview.Option) (auto.PreviewResult, error)
	// Info(ctx context.Context) (auto.StackSummary, error)

	// History(ctx context.Context) ([]auto.UpdateSummary, error)
	// History(ctx context.Context, pageSize int, page int) ([]auto.UpdateSummary, error)
}

// type YamlOutput struct {
// 	Nodes []NodeOutput `yaml:"nodes"`
// }

// type NodeOutput struct {
// 	ID string `yaml:"id"`
// 	IP string `yaml:"ip"`
// }
