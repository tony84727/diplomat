package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tony84727/diplomat"
	"github.com/tony84727/diplomat/pkg/data"
	"github.com/tony84727/diplomat/pkg/emit"
	_ "github.com/tony84727/diplomat/pkg/emit/golang"
	_ "github.com/tony84727/diplomat/pkg/emit/javascript"
	"github.com/tony84727/diplomat/pkg/log"
	"github.com/tony84727/diplomat/pkg/parser/yaml"
	"github.com/tony84727/diplomat/pkg/prepros"
	_ "github.com/tony84727/diplomat/pkg/prepros/chinese"
	_ "github.com/tony84727/diplomat/pkg/prepros/copy"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	watch    bool
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "build",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewColoredLogger()
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
			outDir := filepath.Join(projectDir, "out")
			sourceSet := data.NewFileSystemSourceSet(projectDir)
			configFile, err := sourceSet.GetConfigurationFile()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			configData, err := ioutil.ReadFile(configFile)
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			configParser := yaml.NewConfigurationParser(configData)
			config, err := configParser.GetConfiguration()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			preprocessorConfigs := config.GetPreprocessors()
			preprocessorFactory := prepros.NewComposeFactory(prepros.GlobalRegistry, preprocessorConfigs...)

			allTranslation := data.NewTranslationMerger(data.NewTranslation(""))
			translationFiles, err := sourceSet.GetTranslationFiles()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			for _, t := range translationFiles {
				content, err := ioutil.ReadFile(t)
				if err != nil {
					logger.Error(err.Error())
					os.Exit(1)
				}
				parser := yaml.NewParser(content)
				translation, err := parser.GetTranslation()
				if err != nil {
					logger.Error(err.Error())
					os.Exit(1)
				}
				allTranslation.Merge(translation)
			}
			if err := preprocessorFactory.Build()(allTranslation); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			synthesizer := diplomat.NewSynthesizer(outDir, allTranslation, emit.GlobalRegistry, logger)
			for _, o := range config.GetOutputs() {
				err := synthesizer.Output(o)
				if err != nil {
					logger.Error(err.Error())
					os.Exit(1)
				}
			}
		},
	}
)

func init() {
	buildCmd.Flags().BoolVar(&watch, "watch", false, "watch changes")
}
