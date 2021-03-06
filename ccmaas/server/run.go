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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	newStackCh   = make(chan *auto.Config)
	addNodeCh    = make(chan *NodeConfig)
	addStorageCh = make(chan *StorageConfig)

	mu               sync.Mutex
	etcdir           string
	prjdir           string
	stackControllers StackControllers
)

type StackControllers map[string]*Controller

type NodeConfig struct {
	project string
	stack   string
	node    *auto.Node
}

type StorageConfig struct {
	project string
	stack   string
	storage *auto.Share
}

type Controller struct {
	*auto.Controller
	err error
}

func Run(port int) {
	workdir := viper.GetString("workdir")
	etcdir = path.Join(workdir, "etc")
	prjdir = path.Join(workdir, "projects")
	ctx := context.Background()

	// read configuration files and initialize controllers
	log.Println("INFO", "read configuration in directory", etcdir)
	files, err := ioutil.ReadDir(etcdir)
	if err != nil {
		log.Panic(err)
	}
	for _, f := range files {
		fpath := path.Join(etcdir, f.Name())
		log.Println("INFO", "===", fpath)
		log.Println("INFO", "init controller")
		c, err := InitControllerFromFile(fpath)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
		log.Printf("INFO init stack %q in project %q", c.Stack, c.Project)
		err = c.InitStack(ctx)
		if err != nil {
			log.Println("ERROR", err)
			c.err = err
			continue
		}
		log.Println("INFO", "configure stack")
		err = c.Configure(ctx)
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
		os.Exit(0)
	}

	// start refresh/deploy loop
	go func() {
		ticker := time.Tick(1 * time.Minute)
		for {

			// first refresh then update stack
			for _, c := range stackControllers {
				log.Printf("INFO === stack \"%s-%s\"", c.Project, c.Stack)
				log.Println("INFO", "refreshing")
				err = c.Refresh(ctx)
				if err != nil {
					log.Println("ERROR", err)
					continue
				}
				log.Println("INFO", "updating")
				err = c.Update(ctx)
				if err != nil {
					log.Println("ERROR", err)
					continue
				}
				log.Println("INFO", "resources:")
				c.PrintStackResources()
			}

			select {
			case s := <-newStackCh:
				c, err := InitController(s)
				if err != nil {
					log.Println("ERROR", err)
					continue
				}
				ctx := context.Background()
				err = c.InitStack(ctx)
				if err != nil {
					c.err = err
					log.Println("ERROR", err)
					continue
				}
				err = c.Update(ctx)
				if err != nil {
					c.err = err
					log.Println("ERROR", err)
					continue
				}

			case <-addNodeCh:
				// err = c.AddNode(n)
				// if err != nil {
				// 	handleError(w, http.StatusInternalServerError, err)
				// 	return
				// }

			case <-addStorageCh:
			case <-ticker:
			}

		}
	}()

	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json")
	r.HandleFunc("/{project}/{stack}/new", newStack).Methods("POST")
	r.HandleFunc("/{project}/{stack}/addnode", addNode).Methods("POST")
	r.HandleFunc("/{project}/{stack}/addstorage", addStorage).Methods("POST")
	r.HandleFunc("/{project}/{stack}/state", getState).Methods("GET")
	r.Use(loggingMiddleware)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}

func newStack(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	vars := mux.Vars(r)
	c := auto.Config{
		Stack:   vars["stack"],
		Project: auto.DeployType(vars["project"]),
	}
	err := decoder.Decode(&c.Props)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	// TODO: timeout
	// send stack config to process loop
	newStackCh <- &c

	// request is ok;
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func addNode(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	vars := mux.Vars(r)
	n := auto.Node{}
	err := decoder.Decode(&n)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	node := NodeConfig{
		project: vars["project"],
		stack:   vars["stack"],
		node:    &n,
	}
	addNodeCh <- &node
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(n)
}

func addStorage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	vars := mux.Vars(r)
	s := auto.Share{}
	err := decoder.Decode(&s)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	storage := StorageConfig{
		project: vars["project"],
		stack:   vars["stack"],
		storage: &s,
	}
	addStorageCh <- &storage
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

func getState(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO", "handling", r.URL)
	vars := mux.Vars(r)
	project := vars["project"]
	stack := vars["stack"]
	key := fmt.Sprintf("%s-%s", project, stack)
	if c, ok := stackControllers[key]; ok {
		if c.err != nil {
			handleError(w, http.StatusInternalServerError, c.err)
			return
		}
		err := c.Controller.RuntimeError()
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}
	}
}

func handleError(w http.ResponseWriter, statusCode int, e error) {
	w.WriteHeader(statusCode)
	msg := fmt.Sprintf("%d - %s", statusCode, e)
	w.Write([]byte(msg))
	fmt.Println("ERROR", msg)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func InitControllerFromFile(fpath string) (*Controller, error) {
	mu.Lock()
	defer mu.Unlock()
	workdir := viper.GetString("workdir")
	c, err := auto.NewControllerFromConfigFile(workdir, fpath)
	if err != nil {
		return nil, err
	}
	if stackControllers == nil {
		stackControllers = make(StackControllers)
	}
	s := Controller{c, nil}
	key := fmt.Sprintf("%s-%s", c.Project, c.Stack)
	stackControllers[key] = &s
	return &s, nil
}

func InitController(cfg *auto.Config) (*Controller, error) {
	mu.Lock()
	defer mu.Unlock()

	workdir := viper.GetString("workdir")
	project := string(cfg.Project)
	stack := cfg.Stack
	c, err := auto.NewController(workdir, project, stack)
	if err != nil {
		return nil, err
	}
	if stackControllers == nil {
		stackControllers = make(StackControllers)
	}
	s := Controller{c, nil}
	key := fmt.Sprintf("%s-%s", c.Project, c.Stack)
	stackControllers[key] = &s
	return &s, nil
}

func extractProjectStack(fname string) (project, stack string, err error) {
	if strings.HasSuffix(fname, ".yaml") {
		fname = fname[0 : len(fname)-5]
		projectStackNames := strings.Split(fname, "-")
		if len(projectStackNames) == 2 {
			project = projectStackNames[0]
			stack = projectStackNames[1]
			return
		}
	}
	err = fmt.Errorf("not a configuration file: %q", fname)
	return
}
