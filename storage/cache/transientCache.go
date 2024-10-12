package cache

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type transientSlot map[types.Slot]*evmInt256.Int

func (t transientSlot) Clone() transientSlot {
	replica := transientSlot{}
	for slot, v := range t {
		replica[slot] = v.Clone()
	}

	return replica
}

func (t transientSlot) Merge(cache transientSlot) {
	for slot, v := range cache {
		t[slot] = v
	}
}

type TransientCache map[types.Address]transientSlot

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
		c[address] = transientSlot{}
	}

	c[address][slot] = v
}

func (c TransientCache) Merge(cache TransientCache) {
	for k, v := range cache {
		c[k] = v
	}
}
