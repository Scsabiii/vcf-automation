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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	manager *Manager
	logger  = log.WithField("package", "server")
)

func Run(port int) {
	log.SetFormatter(&log.TextFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyMsg: "message",
		},
	})
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	// load configuration files and initialize controllers
	manager = NewManager()
	manager.ReloadConfigs()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.Headers("Content-Type", "application/json")
	r.HandleFunc("/stacks", stacks).Methods("GET")
	r.HandleFunc("/reload", reload).Methods("GET")
	r.HandleFunc("/{project}/{stack}", getStackOutputs).Methods("GET")
	r.HandleFunc("/{project}/{stack}/error", getStackError).Methods("GET")
	r.HandleFunc("/{project}/{stack}/start", startStack).Methods("GET")
	r.HandleFunc("/{project}/{stack}/stop", stopStack).Methods("GET")
	r.HandleFunc("/{project}/{stack}/reload", reloadStack).Methods("GET")

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
		logger.Printf("listening on port %d", port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err)
		}
	}()

	<-stop
	s.Shutdown(ctx)
}
