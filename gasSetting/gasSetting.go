package gasSetting

import (
	"github.com/SealSC/SealEVM/gasSetting/dynamicGasSetting"
	"github.com/SealSC/SealEVM/gasSetting/intrinsicGasSetting"
	"github.com/SealSC/SealEVM/opcodes"
)

type Setting struct {
	IntrinsicCost     intrinsicGasSetting.IntrinsicGas
	ConstCost         [opcodes.MaxOpCodesCount]uint64
	CommonDynamicCost [opcodes.MaxOpCodesCount]dynamicGasSetting.CommonCalculator
	CallCost          [opcodes.MaxOpCodesCount]dynamicGasSetting.CallGas
	ContractStoreCost dynamicGasSetting.ContractStoreGas
}

var setting = defSetting()

func Set(s *Setting) {
	setting = *s
}

func Get() *Setting {
	return &setting
}
