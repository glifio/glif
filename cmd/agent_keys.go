/*
Copyright © 2023 Glif LTD
*/
package cmd

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/glif-confidential/cli/fevm"
	"github.com/spf13/cobra"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage Glif Agent keys",
	Long:  ``,
}

func init() {
	agentCmd.AddCommand(keysCmd)
}

func deriveAddrFromPk(pk *ecdsa.PrivateKey) (common.Address, address.Address, error) {
	evmAddr, err := fevm.DeriveAddressFromPk(pk)
	if err != nil {
		log.Fatal(err)
	}

	fevmAddr, err := ethtypes.ParseEthAddress(evmAddr.String())
	if err != nil {
		log.Fatal(err)
	}

	delegatedAddr, err := fevmAddr.ToFilecoinAddress()
	if err != nil {
		log.Fatal(err)
	}

	return evmAddr, delegatedAddr, nil
}
