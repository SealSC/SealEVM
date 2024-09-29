package types

import (
	"github.com/SealSC/SealEVM/evmInt256"
)

type Slot = Hash

func Int256ToSlot(i *evmInt256.Int) Slot {
	var s Slot
	s.SetBytes(i.Bytes())
	return s
}
