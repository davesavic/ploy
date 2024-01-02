/*
Copyright © 2023 Dave Savic
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/davesavic/ploy/pkg/ploy"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [pipeline]",
	Short: "Run a ploy pipeline",
	Long:  `Run a ploy pipeline`,
	Run: func(cmd *cobra.Command, args []string) {
		cf, err := os.ReadFile("configuration.json")
		if err != nil {
			log.Fatal(err)
		}

		var cfg ploy.Config
		err = json.Unmarshal(cf, &cfg)
		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			log.Fatal("You must specify a pipeline or task to run")
		}

		for _, arg := range args {
			var executor ploy.PipelineExecutor
			local, _ := cmd.Flags().GetBool("local")

			if local {
				executor = &ploy.LocalPipelineExecutor{Config: cfg}
			} else {
				executor = &ploy.RemotePipelineExecutor{Config: cfg}
			}

			out, err := executor.Execute(arg)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s\n"+out, time.Now().UTC().Format("2006-01-02 15:04:05"))
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("local", "l", false, "Run the pipeline locally")
}
