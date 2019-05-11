package internal

import "github.com/tony84727/diplomat/pkg/data"

type ConfigurationUpdater struct {
}

func (m ConfigurationUpdater) Set(key string, value string) error {
	return nil
}

func NewConfigurationUpdater(origin data.Configuration) *ConfigurationUpdater {
	return &ConfigurationUpdater{}
}

