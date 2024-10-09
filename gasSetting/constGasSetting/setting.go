package constGasSetting

import "github.com/SealSC/SealEVM/opcodes"

func Cost() [opcodes.MaxOpCodesCount]uint64 {
	var constCost [opcodes.MaxOpCodesCount]uint64

	//the gas cost of most opcodes is 3
	for idx, _ := range constCost {
		constCost[idx] = 3
	}

	constCost[opcodes.STOP] = 0

	constCost[opcodes.ADDRESS] = 2
	constCost[opcodes.ORIGIN] = 2
	constCost[opcodes.CALLER] = 2
	constCost[opcodes.CALLVALUE] = 2
	constCost[opcodes.CALLDATASIZE] = 2
	constCost[opcodes.CODESIZE] = 2
	constCost[opcodes.GASPRICE] = 2
	constCost[opcodes.RETURNDATASIZE] = 2
	constCost[opcodes.COINBASE] = 2
	constCost[opcodes.TIMESTAMP] = 2
	constCost[opcodes.NUMBER] = 2
	constCost[opcodes.DIFFICULTY] = 2
	constCost[opcodes.GASLIMIT] = 2
	constCost[opcodes.CHAINID] = 2
	constCost[opcodes.BASEFEE] = 2
	constCost[opcodes.BLOBBASEFEE] = 2
	constCost[opcodes.POP] = 2
	constCost[opcodes.PC] = 2
	constCost[opcodes.MSIZE] = 2
	constCost[opcodes.GAS] = 2
	constCost[opcodes.PUSH0] = 2

	constCost[opcodes.JUMPDEST] = 1

	constCost[opcodes.MUL] = 5
	constCost[opcodes.DIV] = 5
	constCost[opcodes.SDIV] = 5
	constCost[opcodes.MOD] = 5
	constCost[opcodes.SMOD] = 5
	constCost[opcodes.SIGNEXTEND] = 5
	constCost[opcodes.SELFBALANCE] = 5

	constCost[opcodes.ADDMOD] = 8
	constCost[opcodes.MULMOD] = 8
	constCost[opcodes.JUMP] = 8

	constCost[opcodes.BLOCKHASH] = 20

	constCost[opcodes.JUMPI] = 10

	constCost[opcodes.TLOAD] = 100
	constCost[opcodes.TSTORE] = 100

	return constCost
}
