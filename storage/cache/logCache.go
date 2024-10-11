package cache

import "github.com/SealSC/SealEVM/types"

type LogCache []*types.Log

func (l LogCache) Clone() LogCache {
	replica := make(LogCache, len(l))

	for i := 0; i < len(l); i++ {
		replica[i] = l[i].Clone()
	}

	return replica
}
