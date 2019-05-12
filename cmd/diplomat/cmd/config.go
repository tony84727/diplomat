package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tony84727/diplomat"
	"github.com/tony84727/diplomat/cmd/diplomat/internal"
	"github.com/tony84727/diplomat/pkg/parser/yaml"
	"io/ioutil"
	"os"
	"strings"
)

var configCmd = &cobra.Command{
	Use: "config <path> [value]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(projectDir) <= 0 {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			projectDir = pwd
		}
		project := internal.NewProject(projectDir)
		config, err := project.LoadConfig()
		if err != nil {
			return err
		}
		navigator := internal.NewConfigNavigator(config)
		value, err := navigator.Get(strings.Split(args[0], ".")...)
		if err != nil {
			return err
		}
		if len(args) > 1 {
			updater := internal.NewConfigurationUpdater(config)
			if err := updater.Set(args[0], args[1]); err != nil {
				return err
			}
			newConfigContent, err := yaml.Write(updater.Config)
			if err != nil {
				return err
			}
			configPath, err := project.GetConfigurationFile()
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(configPath, newConfigContent, diplomat.DefaultFilePerm); err != nil {
				return err
			}
			fmt.Println("new configuration written")
		} else {
			fmt.Printf("%v", value)
		}
		return nil
	},
	Args: cobra.MaximumNArgs(2),
}

func init() {
	configCmd.Flags().StringVar(&projectDir, "project", "","project directory")
}

var (
	projectDir string
)
