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

package storageCache

import (
	"SealEVM/evmErrors"
	"SealEVM/evmInt256"
)

type Cache map[string] *evmInt256.Int
type EVMResultCallback func(original Cache, final Cache, err error)

type IExternalStorage interface {
	SLoad(k *evmInt256.Int) (*evmInt256.Int, error)
}

type ICache interface {
	SLoad(k *evmInt256.Int) (*evmInt256.Int, error)
	SStore(k *evmInt256.Int, i *evmInt256.Int)
	ResultFeedback(err error)
}

type storageCache struct {
	originalData    Cache
	cachedData      Cache

	extStorage      IExternalStorage
	callback        EVMResultCallback
}

func New(extStorage IExternalStorage, callback EVMResultCallback) ICache {
	return &storageCache{
		extStorage: extStorage,
		callback: callback,
	}
}

func (s *storageCache) SLoad(k *evmInt256.Int) (*evmInt256.Int, error ) {
	if s.originalData == nil || s.cachedData == nil || s.extStorage == nil {
		return nil, evmErrors.StorageNotInitialized
	}

	i, exists := s.cachedData[k.String()]
	if exists {
		return i, nil
	}

	i, err := s.extStorage.SLoad(k)
	if err != nil {
		return nil, evmErrors.NoSuchDataInTheStorage(err)
	}

	cacheString := i.String()
	s.originalData[cacheString] = evmInt256.FromBigInt(i.Int)
	s.cachedData[cacheString] = i

	return i, nil
}

func (s *storageCache) SStore(k *evmInt256.Int, v *evmInt256.Int)  {
	kString := k.String()
	s.cachedData[kString] = v
}

func (s *storageCache) ResultFeedback(err error) {
	s.callback(s.originalData, s.cachedData, err)
}
