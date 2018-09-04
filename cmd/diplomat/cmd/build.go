package cmd

import (
	"fmt"
	"github.com/insufficientchocolate/diplomat"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"path/filepath"
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
	watchCmd = &cobra.Command{
		Use: "watch",
		Short: "watch",
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
			outDir := filepath.Join(projectDir, "out")
			d, errChan, changeListener := diplomat.NewDiplomatWatchDirectory(projectDir)
			go func() {
				for e := range errChan {
					log.Println("error:", e)
				}
			}()
			go func() {
				for range changeListener {
					d.Output(outDir)
				}
			}()
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit
		},
	}
)
