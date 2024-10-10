package gasSetting

import (
	"github.com/SealSC/SealEVM/gasSetting/constGasSetting"
	"github.com/SealSC/SealEVM/gasSetting/dynamicGasSetting"
	"github.com/SealSC/SealEVM/gasSetting/intrinsicGasSetting"
)

func defSetting() Setting {
	s := Setting{
		IntrinsicCost:     intrinsicGasSetting.Cost(),
		ConstCost:         constGasSetting.Cost(),
		CommonDynamicCost: dynamicGasSetting.Common(),
		CallCost:          dynamicGasSetting.Call(),
		ContractStoreCost: dynamicGasSetting.ContractStore(),
	}

	return s
}
