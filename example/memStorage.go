package main

import (
	"bytes"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/types"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

//external storage for example
type memStorage struct {
	storage   storage.SlotCache
	contracts storage.ContractCache
}

func (r *memStorage) GetBalance(address types.Address) (*evmInt256.Int, error) {
	return evmInt256.New(1000000000000000000), nil
}

func (r *memStorage) CanTransfer(from, to types.Address, val *evmInt256.Int) bool {
	return true
}

func (r *memStorage) GetContract(address types.Address) (*environment.Contract, error) {
	return r.contracts[address], nil
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
	return r.contracts[address] != nil
}

func (r *memStorage) ContractEmpty(address types.Address) bool {
	return r.contracts[address] == nil
}

func (r *memStorage) Load(address types.Address, slot types.Slot) (*evmInt256.Int, error) {
	ret := evmInt256.New(0)
	if val, exists := r.storage[address][slot]; exists {
		ret.Set(val.Int)
	}

	return ret, nil
}

func (r *memStorage) NewContract(address types.Address, code []byte) error {
	r.contracts[address] = &environment.Contract{
		Address:  address,
		Code:     bytes.Clone(code),
		CodeHash: r.HashOfCode(code),
		CodeSize: uint64(len(code)),
	}
	return nil
}
