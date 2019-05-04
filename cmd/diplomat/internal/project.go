package internal

import (
	"github.com/tony84727/diplomat/pkg/data"
	"github.com/tony84727/diplomat/pkg/parser/yaml"
	"io/ioutil"
	"os"
)

type Project struct {
	data.SourceSet
}

func (p Project) LoadConfig() (data.Configuration, error) {
	configPath, err := p.SourceSet.GetConfigurationFile()
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	parser := yaml.NewConfigurationParser(content)
	return parser.GetConfiguration()
}

func NewProject(projectDir string) *Project {
	return &Project{data.NewFileSystemSourceSet(projectDir)}
}

func FindProject(projectRoot *string) (project *Project, err error) {
	if projectRoot != nil {
		// might need to implement some check here
		return NewProject(*projectRoot), nil
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		return NewProject(pwd), nil
	}
}
