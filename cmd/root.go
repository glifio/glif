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

	"github.com/ethereum/go-ethereum/common"
	"github.com/glifio/cli/util"
	"github.com/glifio/go-pools/sdk"
	types "github.com/glifio/go-pools/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var PoolsSDK types.PoolsSDK

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.glif/config.toml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	xdgConfigHome := fmt.Sprintf("%s/.glif", home)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".glif" (without extension).
		viper.AddConfigPath(fmt.Sprintf("%s/.glif", home))
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

	if err := sdk.Init(
		rootCmd.Context(),
		&PoolsSDK,
		common.HexToAddress(viper.GetString("routes.agent-police")),
		common.HexToAddress(viper.GetString("routes.miner-registry")),
		common.HexToAddress(viper.GetString("routes.router")),
		common.HexToAddress(viper.GetString("routes.pool-registry")),
		common.HexToAddress(viper.GetString("routes.agent-factory")),
		common.HexToAddress(viper.GetString("routes.ifil")),
		common.HexToAddress(viper.GetString("routes.wfil")),
		common.HexToAddress(viper.GetString("routes.infinity-pool")),
		viper.GetString("ado.address"),
		// using the mock ADO for now
		"Mock",
		viper.GetString("daemon.rpc-url"),
		viper.GetString("daemon.token"),
	); err != nil {
		log.Fatalf("Error initializing Pools SDK: %v\n", err)
	}

	//TODO: check that $HOME/.config/glif exists and create if not
	if err := util.NewKeyStore(fmt.Sprintf("%s/keys.toml", xdgConfigHome)); err != nil {
		log.Fatal(err)
	}

	if err := util.NewAgentStore(fmt.Sprintf("%s/agent.toml", xdgConfigHome)); err != nil {
		log.Fatal(err)
	}
}
