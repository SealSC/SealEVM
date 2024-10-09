package dynamicGasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

type CommonCalculator func(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (memExpSize uint64, gasCost uint64, err error)

func Common() [opcodes.MaxOpCodesCount]CommonCalculator {
	var commDynamicCost [opcodes.MaxOpCodesCount]CommonCalculator

	commDynamicCost[opcodes.EXP] = gasOfExp
	commDynamicCost[opcodes.SHA3] = gasOfKeccak

	commDynamicCost[opcodes.BALANCE] = gasOfBalance
	commDynamicCost[opcodes.EXTCODESIZE] = gasOfExtCodeSize
	commDynamicCost[opcodes.EXTCODEHASH] = gasOfExtCodeHash

	commDynamicCost[opcodes.CALLDATACOPY] = gasOfCopy
	commDynamicCost[opcodes.CODECOPY] = gasOfCopy
	commDynamicCost[opcodes.RETURNDATACOPY] = gasOfCopy

	commDynamicCost[opcodes.EXTCODECOPY] = gasOfExtCodeCopy

	commDynamicCost[opcodes.MLOAD] = gasOfMemory(evmInt256.New(32))
	commDynamicCost[opcodes.MSTORE] = gasOfMemory(evmInt256.New(32))
	commDynamicCost[opcodes.MSTORE8] = gasOfMemory(evmInt256.New(0))

	commDynamicCost[opcodes.MCOPY] = gasOfCopy

	commDynamicCost[opcodes.SLOAD] = gasOfSLoad

	commDynamicCost[opcodes.LOG0] = gasOfLog(0)
	commDynamicCost[opcodes.LOG1] = gasOfLog(1)
	commDynamicCost[opcodes.LOG2] = gasOfLog(2)
	commDynamicCost[opcodes.LOG3] = gasOfLog(3)
	commDynamicCost[opcodes.LOG4] = gasOfLog(4)

	commDynamicCost[opcodes.CREATE] = gasOfCreate(false)
	commDynamicCost[opcodes.CREATE2] = gasOfCreate(true)

	commDynamicCost[opcodes.RETURN] = gasOfMemory(nil)
	commDynamicCost[opcodes.REVERT] = gasOfMemory(nil)
	commDynamicCost[opcodes.SELFDESTRUCT] = gasOfSelfDestruct

	return commDynamicCost
}

func Call() [opcodes.MaxOpCodesCount]CallGas {
	var callCost [opcodes.MaxOpCodesCount]CallGas

	callCost[opcodes.CALL] = gasOfCall
	callCost[opcodes.CALLCODE] = gasOfCall
	callCost[opcodes.STATICCALL] = gasOfCall
	callCost[opcodes.DELEGATECALL] = gasOfCall

	return callCost
}

func SStore() [opcodes.MaxOpCodesCount]SStoreGas {
	var storeCost [opcodes.MaxOpCodesCount]SStoreGas
	storeCost[opcodes.SSTORE] = gasOfSStore
	return storeCost
}

func ContractStore() ContractStoreGas {
	return gasOfContractStore
}
