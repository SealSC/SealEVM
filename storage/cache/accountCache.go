package cache

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type AccountCache map[types.Address]*environment.Account

func (c AccountCache) Clone() AccountCache {
	replica := make(AccountCache)
	for k, v := range c {
		replica[k] = v.Clone()
	}

	return replica
}

func (c AccountCache) Merge(cache AccountCache) {
	for k, v := range cache {
		if v == nil {
			continue
		}

		if c[k] == nil {
			c[k] = v
		} else {
			c[k].Set(v)
		}
	}
}

func (c AccountCache) Set(acc *environment.Account) {
	c[acc.Address] = acc
}

func (c AccountCache) Get(address types.Address) *environment.Account {
	if c[address] == nil {
		return nil
	}

	return c[address]
}

func (c AccountCache) SetSlot(address types.Address, slot types.Slot, value *evmInt256.Int) {
	if c[address] == nil {
		return
	}

	c[address].Slots[slot] = value
}

func (c AccountCache) GetSlot(address types.Address, slot types.Slot) *evmInt256.Int {
	if c[address] == nil {
		return nil
	}

	return c[address].Slots[slot]
}
