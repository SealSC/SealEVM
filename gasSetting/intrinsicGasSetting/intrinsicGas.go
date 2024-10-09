package intrinsicGasSetting

import "github.com/SealSC/SealEVM/environment"

type IntrinsicGas func(data []byte, to *environment.Contract) uint64

func intrinsicGas(data []byte, to *environment.Contract) uint64 {
	var gasCost uint64 = 21000
	if to == nil {
		gasCost += 32000
	}

	dataLen := uint64(len(data))
	if dataLen > 0 {
		for _, val := range data {
			if val != 0 {
				gasCost += 16
			} else {
				gasCost += 4
			}
		}
	}

	return gasCost
}

func Cost() IntrinsicGas {
	return intrinsicGas
}
