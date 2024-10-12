package main

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage/cache"
	"github.com/SealSC/SealEVM/types"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

type memStorage struct {
	accounts cache.AccountCache
}

func (r *memStorage) GetAccount(address types.Address) (*environment.Account, error) {
	if r.accounts[address] == nil {
		r.accounts[address] = environment.NewAccount(address, nil, nil)
	}
	return r.accounts[address], nil
}

func (r *memStorage) HashOfCode(code []byte) types.Hash {
	var ret types.Hash
	ret.SetBytes(crypto.Keccak256(code))
	return ret
}

func (r *memStorage) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	return nil, nil
}

func (r *memStorage) GetChainID() (*evmInt256.Int, error) {
	return evmInt256.New(0), nil
}

func (r *memStorage) CreateAddress(caller types.Address, tx environment.Transaction) types.Address {
	var addr types.Address
	addr.SetBytes(evmInt256.New(uint64(time.Now().UnixNano())).Bytes())
	return addr
}

func (r *memStorage) CreateFixedAddress(caller types.Address, salt types.Hash, code []byte, tx environment.Transaction) types.Address {
	var addr types.Address
	addr.SetBytes(evmInt256.New(uint64(time.Now().UnixNano())).Bytes())
	return addr
}

func (r *memStorage) ContractExist(address types.Address) bool {
	return r.accounts[address] != nil
}

func (r *memStorage) ContractEmpty(address types.Address) bool {
	return r.accounts[address] == nil
}

func (r *memStorage) Load(address types.Address, slot types.Slot) (*evmInt256.Int, error) {
	ret := evmInt256.New(0)
	if r.accounts[address] != nil {
		if r.accounts[address].Slots[slot] != nil {
			ret.Set(r.accounts[address].Slots[slot].Int)
		}
	}

	return ret, nil
}

func (r *memStorage) NewAccount(address types.Address, code []byte) error {
	r.accounts[address] = environment.NewAccount(address, nil, &environment.Contract{
		Code:     code,
		CodeSize: uint64(len(code)),
	})

	r.accounts[address].Contract.CodeHash.SetBytes(crypto.Keccak256(code))
	return nil
}
