package gasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/memory"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/stack"
	"github.com/SealSC/SealEVM/storage"
)

type DynamicGasCalculator func(
	contract *environment.Contract,
	stx *stack.Stack,
	mem *memory.Memory,
	store *storage.Storage,
) (memExpSize uint64, gasCost uint64, err error)

type ResultGasCalculator func(stack *stack.Stack, mem *memory.Memory, store *storage.Storage) (uint64, error)

type Setting struct {
	ConstCost   [opcodes.MaxOpCodesCount]uint64
	DynamicCost [opcodes.MaxOpCodesCount]DynamicGasCalculator
	ResultCost  [opcodes.MaxOpCodesCount]ResultGasCalculator
}

var setting = defSetting()

func Set(s *Setting) {
	setting = *s
}

func Get() Setting {
	return setting
}
