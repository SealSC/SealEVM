package main

import (
	"errors"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"time"
)

//external storage for example
type memStorage struct {
	storage   map[string][]byte
	contracts map[string][]byte
}

func (r *memStorage) GetBalance(address *evmInt256.Int) (*evmInt256.Int, error) {
	return evmInt256.New(1000000000000000000), nil
}

func (r *memStorage) CanTransfer(from, to, val *evmInt256.Int) bool {
	return true
}

func (r *memStorage) GetCode(address *evmInt256.Int) ([]byte, error) {
	return r.contracts[address.AsStringKey()], nil
}

func (r *memStorage) GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error) {
	code, exist := r.contracts[address.AsStringKey()]
	if !exist {
		return nil, errors.New("no code for: 0x" + address.Text(16))
	}
	return evmInt256.New(int64(len(code))), nil
}

func (r *memStorage) GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error) {
	return nil, nil
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

func (r *memStorage) Load(n string, k string) (*evmInt256.Int, error) {
	ret := evmInt256.New(0)
	if val, exists := r.storage[n+k]; exists {
		ret.SetBytes(val)
	}

	return ret, nil
}

func (r *memStorage) NewContract(n string, code []byte) error {
	r.contracts[n] = code
	return nil
}
