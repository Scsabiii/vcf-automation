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

package auto

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"

	pulumistack "github.com/pulumi/pulumi/pkg/v2/resource/stack"
	"github.com/pulumi/pulumi/sdk/v2/go/common/apitype"
	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	"gopkg.in/yaml.v2"
)

type Resource struct {
	Type     string
	URNName  string
	Name     string
	Instance string
	ID       string
}

func printUpdateSummary(s auto.UpdateSummary) {
	if len(*s.ResourceChanges) > 0 {
		log.Println("DEBUG", "resource changes:")
	}
	for k, v := range *s.ResourceChanges {
		log.Println("DEBUG", "\t", k, v)
	}
}

func printStackOutputs(outputs auto.OutputMap) {
	if len(outputs) > 0 {
		log.Println("DEBUG", "stack outputs:")
	}
	for k, v := range outputs {
		log.Println("DEBUG", "\t", k, v.Value)
	}
}

// prints resources from latest checkpoint file
func printStackResources(stackName string) {
	chkpt, err := readCheckpoint(stackName)
	if err != nil {
		log.Println("ERROR", err)
	}
	log.Println("DEBUG", "stack resources:")
	if chkpt.Latest != nil {
		resources := sortResources(chkpt.Latest.Resources)
		printResources(resources, "root", "DEBUG \t")
	}
}

func readCheckpoint(stackName string) (chk *apitype.CheckpointV3, err error) {
	backendURL := os.Getenv("PULUMI_BACKEND_URL")
	if backendURL == "" {
		return nil, ErrBackendURLNotSet
	}
	parsedURL, err := url.Parse(backendURL)
	if err != nil {
		return nil, err
	}
	fname := fmt.Sprintf("%s.json", stackName)
	fpath := path.Join(parsedURL.Path, ".pulumi", "stacks", fname)
	bytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	chk, err = pulumistack.UnmarshalVersionedCheckpointToLatestCheckpoint(bytes)
	if err != nil {
		return nil, err
	}
	return
}

func sortResources(resources []apitype.ResourceV3) map[string][]Resource {
	var parentURN string
	nodes := make(map[string][]Resource)

	for _, r := range resources {
		if r.Parent == "" {
			parentURN = "root"
		} else {
			parentURN = r.Parent.URNName()
		}
		if c, ok := nodes[parentURN]; ok {
			c = append(c, newResource(r))
			nodes[parentURN] = c
		} else {
			c = []Resource{newResource(r)}
			nodes[parentURN] = c
		}
	}

	return nodes
}

func newResource(r apitype.ResourceV3) Resource {
	if name, ok := r.Outputs["name"]; ok {
		n := fmt.Sprintf("%s", name)
		return Resource{r.Type.String(), r.URN.URNName(), r.URN.Name().String(), n, r.ID.String()}
	} else {
		return Resource{r.Type.String(), r.URN.URNName(), r.URN.Name().String(), "", r.ID.String()}
	}
}

func printResources(res map[string][]Resource, resourceURN, prefix string) {
	for _, r := range res[resourceURN] {
		log.Printf("%s %s[%s]: %s %s\n", prefix, r.Name, r.Type, r.Instance, r.ID)
		printResources(res, r.URNName, prefix+"\t")
	}
}

func unmarshalConfig(pluginIn, pluginOut interface{}) error {
	b, err := yaml.Marshal(pluginIn)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, pluginOut)
}
