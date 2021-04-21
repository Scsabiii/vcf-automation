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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sapcc/avocado-automation/pkg/stack"
	log "github.com/sirupsen/logrus"
)

func newStackHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := getConfigFromRequestBody(r.Body)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	c, err := manager.RegisterNewConfig(cfg)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	// request is ok;
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c.Config)
}

func updateStackHandler(w http.ResponseWriter, r *http.Request) {
	c, err := getConfigFromRequestBody(r.Body)
	if err != nil {
		handleError(w, http.StatusUnprocessableEntity, err)
		return
	}
	l, err := manager.Get(c.FileName())
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}
	err = l.UpdateConfig(c)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	// trigger stack update
	l.TriggerUpdateStack()

	// request is ok;
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(l.Config)
}

func getConfigFromRequestBody(body io.ReadCloser) (*stack.Config, error) {
	decoder := json.NewDecoder(body)
	defer body.Close()
	c := stack.Config{}
	err := decoder.Decode(&c)
	if err != nil {
		return nil, err
	}
	if c.Project == "" {
		return nil, fmt.Errorf("project not set")
	}
	if c.Stack == "" {
		return nil, fmt.Errorf("stack not set")
	}
	return &c, nil
}

func getStackState(w http.ResponseWriter, r *http.Request) {
	c, err := getControllerForRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}
	if c.err != nil {
		w.WriteHeader(http.StatusOK)
		errstr := c.err.Error()
		// errstr = errstr[strings.Index(errstr, "stderr: "):]
		w.Write([]byte(errstr))
	}
	w.WriteHeader(http.StatusOK)
}

func updateStack(w http.ResponseWriter, r *http.Request) {
	c, err := getControllerForRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	}
	c.TriggerUpdateStack()
	w.WriteHeader(http.StatusOK)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Info("handling request")
		next.ServeHTTP(w, r)
	})
}

func handleError(w http.ResponseWriter, statusCode int, e error) {
	w.WriteHeader(statusCode)
	msg := fmt.Sprintf("%d - %s", statusCode, e)
	log.WithField("code", statusCode).WithError(e).Error("handling error")
	json.NewEncoder(w).Encode(msg)
}

func getControllerForRequest(r *http.Request) (*StackController, error) {
	vars := mux.Vars(r)
	project := vars["project"]
	stack := vars["stack"]
	fname := fmt.Sprintf("%s-%s.yaml", project, stack)
	return manager.Get(fname)
}
