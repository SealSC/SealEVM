package storage

import (
	"bytes"
	"github.com/SealSC/SealEVM/types"
)

type Contract struct {
	Address  types.Address
	Code     []byte
	CodeHash types.Hash
	CodeSize uint64
}

func (c Contract) Clone() *Contract {
	return &Contract{
		Address:  c.Address,
		Code:     bytes.Clone(c.Code),
		CodeHash: c.CodeHash,
		CodeSize: c.CodeSize,
	}
}

type ContractCache map[types.Address]*Contract

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
	c[contract.Address] = contract.Clone()
}

func (c ContractCache) Get(address types.Address) *Contract {
	if c[address] == nil {
		return nil
	}

	return c[address].Clone()
}
