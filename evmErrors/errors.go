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

package evmErrors

import (
	"errors"
	"fmt"
	"github.com/SealSC/SealEVM/types"
)

var StackUnderFlow = errors.New("stack under flow")
var StackOverFlow = errors.New("stack over flow")
var ClosureDepthOverflow = errors.New("call/create depth overflow (>1024)")
var StorageNotInitialized = errors.New("storage not initialized")
var InvalidEVMInstance = errors.New("invalid EVM instance")
var ReturnDataCopyOutOfBounds = errors.New("return data copy out of bounds")
var JumpOutOfBounds = errors.New("jump out of range")
var InvalidJumpDest = errors.New("invalid jump dest")
var JumpToNoneOpCode = errors.New("jump to non-OpCode")
var OutOfGas = errors.New("out of gas")
var InsufficientBalance = errors.New("insufficient balance")
var WriteProtection = errors.New("write protection")
var RevertErr = errors.New("revert")
var BN256BadPairingInput = errors.New("bn256 bad pairing input")
var InvalidExternalStorageResult = errors.New("external storage return invalid values")

func Panicked(err error) error {
	return errors.New("panic error: " + err.Error())
}

func InvalidOpCode(code byte) error {
	return errors.New(fmt.Sprintf("invalid op code: 0x%X", code))
}

func NoSuchDataInTheStorage(err error) error {
	return errors.New("no such data in the storage: " + err.Error())
}

func InvalidTypeOfStorage() error {
	return errors.New("invalid type of storage for reading or writing")
}

func InvalidPrecompiledAddress(addr types.Address) error {
	return errors.New("invalid precompiled contract address " + addr.String())
}

func DuplicatePrecompiledContract(addr types.Address) error {
	return errors.New("duplicate precompiled contract address " + addr.String())
}

func NoExternalDataBlockStorageSet() error {
	return errors.New("no external data block storage set")
}

var OutOfMemory = errors.New("out of memory")
