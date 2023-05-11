package fevm

import (
	abigen "github.com/glif-confidential/abigen/bindings"
)

var (
	MethodBorrow      = "borrow"
	MethodPay         = "pay"
	MethodAddMiner    = "addMiner"
	MethodRemoveMiner = "removeMiner"
	MethodWithdraw    = "withdraw"
	MethodPushFunds   = "pushFunds"
	MethodPullFunds   = "pullFunds"
)

func MethodStrToBytes(methodStr string) ([4]byte, error) {
	abi, err := abigen.AgentMetaData.GetAbi()
	if err != nil {
		return [4]byte{}, err
	}

	var methodID [4]byte
	copy(methodID[:], abi.Methods[methodStr].ID[:4])

	return methodID, nil
}
