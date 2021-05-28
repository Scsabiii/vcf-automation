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
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type reloadResponse struct {
	Message string          `json:"message,omitempty"`
	Err     string          `json:"err,omitempty"`
	Created []ConfigMessage `json:"created,omitempty"`
	Updated []ConfigMessage `json:"updated,omitempty"`
	Stopped []ConfigMessage `json:"stopped,omitempty"`
}

type ConfigMessage struct {
	ConfigName string `json:"config_name,omitempty"`
	Message    string `json:"message,omitempty"`
	Err        string `json:"err,omitempty"`
}

func reload(w http.ResponseWriter, r *http.Request) {
	created, updated, stopped, err := manager.ReloadConfigs()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	resp := reloadResponse{
		Message: "configs reloaded",
		Created: make([]ConfigMessage, 0),
		Updated: make([]ConfigMessage, 0),
		Stopped: make([]ConfigMessage, 0),
	}
	for fn, e := range created {
		if e == nil {
			resp.Created = append(resp.Created, ConfigMessage{
				ConfigName: fn,
				Message:    "succsess",
			})
		} else {
			resp.Created = append(resp.Created, ConfigMessage{
				ConfigName: fn,
				Err:        e.Error(),
			})
		}
	}
	for fn, e := range updated {
		if e == nil {
			resp.Updated = append(resp.Updated, ConfigMessage{
				ConfigName: fn,
				Message:    "succsess",
			})
		} else {
			resp.Updated = append(resp.Updated, ConfigMessage{
				ConfigName: fn,
				Err:        e.Error(),
			})
		}
	}
	for fn, e := range stopped {
		if e == nil {
			resp.Stopped = append(resp.Stopped, ConfigMessage{
				ConfigName: fn,
				Message:    "succsess",
			})
		} else {
			resp.Stopped = append(resp.Stopped, ConfigMessage{
				ConfigName: fn,
				Err:        e.Error(),
			})
		}
	}
	js, err := json.Marshal(resp)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func startStack(w http.ResponseWriter, r *http.Request) {
	c, err := getControllerByHttpRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	} else {
		c.start()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("stack %s-%s started\n", c.Project, c.Stack)))
	}
}

func stopStack(w http.ResponseWriter, r *http.Request) {
	c, err := getControllerByHttpRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	c.stop()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("stack %s-%s stopped\n", c.Project, c.Stack)))
}

func reloadStack(w http.ResponseWriter, r *http.Request) {
	c, err := getControllerByHttpRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	nc, err := manager.Update(c.ConfigName())
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	nc.triggerUpdateStack()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("stack %s-%s reloaded\n", c.Project, c.Stack)))
}

func stacks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	for k, c := range manager.controllers {
		if c.running {
			w.Write([]byte(fmt.Sprintf("%s: running\n", k)))
		} else {
			w.Write([]byte(fmt.Sprintf("%s: stopped\n", k)))
		}
	}
}

func getStackState(w http.ResponseWriter, r *http.Request) {
	c, err := getControllerByHttpRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	if stackErr := c.GetError(); stackErr != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(stackErr.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no error in stack deployment"))
	}
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

func getControllerByHttpRequest(r *http.Request) (*StackController, error) {
	vars := mux.Vars(r)
	project := vars["project"]
	stack := vars["stack"]
	fname := fmt.Sprintf("%s-%s.yaml", project, stack)
	if sc, ok := manager.Get(fname); ok {
		return sc, nil
	} else {
		err := fmt.Errorf("config not found: %s", fname)
		return nil, err
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	td := viper.GetString("templates_dir")
	// lp := filepath.Join(td, "layout.html")
	fp := filepath.Join(td, filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			http.Error(w, "not found not found", http.StatusNotFound)
			// http.Redirect(w, r, r.URL.Path, http.StatusNotFound)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

}

type pageHandler struct {
	staticPath   string
	templatePath string
}

func (h pageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header()
	w.Write([]byte("ok"))
}
