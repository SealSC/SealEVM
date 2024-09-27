package storage

import (
	"bytes"
	"github.com/SealSC/SealEVM/evmInt256"
)

type Contract struct {
	Address  *evmInt256.Int
	Code     []byte
	CodeHash *evmInt256.Int
	CodeSize uint64
}

func (c Contract) Clone() *Contract {
	return &Contract{
		Address:  c.Address.Clone(),
		Code:     bytes.Clone(c.Code),
		CodeHash: c.CodeHash.Clone(),
		CodeSize: c.CodeSize,
	}
}

type ContractCache map[string]*Contract

func (c ContractCache) Clone() ContractCache {
	replica := make(ContractCache)
	for k, v := range c {
		replica[k] = v.Clone()
	}

	return replica
}

func (c ContractCache) Merge(cache ContractCache) {
	for k, v := range cache {
		if c[k] == nil {
			c[k] = v
		}
	}
}

func (c ContractCache) Set(contract *Contract) {
	c[contract.Address.AsStringKey()] = contract.Clone()
}

func (c ContractCache) Get(address *evmInt256.Int) *Contract {
	if c[address.AsStringKey()] == nil {
		return nil
	}

	return c[address.AsStringKey()].Clone()
}
