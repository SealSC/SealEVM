package cache

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type SlotCache map[types.Slot]*evmInt256.Int

func (c SlotCache) Clone() SlotCache {
	replica := make(SlotCache)

	for k, v := range c {
		replica[k] = v.Clone()
	}

	return replica
}

func (c SlotCache) Merge(cache SlotCache) {
	for k, v := range cache {
		c[k] = v
	}
}

type AccountCacheUnit struct {
	Contract *environment.Contract
	Slots    SlotCache
}

func NewAccountCacheUnit(extContract *environment.Contract) *AccountCacheUnit {
	return &AccountCacheUnit{
		Contract: extContract,
		Slots:    SlotCache{},
	}
}

func (c AccountCacheUnit) Clone() *AccountCacheUnit {
	return &AccountCacheUnit{
		Contract: c.Contract.Clone(),
		Slots:    c.Slots.Clone(),
	}
}

type AccountCache map[types.Address]*AccountCacheUnit

func (c AccountCache) Clone() AccountCache {
	replica := make(AccountCache)
	for k, v := range c {
		replica[k] = v.Clone()
	}

	return replica
}

func (c AccountCache) Merge(cache AccountCache) {
	for k, v := range cache {
		c[k] = v
	}
}

func (c AccountCache) Set(contract *AccountCacheUnit) {
	c[contract.Contract.Address] = contract
}

func (c AccountCache) Get(address types.Address) *AccountCacheUnit {
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
