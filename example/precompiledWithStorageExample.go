//go:build precompiled
// +build precompiled

package main

import (
	"encoding/hex"
	"fmt"
	"math"

	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/environment"
	exampleContracts "github.com/SealSC/SealEVM/example/contracts"
	precompiledwithstorage "github.com/SealSC/SealEVM/example/precompiledWithStorage"
	"github.com/SealSC/SealEVM/executionNote"
	"github.com/SealSC/SealEVM/precompiledContracts"
	"github.com/SealSC/SealEVM/types"
)

func main() {
	// Load SealEVM module
	SealEVM.Load()

	// Prepare external storage
	storage := newStorage()

	// Register cross-transaction data sharing precompiled contract
	crossTxDataShareAddr := types.Address{}

	crossTxDataShareAddr.SetBytes([]byte{0x02, 0x00, 0x01}) // Use address in reserved address space
	fmt.Println("crossTxDataShareAddr:", crossTxDataShareAddr.String())
	err := precompiledContracts.RegisterContractWithStorage(
		crossTxDataShareAddr,
		&precompiledwithstorage.CrossTxDataShare{},
	)
	if err != nil {
		fmt.Println("Register CrossTxDataShare contract failed:", err)
		return
	}

	// Load contract code
	codes := exampleContracts.LoadCrossTxDataShareExampleCodes()

	// Use zero address as caller
	caller := types.Address{}

	// Deploy contract
	contract := deployContracts(caller, codes.InitCode(), nil, storage)
	if contract == nil {
		return
	}

	fmt.Println("\n=== Testing Cross-Transaction Data Sharing ===")

	// Call shareData() to store data
	shareDataMsg := &environment.Message{
		Caller: caller,
		Data:   codes.ShareDataInput(caller),
	}
	evm := newEVM(shareDataMsg, &contract.Address, math.MaxUint64, storage, storage)
	ret, err := evm.Execute()
	if err != nil {
		fmt.Println("call shareData() failed:", err)
		return
	}

	// Save execution result
	storage.StoreResult(&ret.StorageCache)

	// Print shared data event logs
	fmt.Println("\nShared Data Event Logs:")
	printLogs(ret.StorageCache.Logs)

	// Call readSharedData() to read data
	readDataMsg := &environment.Message{
		Caller: caller,
		Data:   codes.ReadDataInput(caller),
	}
	evm = newEVM(readDataMsg, &contract.Address, math.MaxUint64, storage, storage)
	ret, err = evm.Execute()
	if err != nil {
		fmt.Println("call readSharedData() failed:", err)
		return
	}

	// Print read data event logs
	fmt.Println("\nRead Data Event Logs:")
	printLogs(ret.StorageCache.Logs)

	// Print call chain
	fmt.Println("\nStart traversing call chain")
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
	fmt.Println("Call chain traversal completed\n")
}
