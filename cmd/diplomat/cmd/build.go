package cmd

import (
	"fmt"
	"github.com/insufficientchocolate/diplomat/internal"
	"github.com/insufficientchocolate/diplomat/pkg/data"
	"github.com/insufficientchocolate/diplomat/pkg/parser/yaml"
	"github.com/insufficientchocolate/diplomat/pkg/prepros"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

var (
	watch bool
	buildCmd = &cobra.Command{
		Use: "build",
		Short: "build",
		Run: func(cmd *cobra.Command, args []string) {
			var projectDir string
			if len(args) > 0 {
				projectDir = args[0]
			} else {
				var err error
				projectDir, err = os.Getwd()
				if err != nil {
					fmt.Println("cannot get current working directory", err)
					os.Exit(1)
				}
			}
			//outDir := filepath.Join(projectDir, "out")
			sourceSet := data.NewFileSystemSourceSet(projectDir)
			configFile, err := sourceSet.GetConfigurationFile()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			configData, err := ioutil.ReadFile(configFile)
			configParser := yaml.NewConfigurationParser(configData)
			config, err := configParser.GetConfiguration()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			preprocessorConfigs := config.GetPreprocessors()
			preprocessorInstances := make([]internal.PreprocessorFunc,0,len(preprocessorConfigs))
			// reverse order
			for i := len(preprocessorConfigs) -1; i >= 0; i-- {
				p := preprocessorConfigs[i]
				if instance := prepros.Manager.Get(p.GetType()); instance != nil {
					preprocessorInstances = append(preprocessorInstances, func(translation data.Translation) error {
						return instance.Process(translation, p.GetOptions())
					})
				}
			}
			preprocessorPipeline := internal.ComposePreprocessorFunc(preprocessorInstances...)
			allTranslation := data.NewTranslationMerger(data.NewTranslation(""))
			translationFiles, err := sourceSet.GetTranslationFiles()
			if err != nil {
				os.Exit(1)
			}
			for _, t := range translationFiles {
				content, err := ioutil.ReadFile(t)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				parser := yaml.NewParser(content)
				translation, err := parser.GetTranslation()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				allTranslation.Merge(translation)
			}
			if err := preprocessorPipeline(allTranslation); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			walker := data.NewTranslationWalker(allTranslation)
			if err := walker.ForEachTextNode(func(paths []string, textNode data.Translation) error {
				fmt.Println(strings.Join(paths, ".") + " => " + *textNode.GetText())
				return nil
			}); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
)

func init() {
	buildCmd.Flags().BoolVar(&watch, "watch",false,"watch changes")
}
