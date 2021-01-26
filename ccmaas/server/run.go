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
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Run(port int) {
	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json")
	r.HandleFunc("/new", newStack).Methods("POST")
	r.HandleFunc("/{stack}/addnode", addNode).Methods("POST")
	r.HandleFunc("/{stack}/addstorage", addNode).Methods("POST")
	r.HandleFunc("/{stack}/state", getStatus).Methods("GET")
	r.Use(loggingMiddleware)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}

func newStack(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	stackCfg := auto.Config{}
	err := decoder.Decode(&stackCfg)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = stackCfg.Save(false)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stackCfg)
	// exec
}

func addNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	n := auto.Node{}
	err := decoder.Decode(&n)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = auto.AddNode(vars["stack"], n)
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

func getStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stack := vars["stack"]
	fmt.Fprintf(w, "Node: %v\n", stack)
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
