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

package stack

import (
	"context"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

type Stack interface {
	Workspace() auto.Workspace
	Refresh(context.Context) error
	Update(context.Context) (auto.UpResult, error)
	SetConfig(context.Context, string, auto.ConfigValue) error
	SetAllConfig(context.Context, auto.ConfigMap) error
	Outputs(context.Context) (auto.OutputMap, error)

	GetState() interface{}
	GetError() error
	GetOutput(context.Context, string) (string, error)
}
