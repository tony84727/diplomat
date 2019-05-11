package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tony84727/diplomat"
	"github.com/tony84727/diplomat/cmd/diplomat/internal"
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
)

func findProject(args []string) (*internal.Project, error) {
	if len(args) > 0 {
		return internal.NewProject(args[0]), nil
	}
	return internal.FindProject(nil)
}

var (
	watch    bool
	printToStdOut bool
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "build",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.NewColoredLogger()
			project, err := findProject(args)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			outDir := project.Of("out")
			config, err := project.LoadConfig()
			if err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
			preprocessorConfigs := config.GetPreprocessors()
			preprocessorFactory := prepros.NewComposeFactory(prepros.GlobalRegistry, preprocessorConfigs...)

			allTranslation := data.NewTranslationMerger(data.NewTranslation(""))
			translationFiles, err := project.GetTranslationFiles()
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
			var output diplomat.Output
			if printToStdOut {
				output = diplomat.ConsoleOutput{}
			} else {
				output = diplomat.NewOutputDirectory(outDir)
			}
			synthesizer := diplomat.NewSynthesizer(output, allTranslation, emit.GlobalRegistry, logger)
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
	buildCmd.Flags().BoolVarP(&printToStdOut, "print", "p", false,"print outputs to stdout")
}
