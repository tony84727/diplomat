package internal

import (
	"fmt"
	"github.com/tony84727/diplomat/pkg/data"
	"strconv"
	"strings"
)

type ConfigurationUpdater struct {
	Config data.Configuration
}

func (m *ConfigurationUpdater) Set(key string, value string) error {
	navigator := NewConfigNavigator(m.Config)
	paths := strings.Split(key,".")
	out,err := navigator.Get(paths[:len(paths)-1]...)
	if err != nil {
		return err
	}
	switch v := out.(type) {
	case map[interface{}]interface{}:
		v[paths[len(paths) -1]] = value
	case map[string]string:
		v[paths[len(paths) -1]] = value
	case []string:
		i,err := strconv.Atoi(paths[len(paths)])
		if err != nil {
			return err
		}
		v[i] = value
	default:
		return fmt.Errorf("don't know how to handle %T", v)
	}
	return nil
}

func NewConfigurationUpdater(origin data.Configuration) *ConfigurationUpdater {
	return &ConfigurationUpdater{origin}
}

