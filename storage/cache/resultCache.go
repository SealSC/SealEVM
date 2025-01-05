package cache

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type ResultCache struct {
	OriginalAccounts AccountCache
	CachedAccounts   AccountCache

	NewContractAccounts AccountCache

	Logs      *LogCache
	Destructs DestructCache

	DataBlockCache DataBlockCache

	tOriginalData TransientCache
	tCachedData   TransientCache
}

func NewResultCache() ResultCache {
	return ResultCache{
		OriginalAccounts:    AccountCache{},
		CachedAccounts:      AccountCache{},
		NewContractAccounts: AccountCache{},
		Logs:                &LogCache{},
		Destructs:           DestructCache{},
		DataBlockCache:      DataBlockCache{},

		tOriginalData: TransientCache{},
		tCachedData:   TransientCache{},
	}
}

func MergeResultCache(result *ResultCache, to *ResultCache) {
	to.OriginalAccounts.Merge(result.OriginalAccounts)
	to.CachedAccounts.Merge(result.CachedAccounts)
	to.NewContractAccounts.Merge(result.NewContractAccounts)
	to.Destructs.Merge(result.Destructs)
	to.DataBlockCache.Merge(result.DataBlockCache)

	*to.Logs = *result.Logs

	to.tOriginalData.Merge(result.tOriginalData)
	to.tCachedData.Merge(result.tCachedData)
}

func (r *ResultCache) Clone() ResultCache {
	logsClone := r.Logs.Clone()
	replica := ResultCache{
		OriginalAccounts:    r.OriginalAccounts.Clone(),
		CachedAccounts:      r.CachedAccounts.Clone(),
		NewContractAccounts: AccountCache{},

		Logs:      &logsClone,
		Destructs: r.Destructs.Clone(),
		DataBlockCache: r.DataBlockCache.Clone(),
		tOriginalData: r.tOriginalData.Clone(),
		tCachedData:   r.tCachedData.Clone(),
	}

	for addr, acc := range r.NewContractAccounts {
		acc = replica.CachedAccounts.Get(addr)
		if acc != nil {
			replica.NewContractAccounts[addr] = acc
		}
	}

	return replica
}

func (r *ResultCache) XCachedLoad(address types.Address, slot types.Slot, t TypeOfStorage) *evmInt256.Int {
	if t == SStorage {
		return r.CachedAccounts.GetSlot(address, slot)
	} else {
		return r.tCachedData.Get(address, slot)
	}
}

func (r *ResultCache) XOriginalStore(address types.Address, slot types.Slot, val *evmInt256.Int, t TypeOfStorage) {
	if t == SStorage {
		r.OriginalAccounts.SetSlot(address, slot, val)
	} else {
		r.tOriginalData.Set(address, slot, val)
	}
}

func (r *ResultCache) XCachedStore(address types.Address, slot types.Slot, val *evmInt256.Int, t TypeOfStorage) {
	if t == SStorage {
		r.CachedAccounts.SetSlot(address, slot, val)
	} else {
		r.tCachedData.Set(address, slot, val)
	}
}

func (r *ResultCache) RemoveAccount(addr types.Address) {
	if r.CachedAccounts[addr] == nil {
		return
	}

	delete(r.CachedAccounts, addr)
	delete(r.NewContractAccounts, addr)
}

func (r *ResultCache) CacheAccount(acc *environment.Account) *environment.Account {
	if r.CachedAccounts[acc.Address] != nil {
		return r.CachedAccounts[acc.Address]
	}

	cached := acc.Clone()

	r.OriginalAccounts.Set(acc.Clone())
	r.CachedAccounts.Set(cached)

	return cached
}
