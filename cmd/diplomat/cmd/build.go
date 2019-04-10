package cmd

import (
	"fmt"
	"github.com/MinecraftXwinP/diplomat"
	"github.com/spf13/cobra"
	"os"
)

var (
	buildCmd = &cobra.Command{
		Use: "build",
		Short: "build",
		Run: func(cmd *cobra.Command, args []string) {
			var projectDir string
			if len(args) >= 0 {
				projectDir = args[0]
			} else {
				var err error
				projectDir, err = os.Getwd()
				if err != nil {
					fmt.Println("cannot get current working directory", err)
					return
				}
			}
			d, err := diplomat.NewDiplomatForDirectory(projectDir)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			if err := d.Output("out"); err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		},
	}
)
