package gasSetting

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/gasSetting/dynamicGasSetting"
	"github.com/SealSC/SealEVM/gasSetting/intrinsicGasSetting"
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

type Setting struct {
	IntrinsicCost     intrinsicGasSetting.IntrinsicGas
	ConstCost         [opcodes.MaxOpCodesCount]uint64
	CommonDynamicCost [opcodes.MaxOpCodesCount]dynamicGasSetting.CommonCalculator
	CallCost          [opcodes.MaxOpCodesCount]dynamicGasSetting.CallGas
	SStoreCost        [opcodes.MaxOpCodesCount]dynamicGasSetting.SStoreGas
	ContractStoreCost dynamicGasSetting.ContractStoreGas
}

var setting = defSetting()

func Set(s *Setting) {
	setting = *s
}

func Get() *Setting {
	return &setting
}
