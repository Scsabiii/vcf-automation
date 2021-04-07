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
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/gorilla/mux"
	"github.com/sapcc/avocado-automation/pkg/controller"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	workdir string
	cfgpath string
	prjpath string
	manager Manager
	mu      sync.Mutex
)

type Manager map[string]*Controller

type Controller struct {
	*controller.Controller
	err   error
	errCh chan error
	updCh chan bool
}

func Run(port int) {
	workdir = viper.GetString("workdir")
	cfgpath = path.Join(workdir, "etc")
	prjpath = path.Join(workdir, "projects")

	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	// read configuration files and initialize controllers
	log.Infof("read configuration in directory %s", cfgpath)
	files, err := ioutil.ReadDir(cfgpath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// initialize controllers
	log.Println("*** initialize controllers")
	initializeControllers(files)

	log.Println("*** manager: StartRun()")
	for _, c := range manager {
		c.StartRun()
	}

	log.Printf("listening on port %d", port)
	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json")
	r.HandleFunc("/new", newStackHandler).Methods("POST")
	r.HandleFunc("/update", updateStackHandler).Methods("POST")
	r.HandleFunc("/{project}/{stack}/state", getStackState).Methods("GET")
	r.HandleFunc("/{project}/{stack}/update", updateStack).Methods("POST")
	r.Use(loggingMiddleware)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}

func initializeControllers(files []os.FileInfo) {
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		_, err := newControllerFromConfigFile(f.Name())
		if err != nil {
			log.WithError(err).Errorf("%s: controller initialization fails", f.Name())
			continue
		}
		log.Infof("%s: controller initialized", f.Name())
	}
}

func (c *Controller) StartRun() {
	if c.errCh == nil {
		c.errCh = make(chan error, 0)
	}
	if c.updCh == nil {
		c.updCh = make(chan bool, 0)
	}
	go func() {
		for {
			select {
			case err := <-c.errCh:
				c.err = err
			}
		}
	}()
	go c.Controller.Run(c.updCh, c.errCh)
}

func (c *Controller) StartUpdateStack() {
	go func() {
		c.updCh <- true
	}()
}

func getController(fname string) (*Controller, error) {
	mu.Lock()
	defer mu.Unlock()
	if manager == nil {
		return nil, fmt.Errorf("no manager: %s", fname)
	}
	if s, ok := manager[fname]; !ok {
		return nil, fmt.Errorf("no manager: %s", fname)
	} else {
		return s, nil
	}
}

func newControllerFromConfig(cfg *controller.Config) (*Controller, error) {
	mu.Lock()
	defer mu.Unlock()
	l, err := controller.NewController(prjpath, cfgpath, cfg)
	if err != nil {
		return nil, err
	}
	if manager == nil {
		manager = make(Manager)
	}
	c := Controller{l, nil, nil, nil}
	manager[l.FileName()] = &c
	return &c, nil
}

func newControllerFromConfigFile(fname string) (*Controller, error) {
	mu.Lock()
	defer mu.Unlock()
	fpath := path.Join(cfgpath, fname)
	l, err := controller.NewControllerFromConfigFile(prjpath, fpath)
	if err != nil {
		return nil, err
	}
	if manager == nil {
		manager = make(Manager)
	}
	c := Controller{l, nil, nil, nil}
	manager[l.FileName()] = &c
	return &c, nil
}
