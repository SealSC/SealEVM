package cache

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type ResultCache struct {
	OriginalAccounts AccountCache
	CachedAccounts   AccountCache

	Logs         *LogCache
	Destructs    DestructCache
	NewContracts AccountCache

	tOriginalData TransientCache
	tCachedData   TransientCache
}

func NewResultCache() ResultCache {
	return ResultCache{
		OriginalAccounts: AccountCache{},
		CachedAccounts:   AccountCache{},
		tOriginalData:    TransientCache{},
		tCachedData:      TransientCache{},
		Logs:             &LogCache{},
		Destructs:        DestructCache{},
		NewContracts:     AccountCache{},
	}
}

func MergeResultCache(result *ResultCache, to *ResultCache) {
	to.OriginalAccounts.Merge(result.OriginalAccounts)
	to.CachedAccounts.Merge(result.CachedAccounts)
	to.tOriginalData.Merge(result.tOriginalData)
	to.tCachedData.Merge(result.tCachedData)
	to.Destructs.Merge(result.Destructs)
	to.NewContracts.Merge(result.NewContracts)

	*to.Logs = *result.Logs
}

func (r *ResultCache) Clone() ResultCache {
	logsClone := r.Logs.Clone()
	replica := ResultCache{
		OriginalAccounts: r.OriginalAccounts.Clone(),
		CachedAccounts:   r.CachedAccounts.Clone(),
		tOriginalData:    r.tOriginalData.Clone(),
		tCachedData:      r.tCachedData.Clone(),
		Logs:             &logsClone,
		Destructs:        r.Destructs.Clone(),
		NewContracts:     r.NewContracts.Clone(),
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

func (r *ResultCache) CacheContract(contract *environment.Contract) {
	if r.CachedAccounts[contract.Address] != nil {
		return
	}

	original := NewAccountCacheUnit(contract)
	cached := NewAccountCacheUnit(contract)

	r.OriginalAccounts.Set(original)
	r.CachedAccounts.Set(cached)
}
