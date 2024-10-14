package main

import (
	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/types"
	"time"
)

func newEVM(
	msg *environment.Message,
	to *types.Address,
	gas uint64,
	storage storage.IExternalStorage,
) *SealEVM.EVM {
	evm := SealEVM.New(SealEVM.EVMParam{
		MaxStackDepth:  0,
		ExternalStore:  storage,
		ResultCallback: nil,
		Context: &environment.Context{
			Block: environment.Block{
				ChainID:    evmInt256.New(0),
				Coinbase:   types.Address{},
				Timestamp:  uint64(time.Now().Second()),
				Number:     0,
				Difficulty: evmInt256.New(0),
				GasLimit:   evmInt256.New(gas),
				Hash:       types.Hash{},
			},
			Transaction: environment.Transaction{
				Origin:   msg.Caller,
				To:       to,
				GasPrice: evmInt256.New(1),
				GasLimit: evmInt256.New(gas),
			},
			Message: *msg,
		},
	})

	return evm
}
