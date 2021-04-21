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
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"

	"github.com/sapcc/avocado-automation/pkg/stack"
)

type Manager struct {
	controllers      map[string]*StackController
	ProjectDirectory string
	ConfigDirectory  string
	*sync.Mutex
}

type StackController struct {
	*stack.Controller
	err   error
	errCh chan error
	updCh chan bool
	okCh  chan bool
}

func (m Manager) RegisterConfigFile(configFile os.FileInfo) error {
	configFilePath := path.Join(m.ConfigDirectory, configFile.Name())
	c, err := stack.NewControllerFromConfigFile(m.ProjectDirectory, configFilePath)
	if err != nil {
		return err
	}
	m.register(c)
	return nil
}

func (m Manager) RegisterNewConfig(cfg *stack.Config) (*StackController, error) {
	c, err := stack.NewController(m.ProjectDirectory, m.ConfigDirectory, cfg)
	if err != nil {
		return nil, err
	}
	return m.register(c), nil
}

func (m Manager) register(c *stack.Controller) *StackController {
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

func (m Manager) Get(key string) (*StackController, error) {
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
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	if c.errCh == nil {
		c.errCh = make(chan error, 0)
	}
	if c.updCh == nil {
		c.updCh = make(chan bool, 0)
	}
	if c.okCh == nil {
		c.okCh = make(chan bool, 0)
	}
	c.Controller.Run(c.updCh, c.errCh, c.okCh)
	go func() {
		for {
			select {
			case err := <-c.errCh:
				c.err = err
			case <-c.okCh:
				c.err = nil
			case <-sigterm:
				return
			}
		}
	}()
}

func (c *StackController) TriggerUpdateStack() {
	go func() {
		c.updCh <- true
	}()
}
