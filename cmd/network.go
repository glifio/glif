package cmd

type Network string

var (
	Testnet Network = "calibrationnet"
	Mainnet Network = "mainnet"
)

var network string

func NetworkName() Network {
	return Network(network)
}

var mainnetChainID = uint64(314)
var calibnetChainID = uint64(314159)
