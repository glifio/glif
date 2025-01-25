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
	"runtime/debug"

	"github.com/ethereum/go-ethereum/common"
	jnal "github.com/glifio/glif/v2/journal"
	"github.com/glifio/glif/v2/journal/fsjournal"
	"github.com/glifio/glif/v2/util"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/deploy"
	"github.com/glifio/go-pools/sdk"
	types "github.com/glifio/go-pools/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var cfgDir string
var useCalibnet bool // only set in root_calibnet.go
var chainID int64 = constants.MainnetChainID
var PoolsSDK types.PoolsSDK
var journal jnal.Journal

var CommitHash, GoPoolsHash = func() (string, string) {
	var ch string
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				ch = setting.Value
			}
		}

		for _, dep := range info.Deps {
			if dep.Path == "github.com/glifio/go-pools" {
				return ch, dep.Version
			}
		}
	}
	return "", ""
}()

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
	rootCmd.PersistentFlags().Float64("gas-premium-multiply", 1.0, "Multiply the default gas premium by this amount")
	rootCmd.PersistentFlags().Uint64("nonce", 0, "Specify nonce (for replacing transactions)")
	rootCmd.PersistentFlags().Int64("gas-premium", -1, "(advanced) Override gas premium / priority fee per gas")
	rootCmd.PersistentFlags().Uint64("gas-limit", 0, "(advanced) Override gas limit from estimate")
	rootCmd.PersistentFlags().Uint64("gas-fee-cap", 0, "(advanced) Override fee cap / max fee per gas")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if os.Getenv("GLIF_CONFIG_DIR") != "" {
		cfgDir = os.Getenv("GLIF_CONFIG_DIR")
	}
	if cfgDir != "" {
		viper.AddConfigPath(cfgDir)
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

	var err error
	if journal, err = fsjournal.OpenFSJournal(cfgDir, nil); err != nil {
		logFatal(err)
	}

	util.NewKeyStore(fmt.Sprintf("%s/keystore", cfgDir))

	if err := util.NewKeyStoreLegacy(fmt.Sprintf("%s/keys.toml", cfgDir)); err != nil {
		logFatal(err)
	}

	if err := util.NewAgentStore(fmt.Sprintf("%s/agent.toml", cfgDir)); err != nil {
		logFatal(err)
	}

	if err := util.NewAccountsStore(fmt.Sprintf("%s/accounts.toml", cfgDir)); err != nil {
		logFatal(err)
	}

	if err := util.NewBackupsStore(fmt.Sprintf("%s/backups.toml", cfgDir)); err != nil {
		logFatal(err)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logFatalf("No config file found at %s\n", viper.ConfigFileUsed())
		} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// No .glif/config.toml, populate with mainnet defaults
			viper.Set("daemon.rpc-url", deploy.Extern.LotusDialAddr)
			viper.Set("daemon.token", "")
			viper.SafeWriteConfig()
		} else {
			logFatalf("Config file error: %v\n", err)
		}
	}

	viper.WatchConfig()

	if slices.Contains(os.Args[1:], "wallet") &&
		(slices.Contains(os.Args[1:], "create-agent-accounts") ||
			slices.Contains(os.Args[1:], "create-account") ||
			slices.Contains(os.Args[1:], "migrate")) {
		// Skip migration check
	} else {
		err = checkWalletMigrated()
		if err != nil {
			logFatal(err)
		}

		err = checkUnencryptedPrivateKeys()
		if err != nil {
			log.Println(err)
		}

		err = confirmBackupExists()
		if err != nil {
			logFatal(err)
		}
	}

	daemonURL := viper.GetString("daemon.rpc-url")
	daemonToken := viper.GetString("daemon.token")
	adoURL := viper.GetString("ado.address")
	eventsURL := viper.GetString("routes.events-url")

	override := viper.GetBool("routes.override")
	if override || chainID == constants.LocalnetChainID || chainID == constants.AnvilChainID {
		routerAddr := viper.GetString("routes.router")
		router := common.HexToAddress(routerAddr)
		err := sdk.LazyInit(context.Background(), &PoolsSDK, router, adoURL, "ADO", daemonURL, daemonToken, eventsURL)
		if err != nil {
			logFatal(err)
		}
	} else {
		var extern types.Extern
		switch chainID {
		case constants.MainnetChainID:
			extern = deploy.Extern
		case constants.CalibnetChainID:
			extern = deploy.TestExtern
		default:
			logFatalf("Unknown chain id %d", chainID)
		}

		if daemonURL != "" {
			extern.LotusDialAddr = daemonURL
		}
		if daemonToken != "" {
			extern.LotusToken = daemonToken
		}
		if eventsURL != "" {
			extern.EventsURL = eventsURL
		}
		if adoURL != "" {
			extern.AdoAddr = adoURL
		}

		sdk, err := sdk.New(context.Background(), big.NewInt(chainID), extern)
		if err != nil {
			logFatalf("Failed to initialize pools sdk %s", err)
		}
		PoolsSDK = sdk
	}
}
