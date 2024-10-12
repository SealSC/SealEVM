package main

import (
	"encoding/hex"
	"fmt"
	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/crypto/hashes"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage/cache"
	"github.com/SealSC/SealEVM/types"
	"os"
	"time"
)

func logPrinter(logCache *cache.LogCache) {
	for _, l := range *logCache {
		for _, t := range l.Topics {
			fmt.Println("topic:", t)
		}
		fmt.Println("data:", l.Data)
		fmt.Println("data as string:", string(l.Data))
	}
}

//store result to memStorage
func storeResult(result *SealEVM.ExecuteResult, storage *memStorage) {
	for addr, cachedAcc := range result.StorageCache.CachedAccounts {
		storage.accounts[addr] = cachedAcc
	}
}

//create a new evm
func newEvm(code []byte, callData []byte, caller []byte, ms *memStorage) *SealEVM.EVM {
	hash := hashes.Keccak256(code)
	var addr types.Address
	addr.SetBytes(hash)

	_ = ms.NewAccount(addr, code)

	var callerAddr types.Address
	callerAddr.SetBytes(caller)
	evm := SealEVM.New(SealEVM.EVMParam{
		MaxStackDepth:  0,
		ExternalStore:  ms,
		ResultCallback: nil,
		Context: &environment.Context{
			Block: environment.Block{
				ChainID:    evmInt256.New(0),
				Coinbase:   types.Address{},
				Timestamp:  uint64(time.Now().Second()),
				Number:     0,
				Difficulty: evmInt256.New(0),
				GasLimit:   evmInt256.New(10000000),
				Hash:       types.Hash{},
			},
			Transaction: environment.Transaction{
				Origin:   callerAddr,
				To:       &addr,
				GasPrice: evmInt256.New(1),
				GasLimit: evmInt256.New(10000000),
			},
			Message: environment.Message{
				Caller: callerAddr,
				Value:  evmInt256.New(0),
				Data:   callData,
			},
		},
	})

	return evm
}

func main() {
	fmt.Println("The example code is not yet adapted！")
	os.Exit(0)

	//load SealEVM module
	SealEVM.Load()

	//create memStorage
	ms := &memStorage{}
	ms.accounts = cache.AccountCache{}

	//deploy contract
	evm := newEvm(deployCode, nil, caller, ms)
	ret, err := evm.ExecuteContract()

	//check error
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	storeResult(&ret, ms)

	//result data of ret is the deployed code of example contract
	contractCode := ret.ResultData

	//call Counter() to get current counter's value
	evm = newEvm(contractCode, callCounter, caller, ms)
	ret, _ = evm.ExecuteContract()

	//result of Counter()
	fmt.Println("counter: ", hex.EncodeToString(ret.ResultData))

	//call increaseFor("example")
	evm = newEvm(contractCode, callIncreaseFor, caller, ms)
	ret, _ = evm.ExecuteContract()

	//store the result to ms
	storeResult(&ret, ms)

	//the event logs
	logPrinter(ret.StorageCache.Logs)

	//call Counter to get counter's value after increase
	evm = newEvm(contractCode, callCounter, caller, ms)
	ret, err = evm.ExecuteContract()

	//result of Counter()
	fmt.Println("counter: ", hex.EncodeToString(ret.ResultData))
}
