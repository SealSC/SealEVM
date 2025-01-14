package precompiledwithstorage

import (
	"bytes"
	"encoding/hex"
	"math/big"

	"github.com/SealSC/SealEVM/crypto/hashes"
	"github.com/SealSC/SealEVM/evmErrors"
	"github.com/SealSC/SealEVM/storage"
	"github.com/SealSC/SealEVM/types"
)

// CrossTxDataShare is a precompiled contract that supports cross-transaction data sharing
type CrossTxDataShare struct{}

const (
	// Length of function selector
	funcSelectorLength = 4
)

var (
	// Operation codes
	opShare, _ = hex.DecodeString("eb5b7655") // solidity ABI encoded share(bytes32,bytes)
	opRead, _  = hex.DecodeString("61da1439") // solidity ABI encoded read(bytes32)

)

// GasCost calculates gas consumption
func (c *CrossTxDataShare) GasCost(addr types.Address, input []byte, storage storage.IDataBlockStorage) uint64 {
	if len(input) < funcSelectorLength {
		return 0
	}

	switch string(input[:4]) {
	case string(opShare):
		dataSize := len(input[funcSelectorLength:])
		// Base cost is 200 gas
		baseCost := uint64(200)

		// Additional 100 gas per 32-byte data block
		blockCount := (dataSize + 31) / 32
		dataCost := uint64(blockCount) * 100

		return baseCost + dataCost

	case string(opRead):
		// Fixed gas cost of 150 for read operations
		return 150

	default:
		return 0
	}
}

// Execute runs contract logic
func (c *CrossTxDataShare) Execute(caller types.Address, input []byte, dataBlock storage.IDataBlockStorage) ([]byte, error) {
	if len(input) < funcSelectorLength {
		return nil, evmErrors.RevertErr
	}

	switch string(input[:4]) {
	case string(opShare):
		return c.share(caller, input[funcSelectorLength:], dataBlock)
	case string(opRead):
		return c.read(caller, input[funcSelectorLength:], dataBlock)
	default:
		return nil, evmErrors.RevertErr
	}
}

func (c *CrossTxDataShare) getSlotOfCaller(caller types.Address, input []byte) types.Slot {
	slotOfCaller := hashes.Keccak256(bytes.Join([][]byte{caller[:], input[:32]}, []byte{}))
	var slot types.Slot
	copy(slot[:], slotOfCaller[:32])
	return slot
}

// share stores shared data
func (c *CrossTxDataShare) share(caller types.Address, input []byte, dataBlock storage.IDataBlockStorage) ([]byte, error) {
	if len(input) < 64 { // 32(slot) + 32(offset)
		return nil, evmErrors.RevertErr
	}

	// Get slot for caller
	slot := c.getSlotOfCaller(caller, input)

	// Get data from dynamic bytes encoding
	offset := new(big.Int).SetBytes(input[32:64]).Uint64()
	if offset != 64 || len(input) < int(offset+32) {
		return nil, evmErrors.RevertErr
	}

	dataLen := new(big.Int).SetBytes(input[offset : offset+32]).Uint64()
	if len(input) < int(offset+32+dataLen) {
		return nil, evmErrors.RevertErr
	}

	data := make(types.Bytes, dataLen)
	copy(data, input[offset+32:offset+32+dataLen])

	dataBlock.SetDataBlock(slot, data)
	return slot[:], nil
}

// read retrieves shared data
func (c *CrossTxDataShare) read(caller types.Address, input []byte, dataBlock storage.IDataBlockStorage) ([]byte, error) {
	if len(input) < 32 {
		return nil, evmErrors.RevertErr
	}

	// Get slot for caller
	slot := c.getSlotOfCaller(caller, input)

	data, err := dataBlock.GetDataBlock(slot)
	if err != nil {
		return nil, err
	}

	if data == nil {
		// 返回空bytes的ABI编码
		encoded := make([]byte, 64)
		// offset = 32, right aligned
		encoded[31] = 32
		// length = 0, right aligned
		// (already zero)
		return encoded, nil
	}

	// Encode as dynamic bytes
	encoded := make([]byte, 64+len(data))

	// Write offset = 32, right aligned
	encoded[31] = 32

	// Write length, right aligned
	length := big.NewInt(int64(len(data)))
	lengthBytes := length.Bytes()
	copy(encoded[64-len(lengthBytes):64], lengthBytes)

	// Write data
	copy(encoded[64:], data)

	return encoded, nil
}
