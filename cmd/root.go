/*
Copyright © 2023 Glif LTD

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/glif-confidential/cli/fevm"
	"github.com/glif-confidential/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "glif",
	Short: "",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/glif/config.toml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		//TODO: check that $HOME/.config/glif exists and create if not
		// create default config.toml
		// create empty agent.toml
		// create empty keys.toml

		// Search config in home directory with name ".glif" (without extension).
		viper.AddConfigPath(fmt.Sprintf("%s/.config/glif", home))
		viper.AddConfigPath(".")
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("No config file found at %s\n", viper.ConfigFileUsed())
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Fprintln(os.Stderr, "Warning: No config file found.")
		} else {
			log.Fatalf("Config file error: %v\n", err)
		}
	}

	// Pulls in the FEVM connection params
	if err := fevm.InitFEVMConnection(rootCmd.Context()); err != nil {
		log.Fatalf("Error initializing FEVM connection: %v\n", err)
	}

	fevm.Connection().InitNonceCache()

	//TODO: check that $HOME/.config/glif exists and create if not
	if err := util.NewKeyStore("keys.toml"); err != nil {
		log.Fatal(err)
	}

	if err := util.NewAgentStore("agent.toml"); err != nil {
		log.Fatal(err)
	}
}
