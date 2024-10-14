package main

import (
	"fmt"
	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
	"math"
)

func main() {
	//load SealEVM module
	SealEVM.Load()

	//prepare external storage
	storage := newStorage()

	//load codes
	codes := loadCodes()

	//using zero address as caller
	caller := types.Address{}

	//deploy contract, contract would be nil if deploy failed
	contract := deployContracts(caller, codes.initCode, nil, storage)
	if contract == nil {
		return
	}

	//call counter() to get current counter's value
	evm := newEVM(codes.CounterCallMessage(caller), &contract.Address, math.MaxUint64, storage)
	ret, err := evm.Execute()
	if err != nil {
		fmt.Println("call counter() failed:", err)
		return
	}

	//result of Counter(), would be 0
	counter := evmInt256.BytesDataToEVMInt(ret.ResultData)
	fmt.Println("counter: ", counter)

	//call increaseFor("example") to increase the counter
	evm = newEVM(codes.IncreaseCallMessage(caller), &contract.Address, math.MaxUint64, storage)
	ret, _ = evm.Execute()

	if err != nil {
		fmt.Println("call increaseFor(\"example\") failed:", err)
		return
	}

	//store the result
	storage.StoreResult(&ret.StorageCache)

	//the event logs
	printLogs(ret.StorageCache.Logs)

	//call counter() to get counter's value after increase
	evm = newEVM(codes.CounterCallMessage(caller), &contract.Address, math.MaxUint64, storage)
	ret, err = evm.Execute()

	//result of Counter(), would be 1
	counter = evmInt256.BytesDataToEVMInt(ret.ResultData)
	fmt.Println("counter: ", counter)
}
