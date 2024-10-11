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
	"bytes"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage/cache"
	"github.com/SealSC/SealEVM/types"
)

type Storage struct {
	ResultCache     cache.ResultCache
	readOnlyCache   cache.ReadOnlyCache
	externalStorage IExternalStorage
}

func New(extStorage IExternalStorage) *Storage {
	s := &Storage{
		ResultCache:     cache.NewResultCache(),
		externalStorage: extStorage,
		readOnlyCache: cache.ReadOnlyCache{
			BlockHash: cache.SlotCache{},
		},
	}

	return s
}

func (s *Storage) Clone() *Storage {
	replica := &Storage{
		ResultCache:     s.ResultCache.Clone(),
		readOnlyCache:   s.readOnlyCache,
		externalStorage: s.externalStorage,
	}

	return replica
}

func (s *Storage) XLoad(address types.Address, slot types.Slot, t cache.TypeOfStorage) (*evmInt256.Int, error) {
	if s.ResultCache.OriginalAccounts == nil || s.ResultCache.CachedAccounts == nil || s.externalStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	if t != cache.SStorage && t != cache.TStorage {
		return nil, evmErrors.InvalidTypeOfStorage()
	}

	var err error = nil
	i := s.ResultCache.XCachedLoad(address, slot, t)
	if i == nil {
		if t == cache.SStorage {
			i, err = s.externalStorage.Load(address, slot)
		} else {
			i = evmInt256.New(0)
		}

		if err != nil {
			return nil, evmErrors.NoSuchDataInTheStorage(err)
		}

		s.ResultCache.XCachedStore(address, slot, i, t)
		s.ResultCache.XOriginalStore(address, slot, i, t)
	}

	return i, nil
}

func (s *Storage) XStore(address types.Address, slot types.Slot, val *evmInt256.Int, t cache.TypeOfStorage) {
	s.ResultCache.XCachedStore(address, slot, val, t)
}

func (s *Storage) CanTransfer(from types.Address, to types.Address, amount *evmInt256.Int) bool {
	balance, err := s.Balance(from)
	if err != nil {
		return false
	}

	return balance.Cmp(amount.Int) >= 0
}

func (s *Storage) Transfer(fromAddr types.Address, toAddr types.Address, val *evmInt256.Int) error {
	from, err := s.getContract(fromAddr)
	if err != nil {
		return err
	}

	to, err := s.getContract(toAddr)
	if err != nil {
		return err
	}

	if from.Balance.LT(val) {
		return evmErrors.InsufficientBalance
	}

	from.Balance = from.Balance.Sub(val)
	to.Balance = to.Balance.Add(val)

	return nil
}

func (s *Storage) Log(log *types.Log) {
	*s.ResultCache.Logs = append(*s.ResultCache.Logs, log)

	return
}

func (s *Storage) Destruct(address types.Address) {
	s.ResultCache.Destructs[address] = address
}

func (s *Storage) getContract(address types.Address) (*environment.Contract, error) {
	newAcc := s.ResultCache.NewContracts.Get(address)
	if newAcc != nil {
		return newAcc.Contract, nil
	}

	cachedAcc := s.ResultCache.CachedAccounts.Get(address)
	if cachedAcc != nil {
		return cachedAcc.Contract, nil
	}

	extContract, err := s.externalStorage.GetContract(address)
	if err != nil {
		return nil, err
	}

	if extContract == nil {
		return nil, evmErrors.InvalidExternalStorageResult
	}

	cachedAcc = cache.NewAccountCacheUnit(extContract)
	s.ResultCache.CacheContract(extContract)

	return extContract, nil
}

func (s *Storage) Balance(address types.Address) (*evmInt256.Int, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return contract.Balance.Clone(), nil
}

func (s *Storage) GetCode(address types.Address) ([]byte, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return contract.Code, err
}

func (s *Storage) GetCodeSize(address types.Address) (*evmInt256.Int, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return evmInt256.New(contract.CodeSize), err
}

func (s *Storage) HashOfCode(code []byte) types.Hash {
	return s.externalStorage.HashOfCode(code)
}

func (s *Storage) GetCodeHash(address types.Address) (*types.Hash, error) {
	contract, err := s.getContract(address)
	if err != nil {
		return nil, err
	}

	return &contract.CodeHash, err
}

func (s *Storage) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	var slot types.Slot
	slot.SetBytes(block.Bytes())
	if hash, exists := s.readOnlyCache.BlockHash[slot]; exists {
		return hash, nil
	}

	hash, err := s.externalStorage.GetBlockHash(block)
	if err == nil {
		s.readOnlyCache.BlockHash[slot] = hash
	}

	return hash, err
}

func (s *Storage) NewContract(address types.Address, code []byte) {
	s.ResultCache.NewContracts.Set(&cache.AccountCacheUnit{
		Contract: &environment.Contract{
			Address:  address,
			Code:     bytes.Clone(code),
			CodeHash: s.externalStorage.HashOfCode(code),
			CodeSize: uint64(len(code)),
		},
		Slots: cache.SlotCache{},
	})
}

func (s *Storage) CreateAddress(caller types.Address, tx environment.Transaction) types.Address {
	return s.externalStorage.CreateAddress(caller, tx)
}

func (s *Storage) CreateFixedAddress(caller types.Address, salt types.Hash, code []byte, tx environment.Transaction) types.Address {
	return s.externalStorage.CreateFixedAddress(caller, salt, code, tx)
}

func (s *Storage) GetExternalStorage() IExternalStorage {
	return s.externalStorage
}

func (s *Storage) ClearCache() {
	s.ResultCache = cache.NewResultCache()
}

func (s *Storage) CachedContract(addr types.Address) bool {
	return s.ResultCache.CachedAccounts.Get(addr) != nil
}

func (s *Storage) CachedData(addr types.Address, slot types.Slot) (org *evmInt256.Int, current *evmInt256.Int) {
	org = s.ResultCache.CachedAccounts.GetSlot(addr, slot)
	current = s.ResultCache.CachedAccounts.GetSlot(addr, slot)
	return org, current
}

func (s *Storage) ContractExist(addr types.Address) bool {
	return s.externalStorage.ContractExist(addr)
}

func (s *Storage) ContractEmpty(addr types.Address) bool {
	return s.externalStorage.ContractEmpty(addr)
}
