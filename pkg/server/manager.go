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

package server

import (
	"fmt"
	"path"
	"sync"

	"github.com/sapcc/avocado-automation/pkg/stack"
	"github.com/spf13/viper"
)

type Manager struct {
	controllers      map[string]*StackController
	ProjectDirectory string
	ConfigDirectory  string
	sync.Mutex
}

type StackController struct {
	*stack.Controller
	updCh chan bool
}

func NewManager() *Manager {
	workdir := viper.GetString("workdir")
	projectdir := viper.GetString("project_dir")
	if projectdir == "" {
		projectdir = path.Join(workdir, "projects")
	}
	return &Manager{
		ProjectDirectory: projectdir,
		ConfigDirectory:  path.Join(workdir, "etc"),
	}
}

func (m *Manager) RegisterConfig(cfg *stack.Config) (*StackController, error) {
	c, err := stack.NewController(m.ProjectDirectory, m.ConfigDirectory, cfg)
	if err != nil {
		return nil, err
	}
	return m.register(c), nil
}

func (m *Manager) RegisterNewConfig(cfg *stack.Config) (*StackController, error) {
	c, err := stack.NewController(m.ProjectDirectory, m.ConfigDirectory, cfg)
	if err != nil {
		return nil, err
	}
	err = stack.WriteNewConfig(c.ConfigFilePath, c.Config)
	if err != nil {
		return nil, err
	}
	return m.register(c), nil
}

func (m *Manager) register(c *stack.Controller) *StackController {
	m.Lock()
	defer m.Unlock()
	if m.controllers == nil {
		m.controllers = make(map[string]*StackController)
	}
	key := c.Config.FileName()
	sc := StackController{Controller: c}
	m.controllers[key] = &sc
	m.controllers[key].Run()
	return &sc
}

func (m *Manager) Get(key string) (*StackController, error) {
	m.Lock()
	defer m.Unlock()
	errNotFound := fmt.Errorf("%q not found", key)
	if m.controllers == nil {
		return nil, errNotFound
	}
	if c, ok := m.controllers[key]; ok {
		return c, nil
	} else {
		return nil, errNotFound
	}
}

func (c *StackController) Run() {
	if c.updCh == nil {
		c.updCh = make(chan bool, 0)
	}
	go c.Controller.Run(c.updCh)
}

func (c *StackController) TriggerUpdateStack() {
	go func() {
		c.updCh <- true
	}()
}
