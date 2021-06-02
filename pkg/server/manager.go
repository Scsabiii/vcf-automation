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
	"io/ioutil"
	"path"
	"sync"

	"github.com/sapcc/vcf-automation/pkg/stack"
	"github.com/spf13/viper"
)

type Manager struct {
	controllers          map[string]*StackController
	ProjectRootDirectory string
	ConfigRootDirectory  string
	sync.Mutex
}

type StackController struct {
	*stack.Controller
	running bool
	updCh   chan bool
	canCh   chan bool
}

func NewManager() *Manager {
	workdir := viper.GetString("work_dir")
	projectdir := viper.GetString("project_root")
	configdir := viper.GetString("config_dir")
	if projectdir == "" {
		projectdir = path.Join(workdir, "projects")
	}
	if configdir == "" {
		configdir = path.Join(workdir, "etc")
	}
	logger.Debugf("config directory: %s", configdir)
	logger.Debugf("project directory: %s", projectdir)
	return &Manager{
		ProjectRootDirectory: projectdir,
		ConfigRootDirectory:  configdir,
		controllers:          make(map[string]*StackController),
	}
}

func (m *Manager) Get(cfgFileName string) (*StackController, bool) {
	m.Lock()
	defer m.Unlock()
	c, ok := m.controllers[cfgFileName]
	return c, ok
}

func (m *Manager) Load(cfgFileName string) (*StackController, error) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.controllers[cfgFileName]; ok {
		return nil, fmt.Errorf("controller already exists")
	}
	cfgFilePath := path.Join(m.ConfigRootDirectory, cfgFileName)
	c, err := stack.NewControllerFromConfigFile(m.ProjectRootDirectory, cfgFilePath)
	if err != nil {
		return nil, err
	}
	sc := &StackController{Controller: c}
	m.controllers[cfgFileName] = sc
	return sc, nil
}

func (m *Manager) Update(cfgFileName string) (*StackController, error) {
	m.Lock()
	defer m.Unlock()
	if sc, ok := m.controllers[cfgFileName]; !ok {
		return nil, fmt.Errorf("controller not exist")
	} else {
		err := sc.ReloadConfig()
		if err != nil {
			return nil, err
		}
		return sc, nil
	}
}

func (m *Manager) ListConfigFiles() (cfgFiles []string, err error) {
	files, err := ioutil.ReadDir(manager.ConfigRootDirectory)
	if err != nil {
		return
	}
	for _, f := range files {
		if !f.IsDir() {
			cfgFiles = append(cfgFiles, f.Name())
		}
	}
	return
}

func (m *Manager) ReloadConfigs() (created, updated, stopped map[string]error, err error) {
	created = make(map[string]error)
	updated = make(map[string]error)
	stopped = make(map[string]error)
	cfgFiles, err := manager.ListConfigFiles()
	if err != nil {
		logger.Errorf("load configs failed: %v", err)
		return
	}
	for _, f := range cfgFiles {
		if _, ok := manager.Get(f); !ok {
			nc, err := manager.Load(f)
			created[f] = err
			if err != nil {
				logger.Errorf("load %s: %v", f, err)
				continue
			}
			logger.Infof("%s loaded", f)
			nc.start()
		} else {
			nc, err := manager.Update(f)
			updated[f] = err
			if err != nil {
				logger.Errorf("update %s: %v", f, err)
				continue
			}
			logger.Infof("%s updated", f)
			nc.triggerUpdateStack()
		}
	}
	newfiles := make(map[string]struct{})
	for _, f := range cfgFiles {
		newfiles[f] = struct{}{}
	}
	for fname, c := range m.controllers {
		if _, ok := newfiles[fname]; !ok {
			stopped[fname] = nil
			c.stop()
			logger.Infof("%s stopped", fname)
			delete(m.controllers, fname)
		}
	}
	return
}

func (c *StackController) start() {
	if c.updCh == nil {
		c.updCh = make(chan bool)
	}
	if c.canCh == nil {
		c.canCh = make(chan bool)
	}
	c.running = true
	go c.Controller.Run(c.updCh, c.canCh)
}

func (c *StackController) stop() {
	c.running = false
	go func() {
		c.canCh <- true
	}()
}

func (c *StackController) triggerUpdateStack() {
	go func() {
		c.updCh <- true
	}()
}
