package cache

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type ReadOnlyCache struct {
	BlockHash map[types.Slot]*evmInt256.Int
}
