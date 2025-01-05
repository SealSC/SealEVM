package types

type DataBlock map[Slot]Bytes

func (d DataBlock) Clone() DataBlock {
	replica := make(DataBlock)
	for k, v := range d {
		replica[k] = v.Clone()
	}
	return replica
}
