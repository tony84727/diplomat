package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tony84727/diplomat/cmd/diplomat/internal"
	"os"
	"strings"
)

var configCmd = &cobra.Command{
	Use: "config [path]",
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
		fmt.Printf("%s => %v", args[0], value)
		return nil
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	configCmd.Flags().StringVar(&projectDir, "project", "","project directory")
}

var (
	projectDir string
)
