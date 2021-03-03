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
	"path"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	newStackCh   = make(chan *auto.Config)
	addNodeCh    = make(chan *NodeConfig)
	addStorageCh = make(chan *StorageConfig)
	getStateCh   = make(chan int)
)

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
	Error error
}

func Run(port int) {
	workdir := viper.GetString("workdir")
	stackControllers := make(map[string]*Controller)
	ctx := context.Background()

	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json")
	r.HandleFunc("/{project}/{stack}/new", newStack).Methods("POST")
	r.HandleFunc("/{project}/{stack}/addnode", addNode).Methods("POST")
	r.HandleFunc("/{project}/{stack}/addstorage", addStorage).Methods("POST")
	r.HandleFunc("/{project}/{stack}/state", getState).Methods("GET")
	r.Use(loggingMiddleware)

	// read configuration files and initialize controllers
	// Config files are in yaml format, and is named with pattern
	// <project_name>-<stack-name>.yaml.
	log.Println("INFO", "read configuration files in", workdir)
	files, err := ioutil.ReadDir(path.Join(workdir, "etc"))
	if err != nil {
		log.Panic(err)
	}

	for _, f := range files {
		log.Println("INFO", "processing file", f.Name())

		project, stack, err := extractProjectStack(f.Name())
		if err != nil {
			log.Println("WARN", err)
			continue
		}
		key := fmt.Sprintf("%s-%s", project, stack)
		stackControllers[key] = &Controller{}
		c, err := auto.NewController(workdir, project, stack)
		if err != nil {
			log.Println("ERROR", err)
			stackControllers[key].Error = err
			continue
		}
		err = c.InitStack(ctx)
		if err != nil {
			log.Println("ERROR", err)
			stackControllers[key].Error = err
			continue
		}
		stackControllers[key].Controller = c
	}

	// start refresh/deploy loop
	go func() {
		ticker := time.Tick(5 * time.Minute)
		for {

			// first refresh then update stack
			for _, c := range stackControllers {
				log.Printf("INFO refreshing stack \"%s-%s\"\n", c.Project, c.Stack)
				err := c.Refresh(ctx)
				if err != nil {
					log.Println("ERROR", err)
					continue
				}
				log.Printf("INFO updating stack \"%s-%s\"\n", c.Project, c.Stack)
				err = c.Update(ctx)
				if err != nil {
					log.Println("ERROR", err)
					continue
				}
				c.PrintStackResources()
			}

			select {
			case s := <-newStackCh:
				workdir := viper.GetString("workdir")

				key := fmt.Sprintf("%s-%s", s.Project, s.Stack)
				stackControllers[key] = &Controller{}

				c, err := auto.NewController(workdir, string(s.Project), s.Stack)
				if err != nil {
					stackControllers[key].Error = err
					log.Println("ERROR", err)
					continue
				} else {
					stackControllers[key].Controller = c
				}

				ctx := context.Background()
				err = c.InitStack(ctx)
				if err != nil {
					stackControllers[key].Error = err
					log.Println("ERROR", err)
					continue
				}

				err = c.Update(ctx)
				if err != nil {
					stackControllers[key].Error = err
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
			case <-getStateCh:
			case <-ticker:
			}

		}
	}()

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
	// vars := mux.Vars(r)
	// project := vars["project"]
	// stack := vars["stack"]
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
