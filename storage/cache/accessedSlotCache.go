package cache

import (
	"github.com/SealSC/SealEVM/types"
)

type accessedSlots map[types.Slot]struct{}

type AccessedSlotCache map[types.Address]accessedSlots

func (c AccessedSlotCache) Clone() AccessedSlotCache {
	replica := AccessedSlotCache{}
	for addr, slots := range c {
		for slot := range slots {
			replica.Add(addr, slot)
		}
	}

	return replica
}

func (c AccessedSlotCache) Has(address types.Address, slot types.Slot) bool {
	if c[address] == nil {
		return false
	}

	_, accessed := c[address][slot]
	return accessed
}

func (c AccessedSlotCache) Add(address types.Address, slot types.Slot) {
	if c[address] == nil {
		c[address] = accessedSlots{}
	}

	c[address][slot] = struct{}{}
}

func (c AccessedSlotCache) Merge(cache AccessedSlotCache) {
	for addr, slots := range cache {
		for slot := range slots {
			c.Add(addr, slot)
		}
	}
}
