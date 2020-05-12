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
	"SealEVM/environment"
	"SealEVM/evmInt256"
)

//todo: improve cache struct, because they has namespace-liked prefix.
type Cache map[string] *evmInt256.Int

type balance struct {
	Address *evmInt256.Int
	Balance *evmInt256.Int
}

type BalanceCache map[string] *balance

type log struct {
	Topics  [][]byte
	Data    []byte
	Context environment.Context
}

type LogCache map[string] []log

type ResultCache struct {
	OriginalData    Cache
	CachedData      Cache

	Balance         BalanceCache
	Logs            LogCache
	Destructs       Cache
}

type CodeCache map[string] []byte

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
		if to.Balance[k] != nil {
			to.Balance[k].Balance.Add(v.Balance)
		} else {
			to.Balance[k] = v
		}
	}

	for k, v := range src.Logs {
		to.Logs[k] = append(to.Logs[k], v...)
	}

	for k, v := range src.Destructs {
		to.Destructs[k] = v
	}
}

