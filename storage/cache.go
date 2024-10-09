/*
 * Copyright 2020 The SealEVM Authors
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package storage

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type Cache map[types.Slot]*evmInt256.Int

func (c Cache) Clone() Cache {
	replica := make(Cache)

	for k, v := range c {
		replica[k] = v.Clone()
	}

	return replica
}

func (c Cache) Merge(cache Cache) {
	for k, v := range cache {
		c[k] = v
	}
}

type SlotCache map[types.Address]Cache

func (c SlotCache) Clone() SlotCache {
	replica := SlotCache{}

	for addr, slotCache := range c {
		for slot, val := range slotCache {
			replica.Set(addr, slot, val.Clone())
		}
	}

	return replica
}

func (c SlotCache) Get(address types.Address, slot types.Slot) *evmInt256.Int {
	if c[address] == nil {
		return nil
	} else {
		return c[address][slot]
	}
}

func (c SlotCache) Set(address types.Address, slot types.Slot, v *evmInt256.Int) {
	if c[address] == nil {
		c[address] = Cache{}
	}

	c[address][slot] = v
}

func (c SlotCache) Merge(cache SlotCache) {
	for k, v := range cache {
		c[k] = v
	}
}

type balance struct {
	Address types.Address
	Balance *evmInt256.Int
}

type LogCache []*types.Log

func (l LogCache) Clone() LogCache {
	replica := make(LogCache, len(l))

	for i := 0; i < len(l); i++ {
		replica[i] = l[i].Clone()
	}

	return replica
}

type DestructCache map[types.Address]types.Address

func (d DestructCache) Clone() DestructCache {
	replica := DestructCache{}

	for k, v := range d {
		replica[k] = v
	}

	return replica
}

func (d DestructCache) Merge(cache DestructCache) {
	for k, v := range cache {
		d[k] = v
	}
}

type ResultCache struct {
	OriginalData SlotCache
	CachedData   SlotCache

	TOriginalData SlotCache
	TCachedData   SlotCache

	Logs         *LogCache
	Destructs    DestructCache
	NewContracts ContractCache
}

func NewResultCache() ResultCache {
	return ResultCache{
		OriginalData:  SlotCache{},
		CachedData:    SlotCache{},
		TOriginalData: SlotCache{},
		TCachedData:   SlotCache{},
		Logs:          &LogCache{},
		Destructs:     DestructCache{},
		NewContracts:  ContractCache{},
	}
}

func (r *ResultCache) Clone() ResultCache {
	logsClone := r.Logs.Clone()
	replica := ResultCache{
		OriginalData:  r.OriginalData.Clone(),
		CachedData:    r.CachedData.Clone(),
		TOriginalData: r.TOriginalData.Clone(),
		TCachedData:   r.TCachedData.Clone(),
		Logs:          &logsClone,
		Destructs:     r.Destructs.Clone(),
		NewContracts:  r.NewContracts.Clone(),
	}

	return replica
}

func (r *ResultCache) XOriginalLoad(address types.Address, slot types.Slot, t TypeOfStorage) *evmInt256.Int {
	if t == SStorage {
		return r.OriginalData.Get(address, slot)
	} else {
		return r.TOriginalData.Get(address, slot)
	}
}

func (r *ResultCache) XCachedLoad(address types.Address, slot types.Slot, t TypeOfStorage) *evmInt256.Int {
	if t == SStorage {
		return r.CachedData.Get(address, slot)
	} else {
		return r.TCachedData.Get(address, slot)
	}
}

func (r *ResultCache) XOriginalStore(address types.Address, slot types.Slot, val *evmInt256.Int, t TypeOfStorage) {
	if t == SStorage {
		r.OriginalData.Set(address, slot, val)
	} else {
		r.TOriginalData.Set(address, slot, val)
	}
}

func (r *ResultCache) XCachedStore(address types.Address, slot types.Slot, val *evmInt256.Int, t TypeOfStorage) {
	if t == SStorage {
		r.CachedData.Set(address, slot, val)
	} else {
		r.TCachedData.Set(address, slot, val)
	}
}

type CodeCache map[string][]byte

type readOnlyCache struct {
	BlockHash Cache
	Contracts ContractCache
}

func MergeResultCache(src *ResultCache, to *ResultCache) {
	to.OriginalData.Merge(src.OriginalData)
	to.CachedData.Merge(src.CachedData)
	to.TOriginalData.Merge(src.TOriginalData)
	to.TCachedData.Merge(src.TCachedData)
	to.Destructs.Merge(src.Destructs)
	to.NewContracts.Merge(src.NewContracts)

	*to.Logs = *src.Logs
}
