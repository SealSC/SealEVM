package main

import (
	"bytes"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

//external storage for example
type memStorage struct {
	storage   map[string][]byte
	contracts storage.ContractCache
}

func (r *memStorage) GetBalance(address *evmInt256.Int) (*evmInt256.Int, error) {
	return evmInt256.New(1000000000000000000), nil
}

func (r *memStorage) CanTransfer(from, to, val *evmInt256.Int) bool {
	return true
}

func (r *memStorage) GetContract(address *evmInt256.Int) (*storage.Contract, error) {
	return r.contracts[address.AsStringKey()], nil
}

func (r *memStorage) HashOfCode(code []byte) *evmInt256.Int {
	ret := evmInt256.New(0)
	ret.SetBytes(crypto.Keccak256(code))
	return ret
}

func (r *memStorage) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	return nil, nil
}

func (r *memStorage) GetChainID() (*evmInt256.Int, error) {
	return evmInt256.New(0), nil
}

func (r *memStorage) CreateAddress(caller *evmInt256.Int, tx environment.Transaction) *evmInt256.Int {
	return evmInt256.New(time.Now().UnixNano())
}

func (r *memStorage) CreateFixedAddress(caller *evmInt256.Int, salt *evmInt256.Int, code []byte, tx environment.Transaction) *evmInt256.Int {
	return evmInt256.New(time.Now().UnixNano())
}

func (r *memStorage) Load(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error) {
	ret := evmInt256.New(0)
	if val, exists := r.storage[n.AsStringKey()+k.AsStringKey()]; exists {
		ret.SetBytes(val)
	}

	return ret, nil
}

func (r *memStorage) NewContract(n *evmInt256.Int, code []byte) error {
	r.contracts[n.AsStringKey()] = &storage.Contract{
		Address:  n.Clone(),
		Code:     bytes.Clone(code),
		CodeHash: r.HashOfCode(code),
		CodeSize: uint64(len(code)),
	}
	return nil
}
