package cache

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type TransientCache map[types.Address]SlotCache

func (c TransientCache) Clone() TransientCache {
	replica := TransientCache{}

	for addr, slotCache := range c {
		for slot, val := range slotCache {
			replica.Set(addr, slot, val.Clone())
		}
	}

	return replica
}

func (c TransientCache) Get(address types.Address, slot types.Slot) *evmInt256.Int {
	if c[address] == nil {
		return nil
	} else {
		return c[address][slot]
	}
}

func (c TransientCache) Set(address types.Address, slot types.Slot, v *evmInt256.Int) {
	if c[address] == nil {
		c[address] = SlotCache{}
	}

	c[address][slot] = v
}

func (c TransientCache) Merge(cache TransientCache) {
	for k, v := range cache {
		c[k] = v
	}
}
