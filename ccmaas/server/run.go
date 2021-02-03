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
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var workDir string

func Run(wd string, port int) {
	workDir = wd
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
	project := vars["project"]
	stack := vars["stack"]

	c, err := auto.NewController(workDir, project, stack)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}

	err = decoder.Decode(&c.Config)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	err = c.WriteConfig()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c.Config)
	// exec
}

func addNode(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	vars := mux.Vars(r)
	project := vars["project"]
	stack := vars["stack"]

	c, err := auto.NewController(workDir, project, stack)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}

	n := auto.Node{}
	err = decoder.Decode(&n)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = c.AddNode(n)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(n)
	// exec
}

func addStorage(w http.ResponseWriter, r *http.Request) {
}

func getState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	stack := vars["stack"]

	ctx := context.Background()

	c, err := auto.NewController(workDir, project, stack)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}
	s, err := c.InitStack(ctx)

	c.GetState(ctx, s)
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
