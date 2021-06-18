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
	"strings"

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

type StackSummary struct {
	Name       string            `json:"name,omitempty"`
	ConfigFile string            `json:"config_file,omitempty"`
	Status     string            `json:"status,omitempty"`
	HasError   bool              `json:"has_error,omitempty"`
	Outputs    map[string]string `json:"outputs,omitempty"`
	Links      []Link            `json:"links,omitempty"`
}

type Link struct {
	Name        string `json:"name,omitempty"`
	Url         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

func reload(w http.ResponseWriter, r *http.Request) {
	messages := manager.ReloadConfigs()
	err := writeJson(w, messages)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}

func startStack(w http.ResponseWriter, r *http.Request) {
	c, err := getControllerByHttpRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
	} else {
		c.start()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("stack %s-%s started\n", c.ProjectType, c.StackName)))
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
	w.Write([]byte(fmt.Sprintf("stack %s-%s stopped\n", c.ProjectType, c.StackName)))
}

func reloadStack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]
	stack := vars["stack"]
	nc, err := manager.Update(project, stack)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	nc.triggerUpdateStack()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("stack %s-%s reloaded\n", project, stack)))
}

func jsonFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, err := getControllerByHttpRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	s, err := c.Controller.GetOutput(vars["key"])
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = writeJson(w, s)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}

// TODO: show stck summary only for the project
func stackSummaries(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	ss := make([]StackSummary, 0)
	for k, c := range manager.controllers {
		project, stack := c.GetProjectStackName()
		httpBase := "http://localhost:8080"
		uriBase := httpBase + fmt.Sprintf("/%s/%s", project, stack)
		status := "stopped"
		if c.running {
			status = "running"
		}
		hasError := false
		if c.GetError() != nil {
			hasError = true
		}
		links := make([]Link, 0)
		links = append(links, Link{"cloud-builder", uriBase + "/cloud-builder.json", "payload for cloud builder"})
		links = append(links, Link{"state", uriBase + "/state", "resources deployed by automation"})
		links = append(links, Link{"error", uriBase + "/error", ""})
		links = append(links, Link{"start", uriBase + "/start", "restart automation controller loop"})
		links = append(links, Link{"stop", uriBase + "/stop", "pause automation controller"})
		links = append(links, Link{"reload", uriBase + "/reload", "force controller to reload configuration"})
		s := StackSummary{
			Name:       k,
			ConfigFile: c.ConfigPath,
			Status:     status,
			HasError:   hasError,
			Links:      links,
		}
		ss = append(ss, s)
	}
	err := writeJson(w, ss)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}

func getStackError(w http.ResponseWriter, r *http.Request) {
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

func getStackOutputs(w http.ResponseWriter, r *http.Request) {
	c, err := getControllerByHttpRequest(r)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	o, err := c.GetOutputs()
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	err = writeJson(w, o)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
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

func writeJson(w http.ResponseWriter, data interface{}) error {
	var err error
	var b []byte
	if jsontext, ok := data.(string); ok {
		// when data is a json string
		var objmap map[string]*json.RawMessage
		err = json.Unmarshal([]byte(jsontext), &objmap)
		if err != nil {
			return err
		}
		b, err = json.Marshal(objmap)
		if err != nil {
			return err
		}
	} else {
		// when data is struct
		b, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return nil
}

func handleError(w http.ResponseWriter, statusCode int, e error) {
	w.WriteHeader(statusCode)
	msg := fmt.Sprintf("%d - %s", statusCode, e)
	log.WithField("code", statusCode).WithError(e).Error("handling error")
	json.NewEncoder(w).Encode(msg)
}

func handleErrors(w http.ResponseWriter, statusCode int, errs []error) {
	msg := make([]string, 0)
	for _, e := range errs {
		msg = append(msg, e.Error())
	}
	err := fmt.Errorf(strings.Join(msg, "\n"))
	handleError(w, http.StatusInternalServerError, err)
}

func getControllerByHttpRequest(r *http.Request) (*StackController, error) {
	vars := mux.Vars(r)
	project := vars["project"]
	stack := vars["stack"]
	if sc, ok := manager.Get(project, stack); ok {
		return sc, nil
	} else {
		err := fmt.Errorf("controller not exist: %s/%s", project, stack)
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
