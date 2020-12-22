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

package environment

import "github.com/SealSC/SealEVM/evmInt256"

type Block struct {
	Coinbase    *evmInt256.Int
	Timestamp   *evmInt256.Int
	Number      *evmInt256.Int
	Difficulty  *evmInt256.Int
	GasLimit    *evmInt256.Int
	Hash        *evmInt256.Int
}
