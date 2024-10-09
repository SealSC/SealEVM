package types

import "github.com/SealSC/SealEVM/evmInt256"

type Topic = Hash

func Int256ToTopic(i *evmInt256.Int) Topic {
	var s Topic
	s.SetBytes(i.Bytes())
	return s
}

type Log struct {
	Address Address
	Topics  []Topic
	Data    []byte
}

func (l Log) Clone() *Log {
	replica := &Log{}

	replica.Topics = make([]Topic, len(l.Topics))
	for idx, t := range l.Topics {
		replica.Topics[idx] = t
	}

	d := make([]byte, len(l.Data))
	copy(d, l.Data)
	replica.Data = d

	replica.Address = l.Address
	return replica
}