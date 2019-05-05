package internal

import (
	"github.com/tony84727/diplomat/pkg/data"
)

type ConfigNavigator struct {
	config data.Configuration
}

func NewConfigNavigator(config data.Configuration) *ConfigNavigator {
	return &ConfigNavigator{config: config}
}

func (c ConfigNavigator) Get(paths ...string) (interface{}, error) {
}
