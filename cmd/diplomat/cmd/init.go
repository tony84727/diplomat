package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const exampleOutlineFile = `version: '1'
preprocessors:
- type: chinese
  options:
    - mode: t2s
      from: zh-TW
      to: zh-CN
output:
  - selectors:
      - admin
      - manage
    templates:
      - type: js
        options:
          filename: "{{.Lang}}.locale.js"`

var (
	outline string
	initCmd = &cobra.Command{
		Use:   "init",
		Short: "init diplomat project",
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
			outlinePath := filepath.Join(projectDir, outline)
			outlineDir := filepath.Dir(outlinePath)
			if err := os.MkdirAll(outlineDir, 0755); err != nil {
				fmt.Println(err)
				return
			}
			if _, err := os.Stat(outlinePath); err != nil {
				fmt.Println(err)
				return
			}
			outlineFile, err := os.Create(outlinePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer outlineFile.Close()
			if _, err := outlineFile.WriteString(exampleOutlineFile); err != nil {
				fmt.Println(err)
			}
		},
	}
)

func init() {
	initCmd.Flags().StringVarP(&outline, "outline", "o", "outline.yaml", "path to configuration file of diplomat")
}
