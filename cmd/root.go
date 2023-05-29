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
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math/big"
	"os"

	"github.com/glifio/cli/util"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/deploy"
	"github.com/glifio/go-pools/sdk"
	types "github.com/glifio/go-pools/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgDir string
var useCalibnet bool // only set in root_calibnet.go
var PoolsSDK types.PoolsSDK

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "glif",
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
	rootCmd.PersistentFlags().StringVar(&cfgDir, "config-dir", "", "config directory")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgDir != "" {
		viper.AddConfigPath(cfgDir)
	} else if os.Getenv("GLIF_CONFIG_DIR") != "" {
		viper.AddConfigPath(os.Getenv("GLIF_CONFIG_DIR"))
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		cfgDir = fmt.Sprintf("%s/.glif", home)
		if useCalibnet {
			cfgDir = fmt.Sprintf("%s/.glif/%s", home, "calibnet")
		}

		// Search config in home directory with name ".glif" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.AddConfigPath(".")
	}

	viper.SetConfigType("toml")
	viper.SetConfigName("config")

	if err := util.NewKeyStore(fmt.Sprintf("%s/keys.toml", cfgDir)); err != nil {
		log.Fatal(err)
	}

	if err := util.NewAgentStore(fmt.Sprintf("%s/agent.toml", cfgDir)); err != nil {
		log.Fatal(err)
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

	chainID := viper.GetInt64("chain.chain-id")
	var extern types.Extern

	switch chainID {
	case constants.MainnetChainID:
		extern = deploy.Extern
	case constants.CalibnetChainID:
		extern = deploy.TestExtern
	default:
		log.Fatalf("Unknown chain id %d", chainID)
	}

	sdk, err := sdk.New(context.Background(), big.NewInt(chainID), extern)
	if err != nil {
		log.Fatalf("Failed to initialize pools sdk %s", err)
	}
	PoolsSDK = sdk
}
