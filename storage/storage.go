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
	"errors"
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
	externalDataBlockStorage IExternalDataBlockStorage
}

func New(extStorage IExternalStorage, extDataBlockStorage IExternalDataBlockStorage) *Storage {
	s := &Storage{
		ResultCache:     cache.NewResultCache(),
		externalStorage: extStorage,
		externalDataBlockStorage: extDataBlockStorage,
		readOnlyCache: cache.ReadOnlyCache{
			BlockHash: map[types.Slot]*evmInt256.Int{},
		},
	}

	return s
}

func (s *Storage) Clone() *Storage {
	replica := &Storage{
		ResultCache:     s.ResultCache.Clone(),
		readOnlyCache:   s.readOnlyCache,
		externalStorage: s.externalStorage,
		externalDataBlockStorage: s.externalDataBlockStorage,
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
	if val.IsZero() {
		return nil
	}

	from, err := s.GetAccount(fromAddr)
	if err != nil {
		return err
	}

	to, err := s.GetAccount(toAddr)
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
}

func (s *Storage) Destruct(address types.Address) {
	s.ResultCache.Destructs[address] = address
}

func (s *Storage) GetAccount(address types.Address) (*environment.Account, error) {
	cachedAcc := s.ResultCache.CachedAccounts.Get(address)
	if cachedAcc != nil {
		return cachedAcc, nil
	}

	extAcc, err := s.externalStorage.GetAccount(address)
	if err != nil {
		return nil, err
	}

	if extAcc == nil {
		return nil, evmErrors.NoSuchDataInTheStorage(errors.New("external return nil"))
	}

	s.ResultCache.CacheAccount(extAcc)

	return extAcc, nil
}

func (s *Storage) AccountWithoutCache(addr types.Address) (*environment.Account, error) {
	acc, err := s.externalStorage.GetAccount(addr)
	if err != nil {
		return nil, err
	}

	if acc == nil {
		return nil, evmErrors.NoSuchDataInTheStorage(errors.New("external storage return nil"))
	}

	return acc, err
}

func (s *Storage) Balance(address types.Address) (*evmInt256.Int, error) {
	acc, err := s.GetAccount(address)
	if err != nil {
		return nil, err
	}

	if acc.Balance == nil {
		return evmInt256.New(0), nil
	}

	return acc.Balance.Clone(), nil
}

func (s *Storage) GetCode(address types.Address) ([]byte, error) {
	acc, err := s.GetAccount(address)
	if err != nil {
		return nil, err
	}

	if acc.Contract == nil {
		return nil, err
	}

	return acc.Contract.Code, err
}

func (s *Storage) GetCodeSize(address types.Address) (*evmInt256.Int, error) {
	acc, err := s.GetAccount(address)
	if err != nil {
		return nil, err
	}

	if acc.Contract == nil {
		return evmInt256.New(0), err
	}

	return evmInt256.New(acc.Contract.CodeSize), err
}

func (s *Storage) HashOfCode(code []byte) types.Hash {
	return s.externalStorage.HashOfCode(code)
}

func (s *Storage) GetCodeHash(address types.Address) (*types.Hash, error) {
	acc, err := s.GetAccount(address)
	if err != nil {
		return nil, err
	}

	if acc.Contract == nil {
		return &types.Hash{}, nil
	}

	return &acc.Contract.CodeHash, nil
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
	return s.externalStorage.AccountExist(addr)
}

func (s *Storage) ContractEmpty(addr types.Address) bool {
	return s.externalStorage.AccountEmpty(addr)
}

func (s *Storage) RemoveCachedAccount(addr types.Address) {
	s.ResultCache.RemoveAccount(addr)
}

func (s *Storage) CacheAccount(acc *environment.Account, newContract bool) {
	if s.ResultCache.CachedAccounts.Get(acc.Address) != nil {
		return
	}

	cached := s.ResultCache.CacheAccount(acc)

	if newContract {
		s.ResultCache.NewContractAccounts.Set(cached)
	}
}

func (s *Storage) UpdateAccountContract(address types.Address, code []byte) {
	acc := s.ResultCache.CachedAccounts.Get(address)
	if acc == nil {
		return
	}

	newContract := &environment.Contract{
		Code:     code,
		CodeHash: s.HashOfCode(code),
		CodeSize: uint64(len(code)),

		InitCode: bytes.Clone(acc.Contract.Code),
	}

	acc.Contract = newContract
}

type IDataBlockStorage interface {
	GetDataBlock(slot types.Slot) (types.Bytes, error)
	SetDataBlock(slot types.Slot, data types.Bytes)
}

type dataBlockStorage struct {
	address types.Address
	dataBlock types.DataBlock
	externalDataBlockStorage IExternalDataBlockStorage
}

func (s *dataBlockStorage) GetDataBlock(slot types.Slot) (types.Bytes, error) {
	if s.dataBlock[slot] != nil {
		return s.dataBlock[slot], nil
	}

	data, err := s.externalDataBlockStorage.GetDataBlock(s.address, slot)
	if err != nil {
		return nil, err
	}

	s.dataBlock[slot] = data
	return data, nil
}

func (s *dataBlockStorage) SetDataBlock(slot types.Slot, data types.Bytes) {
	s.dataBlock[slot] = data
}

func (s *Storage) CloneDataBlockStorage(address types.Address) IDataBlockStorage {
	if s.ResultCache.DataBlockCache[address] == nil {
		s.ResultCache.DataBlockCache[address] = make(types.DataBlock)
	}

	return &dataBlockStorage{
		address: address,
		dataBlock: s.ResultCache.DataBlockCache[address].Clone(),
		externalDataBlockStorage: s.externalDataBlockStorage,
	}
}

func (s *Storage) NewDataBlockStorage(address types.Address) IDataBlockStorage {
	if s.ResultCache.DataBlockCache[address] == nil {
		s.ResultCache.DataBlockCache[address] = make(types.DataBlock)
	}

	return &dataBlockStorage{
		address: address,
		dataBlock: s.ResultCache.DataBlockCache[address],
		externalDataBlockStorage: s.externalDataBlockStorage,
	}
}

func (s *Storage) HasExternalDataBlockStorage() bool {
	return s.externalDataBlockStorage != nil
}

func (s *Storage) GetDataBlockStorage(address types.Address) types.DataBlock {
	if s.ResultCache.DataBlockCache[address] != nil {
		return s.ResultCache.DataBlockCache[address]
	}

	dataBlock := make(types.DataBlock)
	s.ResultCache.DataBlockCache[address] = dataBlock
	return dataBlock
}
