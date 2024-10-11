package cache

import "github.com/SealSC/SealEVM/types"

type DestructCache map[types.Address]types.Address

func (d DestructCache) Clone() DestructCache {
	replica := DestructCache{}

	for k, v := range d {
		replica[k] = v
	}

	return replica
}

func (d DestructCache) Merge(cache DestructCache) {
	for k, v := range cache {
		d[k] = v
	}
}
