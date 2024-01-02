/*
Copyright © 2023 Dave Savic
*/
package cmd

import (
	"encoding/json"
	"github.com/davesavic/ploy/pkg/ploy"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a ploy script",
	Long:  `Initialise a ploy script`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := ploy.Config{
			Params: map[string]string{
				"message": "hello, world!",
			},
			Servers: map[string]ploy.Server{
				"staging": {
					Host:       "111.111.111.111",
					Port:       22,
					User:       "ploy",
					PrivateKey: "/home/user/.ssh/id_rsa",
				},
			},
			Tasks: map[string][]string{
				"print-message": {
					"echo '{{message}}'",
				},
			},
			Pipelines: map[string]ploy.Pipeline{
				"say-hello": {
					Tasks: []string{
						"print-message",
					},
					Servers: []string{
						"staging",
					},
				},
			},
		}

		jsCfg, _ := json.MarshalIndent(cfg, "", "	")
		err := os.WriteFile("configuration.json", jsCfg, 0666)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
