package main

import (
	"fmt"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
	"math"
)

func deployContracts(
	caller types.Address,
	initCode []byte,
	val *evmInt256.Int,
	storage *extStorage,
) *environment.Account {
	//create the EVM instance
	evm := newEVM(&environment.Message{
		Caller: caller,
		Value:  val,
		Data:   initCode,
	}, nil, math.MaxUint64, storage, storage)

	//execute deploy tx
	ret, err := evm.Execute()

	//check error
	if err != nil {
		fmt.Println("deploy contract failed:", err.Error())
		return nil
	}
	fmt.Println("deploy contract success@address:", ret.ContractAddress)

	//store the account result
	storage.StoreResult(&ret.StorageCache)

	//get the contract deployed through the initCode
	deployed, _ := storage.GetAccount(*ret.ContractAddress)
	return deployed
}
