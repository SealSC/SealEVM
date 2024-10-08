package gasSetting

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
)

func setDynamicGasCalculators(s *Setting) {
	s.DynamicCost[opcodes.EXP] = gasOfExp
	s.DynamicCost[opcodes.SHA3] = gasOfKeccak

	s.DynamicCost[opcodes.BALANCE] = gasOfBalance
	s.DynamicCost[opcodes.EXTCODESIZE] = gasOfExtCodeSize
	s.DynamicCost[opcodes.EXTCODEHASH] = gasOfExtCodeHash

	s.DynamicCost[opcodes.CALLDATACOPY] = gasOfCopy
	s.DynamicCost[opcodes.CODECOPY] = gasOfCopy
	s.DynamicCost[opcodes.RETURNDATACOPY] = gasOfCopy

	s.DynamicCost[opcodes.EXTCODECOPY] = gasOfExtCodeCopy

	s.DynamicCost[opcodes.MLOAD] = gasOfMemory(evmInt256.New(32))
	s.DynamicCost[opcodes.MSTORE] = gasOfMemory(evmInt256.New(32))
	s.DynamicCost[opcodes.MSTORE8] = gasOfMemory(evmInt256.New(0))

	s.DynamicCost[opcodes.MCOPY] = gasOfCopy

	s.DynamicCost[opcodes.SLOAD] = gasOfSLoad
	s.ConstCost[opcodes.SSTORE] = 0

	s.DynamicCost[opcodes.LOG0] = gasOfLog(0)
	s.DynamicCost[opcodes.LOG1] = gasOfLog(1)
	s.DynamicCost[opcodes.LOG2] = gasOfLog(2)
	s.DynamicCost[opcodes.LOG3] = gasOfLog(3)
	s.DynamicCost[opcodes.LOG4] = gasOfLog(4)

	s.DynamicCost[opcodes.CREATE] = gasOfCreate(false)
	s.DynamicCost[opcodes.CREATE2] = gasOfCreate(true)

	s.ConstCost[opcodes.CALL] = 0
	s.ConstCost[opcodes.CALLCODE] = 0
	s.ConstCost[opcodes.DELEGATECALL] = 0
	s.ConstCost[opcodes.STATICCALL] = 0

	s.DynamicCost[opcodes.RETURN] = gasOfMemory(nil)
	s.DynamicCost[opcodes.REVERT] = gasOfMemory(nil)
	s.DynamicCost[opcodes.SELFDESTRUCT] = gasOfSelfDestruct
}

func defSetting() Setting {
	s := Setting{}

	//the gas cost of most opcodes is 3
	for idx, _ := range s.ConstCost {
		s.ConstCost[idx] = 3
	}

	s.ConstCost[opcodes.STOP] = 0

	s.ConstCost[opcodes.ADDRESS] = 2
	s.ConstCost[opcodes.ORIGIN] = 2
	s.ConstCost[opcodes.CALLER] = 2
	s.ConstCost[opcodes.CALLVALUE] = 2
	s.ConstCost[opcodes.CALLDATASIZE] = 2
	s.ConstCost[opcodes.CODESIZE] = 2
	s.ConstCost[opcodes.GASPRICE] = 2
	s.ConstCost[opcodes.RETURNDATASIZE] = 2
	s.ConstCost[opcodes.COINBASE] = 2
	s.ConstCost[opcodes.TIMESTAMP] = 2
	s.ConstCost[opcodes.NUMBER] = 2
	s.ConstCost[opcodes.DIFFICULTY] = 2
	s.ConstCost[opcodes.GASLIMIT] = 2
	s.ConstCost[opcodes.CHAINID] = 2
	s.ConstCost[opcodes.BASEFEE] = 2
	s.ConstCost[opcodes.BLOBBASEFEE] = 2
	s.ConstCost[opcodes.POP] = 2
	s.ConstCost[opcodes.PC] = 2
	s.ConstCost[opcodes.MSIZE] = 2
	s.ConstCost[opcodes.GAS] = 2
	s.ConstCost[opcodes.PUSH0] = 2

	s.ConstCost[opcodes.JUMPDEST] = 1

	s.ConstCost[opcodes.MUL] = 5
	s.ConstCost[opcodes.DIV] = 5
	s.ConstCost[opcodes.SDIV] = 5
	s.ConstCost[opcodes.MOD] = 5
	s.ConstCost[opcodes.SMOD] = 5
	s.ConstCost[opcodes.SIGNEXTEND] = 5
	s.ConstCost[opcodes.SELFBALANCE] = 5

	s.ConstCost[opcodes.ADDMOD] = 8
	s.ConstCost[opcodes.MULMOD] = 8
	s.ConstCost[opcodes.JUMP] = 8

	s.ConstCost[opcodes.BLOCKHASH] = 20

	s.ConstCost[opcodes.JUMPI] = 10

	s.ConstCost[opcodes.TLOAD] = 100
	s.ConstCost[opcodes.TSTORE] = 100

	setDynamicGasCalculators(&s)
	return s
}
