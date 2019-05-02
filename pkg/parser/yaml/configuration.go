package yaml

import (
	"github.com/tony84727/diplomat/pkg/data"
	"gopkg.in/yaml.v2"
)

type ConfigurationParser struct {
	content       []byte
	configuration data.Configuration
}

func (c ConfigurationParser) GetConfiguration() (data.Configuration, error) {
	if c.configuration != nil {
		return c.configuration, nil
	}
	if err := c.parse(); err != nil {
		return nil, err
	}
	return c.configuration, nil
}

func (c *ConfigurationParser) parse() error {
	var configFile configurationFile
	if err := yaml.Unmarshal(c.content, &configFile); err != nil {
		return err
	}
	c.configuration = configFile
	return nil
}

func NewConfigurationParser(content []byte) *ConfigurationParser {
	return &ConfigurationParser{content: content}
}

func Write(configuration data.Configuration) ([]byte, error) {
	if sc, ok := configuration.(data.SimpleConfiguration); ok {
		return yaml.Marshal(sc)
	}
	return yaml.Marshal(configuration)
}
