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
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
)

type Cache map[string]*evmInt256.Int
type CacheUnderNamespace map[string]Cache

func (c CacheUnderNamespace) Get(namespace string, key string) *evmInt256.Int {
	if c[namespace] == nil {
		return nil
	} else {
		return c[namespace][key]
	}
}

func (c CacheUnderNamespace) Set(namespace string, key string, v *evmInt256.Int) {
	if c[namespace] == nil {
		c[namespace] = Cache{}
	}

	c[namespace][key] = v
}

type balance struct {
	Address *evmInt256.Int
	Balance *evmInt256.Int
}

type BalanceCache map[string]*balance

type Log struct {
	Topics  [][]byte
	Data    []byte
	Context environment.Context
}

type LogCache map[string][]Log

type ResultCache struct {
	OriginalData CacheUnderNamespace
	CachedData   CacheUnderNamespace

	TOriginalData CacheUnderNamespace
	TCachedData   CacheUnderNamespace

	Balance   BalanceCache
	Logs      LogCache
	Destructs Cache
}

func (r *ResultCache) XOriginalLoad(namespace string, key string, t TypeOfStorage) *evmInt256.Int {
	if t == SStorage {
		return r.OriginalData.Get(namespace, key)
	} else {
		return r.TOriginalData.Get(namespace, key)
	}
}

func (r *ResultCache) XCachedLoad(namespace string, key string, t TypeOfStorage) *evmInt256.Int {
	if t == SStorage {
		return r.CachedData.Get(namespace, key)
	} else {
		return r.TCachedData.Get(namespace, key)
	}
}

func (r *ResultCache) XOriginalStore(namespace string, key string, v *evmInt256.Int, t TypeOfStorage) {
	if t == SStorage {
		r.OriginalData.Set(namespace, key, v)
	} else {
		r.TOriginalData.Set(namespace, key, v)
	}
}

func (r *ResultCache) XCachedStore(namespace string, key string, v *evmInt256.Int, t TypeOfStorage) {
	if t == SStorage {
		r.CachedData.Set(namespace, key, v)
	} else {
		r.TCachedData.Set(namespace, key, v)
	}
}

type CodeCache map[string][]byte

type readOnlyCache struct {
	Code      CodeCache
	CodeSize  Cache
	CodeHash  Cache
	BlockHash Cache
}

func MergeResultCache(src *ResultCache, to *ResultCache) {
	for k, v := range src.OriginalData {
		to.OriginalData[k] = v
	}

	for k, v := range src.CachedData {
		to.CachedData[k] = v
	}

	for k, v := range src.Balance {
		to.Balance[k] = v
	}

	//TODO check whether there are duplicate logs
	for k, v := range src.Logs {
		to.Logs[k] = v
	}

	for k, v := range src.Destructs {
		to.Destructs[k] = v
	}
}
