package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type EsxiConfig struct {
	Rack        Rack        `yaml:"rack"`
	Node        Node        `yaml:"node"`
	DeployProps DeployProps `yaml:"deployment"`
}

type Rack struct {
	Name string `yaml:"name"`
}

type Node struct {
	Name string `yaml:"name"`
	UUID string `yaml:"uuid"`
}

type DeployProps struct {
	Region     string `yaml:"region"`
	Domain     string `yaml:"domain"`
	Project    string `yaml:"project"`
	ImageName  string `yaml:"image"`
	FlavorName string `yaml:"flavor"`
}

func GetEsxiConfig(filePath string) (cfg EsxiConfig, err error) {
	yamlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlBytes, &cfg)
	if err != nil {
		return
	}
	return
}
