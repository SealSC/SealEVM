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

package instructions

import (
	"github.com/SealSC/SealEVM/types"
	"math/big"

	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/precompiledContracts"
	"github.com/SealSC/SealEVM/utils"
)

func loadEnvironment() {
	instructionTable[opcodes.ADDRESS] = opCodeInstruction{
		action:            addressAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.BALANCE] = opCodeInstruction{
		action:            balanceAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.SELFBALANCE] = opCodeInstruction{
		action:            selfBalanceAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.ORIGIN] = opCodeInstruction{
		action:            originAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.CALLER] = opCodeInstruction{
		action:            callerAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.CALLVALUE] = opCodeInstruction{
		action:            callValueAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.CALLDATALOAD] = opCodeInstruction{
		action:            callDataLoadAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.CALLDATASIZE] = opCodeInstruction{
		action:            callDataSizeAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.CALLDATACOPY] = opCodeInstruction{
		action:            callDataCopyAction,
		requireStackDepth: 3,
		enabled:           true,
	}

	instructionTable[opcodes.CODESIZE] = opCodeInstruction{
		action:            codeSizeAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.CODECOPY] = opCodeInstruction{
		action:            codeCopyAction,
		requireStackDepth: 3,
		enabled:           true,
	}

	instructionTable[opcodes.GASPRICE] = opCodeInstruction{
		action:            gasPriceAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.EXTCODESIZE] = opCodeInstruction{
		action:            extCodeSizeAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.EXTCODECOPY] = opCodeInstruction{
		action:            extCodeCopyAction,
		requireStackDepth: 4,
		enabled:           true,
	}

	instructionTable[opcodes.RETURNDATASIZE] = opCodeInstruction{
		action:            returnDataSizeAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.RETURNDATACOPY] = opCodeInstruction{
		action:            returnDataCopyAction,
		requireStackDepth: 3,
		enabled:           true,
	}

	instructionTable[opcodes.EXTCODEHASH] = opCodeInstruction{
		action:            extCodeHashAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.CHAINID] = opCodeInstruction{
		action:            chainIDAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.BASEFEE] = opCodeInstruction{
		action:            baseFeeAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.BLOCKHASH] = opCodeInstruction{
		action:            blockHashAction,
		requireStackDepth: 1,
		enabled:           true,
	}

	instructionTable[opcodes.COINBASE] = opCodeInstruction{
		action:            coinbaseAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.TIMESTAMP] = opCodeInstruction{
		action:            timestampAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.NUMBER] = opCodeInstruction{
		action:            numberAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.DIFFICULTY] = opCodeInstruction{
		action:            difficultyAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.GASLIMIT] = opCodeInstruction{
		action:            gasLimitAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

	instructionTable[opcodes.GAS] = opCodeInstruction{
		action:            gasAction,
		willIncreaseStack: 1,
		enabled:           true,
	}

}

func addressAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Contract.Address.Int256())
	return nil, nil
}

func balanceAction(ctx *instructionsContext) ([]byte, error) {
	addr := ctx.stack.Peek()
	balance, err := ctx.storage.Balance(types.Int256ToAddress(addr))
	if err != nil {
		return nil, err
	}

	addr.Set(balance.Int)
	return nil, nil
}

func selfBalanceAction(ctx *instructionsContext) ([]byte, error) {
	balance, err := ctx.storage.Balance(ctx.environment.Contract.Address)
	if err != nil {
		return nil, err
	}

	ctx.stack.Push(balance)
	return nil, nil
}

func originAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Transaction.Origin.Int256())
	return nil, nil
}

func callerAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Message.Caller.Int256())
	return nil, nil
}

func callValueAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Message.Value.Clone())
	return nil, nil
}

func callDataLoadAction(ctx *instructionsContext) ([]byte, error) {
	i := ctx.stack.Peek()
	data := utils.GetDataFrom(ctx.environment.Message.Data, i.Uint64(), 32)

	i.SetBytes(data)
	return nil, nil
}

func callDataSizeAction(ctx *instructionsContext) ([]byte, error) {
	i := evmInt256.New(0)
	s := ctx.environment.Message.DataSize()

	i.SetUint64(s)
	ctx.stack.Push(i)
	return nil, nil
}

func callDataCopyAction(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	dOffset := ctx.stack.Pop()
	size := ctx.stack.Pop()

	//gas check
	offset, _, _, err := ctx.memoryGasCostAndMalloc(mOffset, size)
	if err != nil {
		return nil, err
	}

	data := utils.GetDataFrom(ctx.environment.Message.Data, dOffset.Uint64(), size.Uint64())
	err = ctx.memory.Store(offset, data)
	return nil, err
}

func codeSizeAction(ctx *instructionsContext) ([]byte, error) {
	s := evmInt256.New(int64(len(ctx.environment.Contract.Code)))
	ctx.stack.Push(s)
	return nil, nil
}

func codeCopyAction(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	dOffset := ctx.stack.Pop()
	size := ctx.stack.Pop()

	//gas check
	offset, _, _, err := ctx.memoryGasCostAndMalloc(mOffset, size)
	if err != nil {
		return nil, err
	}

	data := utils.GetDataFrom(ctx.environment.Contract.Code, dOffset.Uint64(), size.Uint64())
	err = ctx.memory.Store(offset, data)
	return nil, err
}

func gasPriceAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Transaction.GasPrice.Clone())
	return nil, nil
}

func extCodeSizeAction(ctx *instructionsContext) ([]byte, error) {
	addrInt := ctx.stack.Peek()

	addr := types.Int256ToAddress(addrInt)
	if precompiledContracts.IsPrecompiledContract(addr) {
		addrInt.SetUint64(0)
		return nil, nil
	}

	s, err := ctx.storage.GetCodeSize(addr)
	if err != nil {
		return nil, err
	}

	if s == nil {
		addrInt.SetUint64(0)
	} else {
		addrInt.Set(s.Int)
	}
	return nil, nil
}

func extCodeCopyAction(ctx *instructionsContext) ([]byte, error) {
	addrInt := ctx.stack.Pop()
	mOffset := ctx.stack.Pop()
	dOffset := ctx.stack.Pop()
	size := ctx.stack.Pop()

	addr := types.Int256ToAddress(addrInt)
	if precompiledContracts.IsPrecompiledContract(addr) {
		return nil, nil
	}

	code, err := ctx.storage.GetCode(addr)
	if err != nil {
		return nil, err
	}

	//gas check
	offset, _, _, err := ctx.memoryGasCostAndMalloc(mOffset, size)
	if err != nil {
		return nil, err
	}

	data := utils.GetDataFrom(code, dOffset.Uint64(), size.Uint64())
	err = ctx.memory.Store(offset, data)
	return nil, err
}

func returnDataSizeAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(evmInt256.New(int64(len(ctx.lastReturn))))
	return nil, nil
}

func returnDataCopyAction(ctx *instructionsContext) ([]byte, error) {
	mOffset := ctx.stack.Pop()
	dOffset := ctx.stack.Pop()
	dLen := ctx.stack.Pop()

	end := big.NewInt(0).Add(dOffset.Int, dLen.Int)
	if !end.IsUint64() || end.Uint64() > uint64(len(ctx.lastReturn)) {
		return nil, evmErrors.ReturnDataCopyOutOfBounds
	}

	//gas check
	offset, size, _, err := ctx.memoryGasCostAndMalloc(mOffset, dLen)
	if err != nil {
		return nil, err
	}

	err = ctx.memory.Store(offset, ctx.lastReturn[dOffset.Uint64():size])
	return nil, err
}

func extCodeHashAction(ctx *instructionsContext) ([]byte, error) {
	addrInt := ctx.stack.Peek()

	addr := types.Int256ToAddress(addrInt)
	if precompiledContracts.IsPrecompiledContract(addr) {
		addrInt.SetBytes(utils.ZeroHash)
		return nil, nil
	}

	codeHash, err := ctx.storage.GetCodeHash(addr)
	if err != nil {
		return nil, err
	}

	if codeHash == nil {
		addrInt.SetUint64(0)
	} else {
		addrInt.SetBytes(codeHash[:])
	}
	return nil, nil
}

func chainIDAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Block.ChainID.Clone())
	return nil, nil
}

func baseFeeAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Block.BaseFee.Clone())
	return nil, nil
}

func blockHashAction(ctx *instructionsContext) ([]byte, error) {
	blk := ctx.stack.Peek()
	blkHash, err := ctx.storage.GetBlockHash(blk)
	if err != nil {
		return nil, err
	}

	blk.Set(blkHash.Int)
	return nil, nil
}

func coinbaseAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Block.Coinbase.Clone())
	return nil, nil
}

func timestampAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Block.Timestamp.Clone())
	return nil, nil
}

func numberAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Block.Number.Clone())
	return nil, nil
}

func difficultyAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Block.Difficulty.Clone())
	return nil, nil
}

func gasLimitAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.environment.Block.GasLimit.Clone())
	return nil, nil
}

func gasAction(ctx *instructionsContext) ([]byte, error) {
	ctx.stack.Push(ctx.gasRemaining.Clone())
	return nil, nil
}
