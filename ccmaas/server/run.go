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
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var (
	newStackCh   = make(chan *CachedController)
	addNodeCh    = make(chan *NodeConfig)
	addStorageCh = make(chan *StorageConfig)

	mu        sync.Mutex
	workdir   string
	etcdir    string
	prjdir    string
	contCache ControllerCache
)

type ControllerCache map[string]*CachedController

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

type CachedController struct {
	*auto.Controller
	err error
}

func Run(port int) {
	workdir = viper.GetString("workdir")
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
	}

	// start refrsh/deploy loop
	go func() {
		ticker := time.Tick(5 * time.Minute)
		log.Println("start deploy loop")
		for {

			// first refresh then update stack
			for _, l := range contCache {
				log.Println("INFO ===", "refresh and update stack", l.Project, l.Stack)
				log.Println("INFO", "refreshing")
				err = l.Refresh(ctx)
				if err != nil {
					log.Println("ERROR", err)
					continue
				}
				log.Println("INFO", "updating")
				err = l.Update(ctx)
				if err != nil {
					log.Println("ERROR", err)
					continue
				}
				log.Println("INFO", "resources:")
				l.PrintStackResources()
			}

			select {
			case l := <-newStackCh:
				log.Println("INFO ===", "new stack created", l.Project, l.Stack)
				log.Println("INFO", "updating")
				err = l.Update(ctx)
				if err != nil {
					l.err = err
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
	r.HandleFunc("/new", newStack).Methods("POST")
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

	// decode and validate payload as auto.Config
	c := auto.Config{}
	err := decoder.Decode(&c)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	if c.Project == "" {
		err = fmt.Errorf("project not set")
		handleError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if c.Project != auto.DeployEsxi && c.Project != auto.DeployExample {
		err = fmt.Errorf("project %q not supported", c.Project)
		handleError(w, http.StatusUnprocessableEntity, err)
		return
	}
	if c.Stack == "" {
		err = fmt.Errorf("stack not set")
		handleError(w, http.StatusUnprocessableEntity, err)
		return
	}

	// TODO: timeout
	l, err := InitController(&c)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	// send stack controller to process loop
	newStackCh <- l

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
	if c, ok := contCache[key]; ok {
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
	log.Println("ERROR", msg)
	json.NewEncoder(w).Encode(msg)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func InitControllerFromFile(fpath string) (*CachedController, error) {
	mu.Lock()
	defer mu.Unlock()
	c, err := auto.NewControllerFromConfigFile(prjdir, fpath)
	if err != nil {
		return nil, err
	}
	if contCache == nil {
		contCache = make(ControllerCache)
	}
	s := CachedController{c, nil}
	fname := fmt.Sprintf("%s-%s.yaml", c.Project, c.Stack)
	contCache[fname] = &s
	return &s, nil
}

func InitController(cfg *auto.Config) (c *CachedController, err error) {
	mu.Lock()
	defer mu.Unlock()

	// create and initialize controller
	ctx := context.Background()
	l := auto.Controller{Config: cfg, ProjectDir: prjdir}
	err = l.InitStack(ctx)
	if err != nil {
		return
	}
	err = l.Configure(ctx)
	if err != nil {
		return
	}

	// write cfg file to disk
	fname := fmt.Sprintf("%s-%s.yaml", l.Project, l.Stack)
	fpath := path.Join(etcdir, fname)
	err = auto.WriteConfig(fpath, cfg, false)
	if err != nil {
		return
	}

	// cache controller
	if contCache == nil {
		contCache = make(ControllerCache)
	}
	c = &CachedController{&l, nil}
	contCache[fname] = c
	return
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
