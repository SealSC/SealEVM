package cache

import "github.com/SealSC/SealEVM/types"

type DataBlockCache map[types.Address]types.DataBlock

func (c DataBlockCache) Clone() DataBlockCache {
	replica := make(DataBlockCache)

	for k, v := range c {
		replica[k] = v.Clone()
	}

	return replica
}

func (c DataBlockCache) Merge(cache DataBlockCache) {
	for k, v := range cache {
		c[k] = v
	}
}
