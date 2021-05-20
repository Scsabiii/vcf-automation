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
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sapcc/avocado-automation/pkg/stack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	manager *Manager
)

func Run(port int) {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	manager = NewManager()
	log.Infof("configuration directory %s", manager.ConfigDirectory)
	log.Infof("project directory: %s", manager.ProjectDirectory)

	// read configuration files and initialize controllers
	files, err := ioutil.ReadDir(manager.ConfigDirectory)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		c, err := stack.ReadConfig(path.Join(manager.ConfigDirectory, f.Name()))
		if err != nil {
			log.Error(err)
			continue
		}
		_, err = manager.RegisterConfig(c)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("register %s done", f.Name())
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Headers("Content-Type", "application/json")
	r.HandleFunc("/new", newStackHandler).Methods("POST")
	r.HandleFunc("/update", updateStackHandler).Methods("POST")
	r.HandleFunc("/{project}/{stack}/state", getStackState).Methods("GET")
	r.HandleFunc("/{project}/{stack}/update", updateStack).Methods("POST")

	// sd := "/static/"
	// sdpath := viper.GetString("static_dir")
	// r.PathPrefix(sd).Handler(http.StripPrefix(sd, http.FileServer(http.Dir(sdpath))))

	h := pageHandler{
		staticPath:   viper.GetString("static_path"),
		templatePath: viper.GetString("template_path"),
	}
	r.Handle("/", h)

	s := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("localhost:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("listening on port %d", port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err)
		}
	}()

	<-stop
	s.Shutdown(ctx)
}
