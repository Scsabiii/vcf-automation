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
	"fmt"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
)

type ManagementStack struct {
	auto.Stack
	state ManagementStackState
}

type ManagementStackState struct {
	err error
}

func InitManagementStack(ctx context.Context, stackName, projectDir string) (*ManagementStack, error) {
	s, err := auto.UpsertStackLocalSource(ctx, stackName, projectDir)
	if err != nil {
		return nil, fmt.Errorf("Failed to create or select stack: %v\n", err)
	}
	return &ManagementStack{s, ManagementStackState{}}, nil
}

func (s *ManagementStack) Configure(ctx context.Context, cfg *Config) error {
	return configureOpenstack(ctx, s.Stack, cfg)
}

func (s *ManagementStack) GenYaml(ctx context.Context, cfg *Config) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (s *ManagementStack) Refresh(ctx context.Context) error {
	_, err := s.Stack.Refresh(ctx)
	if err != nil {
		s.state.err = err
		return err
	}
	return nil
}

func (s *ManagementStack) Update(ctx context.Context) (auto.UpResult, error) {
	res, err := s.Stack.Up(ctx)
	if err != nil {
		s.state.err = err
		return auto.UpResult{}, err
	}
	return res, nil
}

func (s *ManagementStack) GetState() interface{} {
	return nil
}

func (s *ManagementStack) GetError() error {
	return nil
}
