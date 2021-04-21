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

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	manager Manager
)

func Run(port int) {
	workdir := viper.GetString("workdir")

	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	manager = Manager{
		ProjectDirectory: path.Join(workdir, "projects"),
		ConfigDirectory:  path.Join(workdir, "etc"),
	}

	// read configuration files and initialize controllers
	log.Infof("read configuration in directory %s", manager.ConfigDirectory)
	if files, err := ioutil.ReadDir(manager.ConfigDirectory); err == nil {
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if err := manager.RegisterConfigFile(f); err == nil {
				log.Infof("register %s done", f.Name())
			} else {
				log.Error(err)
			}
		}
	} else {
		log.Error(err)
		os.Exit(1)
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
