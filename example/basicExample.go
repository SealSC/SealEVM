//go:build basic
// +build basic

package main

import (
	"encoding/hex"
	"fmt"
	"math"

	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/executionNote"
	"github.com/SealSC/SealEVM/types"
	"github.com/SealSC/SealEVM/example/contracts"
)

func main() {
	//load SealEVM module
	SealEVM.Load()

	//prepare external storage
	storage := newStorage()

	//load codes
	codes := exampleContracts.LoadBasicExampleCodes()

	//using zero address as caller
	caller := types.Address{}

	//deploy contract, contract would be nil if deploy failed
	contract := deployContracts(caller, codes.InitCode(), nil, storage)
	if contract == nil {
		return
	}

	//call counter() to get current counter's value
	evm := newEVM(codes.CounterCallMessage(caller), &contract.Address, math.MaxUint64, storage, storage)
	ret, err := evm.Execute()
	if err != nil {
		fmt.Println("call counter() failed:", err)
		return
	}

	//result of Counter(), would be 0
	counter := evmInt256.BytesDataToEVMInt(ret.ResultData)
	fmt.Println("counter: ", counter)

	//call increaseFor("example") to increase the counter
	evm = newEVM(codes.IncreaseCallMessage(caller), &contract.Address, math.MaxUint64, storage, storage)
	ret, _ = evm.Execute()

	if err != nil {
		fmt.Println("call increaseFor(\"example\") failed:", err)
		return
	}

	//store the result
	storage.StoreResult(&ret.StorageCache)

	//Traversing the call chain
	fmt.Println("\nstart traversing the call chain")
	if ret.Note != nil {
		ret.Note.Walk(func(note *executionNote.Note, depth uint64) {
			fmt.Println("\n----------------------->")
			fmt.Println("  depth:", depth)
			fmt.Println("   type:", note.Type)
			fmt.Println("   from:", note.From)
			fmt.Println("     to:", note.To)
			fmt.Println("    gas:", note.Gas)
			fmt.Println("    val:", note.Val)
			fmt.Println("   data:", hex.EncodeToString(note.Input))
			fmt.Println("execErr:", note.ExecutionError)
			fmt.Println("execRet:", hex.EncodeToString(note.ReturnData))
			fmt.Println("<-----------------------")
		})
	}
	fmt.Println("traversing the call chain ended\n")

	//the event logs
	printLogs(ret.StorageCache.Logs)

	//call counter() to get counter's value after increase
	evm = newEVM(codes.CounterCallMessage(caller), &contract.Address, math.MaxUint64, storage, storage)
	ret, err = evm.Execute()

	//result of Counter(), would be 1
	counter = evmInt256.BytesDataToEVMInt(ret.ResultData)
	fmt.Println("counter: ", counter)
}
