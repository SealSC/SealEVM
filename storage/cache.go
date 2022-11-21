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

	Balance   BalanceCache
	Logs      LogCache
	Destructs Cache
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

func DupResultCache(src *ResultCache) ResultCache {
	dst := ResultCache{}

	dst.OriginalData = make(CacheUnderNamespace)
	for k, v := range src.OriginalData {
		oriData := make(Cache)
		for mk, mv := range v {
			oriData[mk] = mv.Clone()
		}
		dst.OriginalData[k] = oriData
	}

	dst.CachedData = make(CacheUnderNamespace)
	for k, v := range src.CachedData {
		//TODO v is a map
		cachedData := make(Cache)
		for mk, mv := range v {
			cachedData[mk] = mv.Clone()
		}
		dst.CachedData[k] = cachedData
	}

	dst.Balance = make(BalanceCache)
	for k, v := range src.Balance {
		dst.Balance[k] = &balance{
			Address: v.Address.Clone(),
			Balance: v.Balance.Clone(),
		}
	}

	dst.Logs = make(LogCache)
	for k, v := range src.Logs {
		logSlice := make([]Log, len(v))
		copy(logSlice, v)
		for sk, sv := range v {
			//Topics
			logSlice[sk].Topics = make([][]byte, len(sv.Topics))
			copy(logSlice[sk].Topics, sv.Topics)
			for tsk, tsv := range sv.Topics {
				logSlice[sk].Topics[tsk] = make([]byte, len(tsv))
				copy(logSlice[sk].Topics[tsk], tsv)
			}
			//Data
			logSlice[sk].Data = make([]byte, len(sv.Data))
			copy(logSlice[sk].Data, sv.Data)

			//Context
			logSlice[sk].Context = sv.Context
		}
		dst.Logs[k] = logSlice
	}

	dst.Destructs = make(Cache)
	for k, v := range src.Destructs {
		dst.Destructs[k] = v.Clone()
	}

	return dst
}
