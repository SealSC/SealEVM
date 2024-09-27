package storage

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/types"
)

type ContractCache map[types.Address]*environment.Contract

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

func (c ContractCache) Set(contract *environment.Contract) {
	c[contract.Address] = contract.Clone()
}

func (c ContractCache) Get(address types.Address) *environment.Contract {
	if c[address] == nil {
		return nil
	}

	return c[address].Clone()
}
