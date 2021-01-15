package auto

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DeployType string

const (
	DeployEsxi    DeployType = "esxi"
	DeployVCF     DeployType = "vcf"
	DeployExample DeployType = "example"
)

type Config struct {
	Type  DeployType  `yaml:"type"`
	Props DeployProps `yaml:"props"`
	Nodes []Node      `yaml:"nodes"`
	Stack string
}

type Node struct {
	Name       string `yaml:"name"`
	UUID       string `yaml:"uuid"`
	ImageName  string `yaml:"image"`
	FlavorName string `yaml:"flavor"`
}

type DeployProps struct {
	Region   string `yaml:"region"`
	Domain   string `yaml:"domain"`
	Project  string `yaml:"project"`
	UserName string `yaml:"user"`
	Password string
}

func ReadConfig(path string) (cfg Config, err error) {
	yamlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlBytes, &cfg)
	if err != nil {
		return
	}
	return
}
