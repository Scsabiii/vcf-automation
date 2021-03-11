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
	"ccmaas/auto"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	workdir  string
	cfgpath  string
	prjpath  string
	mgrCache ManagerCache
	stackCh  chan *StackManager
	mu       sync.Mutex
)

type ManagerCache map[string]*StackManager

type StackManager struct {
	*auto.Controller
	err error
}

func Run(port int) {
	stackCh = make(chan *StackManager)
	workdir = viper.GetString("workdir")
	cfgpath = path.Join(workdir, "etc")
	prjpath = path.Join(workdir, "projects")

	// read configuration files and initialize controllers
	log.Println("INFO", "read configuration in directory", cfgpath)
	files, err := ioutil.ReadDir(cfgpath)
	if err != nil {
		log.Panic(err)
	}
	initialize(files)

	// start refrsh/deploy loop
	log.Println("INFO", "go run loop")
	go run()

	// http server
	log.Printf("INFO listening on port %d\n", port)
	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json")
	r.HandleFunc("/new", newStackHandler).Methods("POST")
	r.HandleFunc("/update", updateStackHandler).Methods("POST")
	r.HandleFunc("/{project}/{stack}/state", getState).Methods("GET")
	r.Use(loggingMiddleware)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}

func initialize(files []os.FileInfo) {
	ctx := context.Background()
	for _, f := range files {
		log.Println("INFO", "===", f.Name())
		m, err := newManagerFromConfigFile(f.Name())
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
		log.Printf("INFO init stack %q in project %q", m.Stack, m.Project)
		err = m.InitStack(ctx)
		if err != nil {
			log.Println("ERROR", err)
			m.err = err
			continue
		}
		log.Println("INFO", "configure stack")
		err = m.ConfigureStack(ctx)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
	}
}

func run() {
	log.Println("start deploy loop")
	ctx := context.Background()
	ticker := time.Tick(5 * time.Minute)
	for {
		// Refresh and update all managed stacks
		for _, m := range mgrCache {
			log.Println("INFO ===", "refresh and update stack", m.Project, m.Stack)
			log.Println("INFO", "refreshing")
			err := m.RefreshStack(ctx)
			if err != nil {
				log.Println("ERROR", err)
				continue
			}
			log.Println("INFO", "updating")
			err = m.UpdateStack(ctx)
			if err != nil {
				log.Println("ERROR", err)
				continue
			}
			log.Println("INFO", "resources:")
			m.PrintStackResources()
		}

		select {
		case m := <-stackCh:
			// Configure stack when stack is created or updated
			log.Println("INFO", "configure stack")
			err := m.ConfigureStack(ctx)
			if err != nil {
				m.err = err
				log.Println("ERROR", err)
				return
			}
		case <-ticker:
		}
	}
}

func getManager(fname string) (*StackManager, error) {
	mu.Lock()
	defer mu.Unlock()
	if mgrCache == nil {
		return nil, fmt.Errorf("no manager: %s", fname)
	}
	if s, ok := mgrCache[fname]; !ok {
		return nil, fmt.Errorf("no manager: %s", fname)
	} else {
		return s, nil
	}
}

func newManagerFromConfig(cfg *auto.Config) (*StackManager, error) {
	mu.Lock()
	defer mu.Unlock()
	ctx := context.Background()
	// create and initialize controller
	l, err := auto.NewController(prjpath, cfgpath, cfg)
	if err != nil {
		return nil, err
	}
	// initialize stack
	err = l.InitStack(ctx)
	if err != nil {
		return nil, err
	}
	if mgrCache == nil {
		mgrCache = make(ManagerCache)
	}
	m := StackManager{l, nil}
	mgrCache[l.FileName()] = &m
	return &m, nil
}

func newManagerFromConfigFile(fname string) (*StackManager, error) {
	mu.Lock()
	defer mu.Unlock()
	fpath := path.Join(cfgpath, fname)
	l, err := auto.NewControllerFromConfigFile(prjpath, fpath)
	if err != nil {
		return nil, err
	}
	if mgrCache == nil {
		mgrCache = make(ManagerCache)
	}
	m := StackManager{l, nil}
	mgrCache[l.FileName()] = &m
	return &m, nil
}
