package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/SealSC/SealEVM/crypto/hashes"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage/cache"
	"github.com/SealSC/SealEVM/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type extStorage struct {
	Accounts   cache.AccountCache
	DataBlocks map[types.Address]types.DataBlock
}

func newStorage() *extStorage {
	return &extStorage{
		Accounts:   cache.AccountCache{},
		DataBlocks: make(map[types.Address]types.DataBlock),
	}
}

func (r *extStorage) GetAccount(address types.Address) (*environment.Account, error) {
	acc := r.Accounts.Get(address)
	if acc != nil {
		return acc, nil
	}

	return environment.NewAccount(address, nil, nil), nil
}

func (r *extStorage) HashOfCode(code []byte) types.Hash {
	hash := crypto.Keccak256(code)
	var ret types.Hash
	ret.SetBytes(hash)
	return ret
}

func (r *extStorage) GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error) {
	return evmInt256.New(0), nil
}

func (r *extStorage) GetChainID() (*evmInt256.Int, error) {
	return evmInt256.New(0), nil
}

func (r *extStorage) CreateAddress(caller types.Address, tx environment.Transaction) types.Address {
	var ret types.Address
	now := binary.BigEndian.AppendUint64(nil, uint64(time.Now().UnixNano()))
	addr := hashes.Keccak256(now)
	ret.SetBytes(addr)
	return ret
}

func (r *extStorage) CreateFixedAddress(caller types.Address, salt types.Hash, code []byte, tx environment.Transaction) types.Address {
	var ret types.Address
	ret.SetBytes(hashes.Keccak256(salt[:]))
	return ret
}

func (r *extStorage) AccountExist(address types.Address) bool {
	return r.Accounts[address] != nil
}

func (r *extStorage) AccountEmpty(address types.Address) bool {
	return r.Accounts[address] == nil
}

func (r *extStorage) Load(addr types.Address, slot types.Slot) (*evmInt256.Int, error) {
	ret := evmInt256.New(0)
	if r.Accounts[addr] != nil {
		if r.Accounts[addr].Slots[slot] != nil {
			ret.Set(r.Accounts[addr].Slots[slot].Int)
		}
	}

	return ret, nil
}

func (r *extStorage) SetEOA(addr types.Address, balance *evmInt256.Int) *environment.Account {
	r.Accounts[addr] = environment.NewAccount(addr, balance, nil)

	return r.Accounts[addr]
}

func (r *extStorage) SetAccount(acc *environment.Account) *environment.Account {
	r.Accounts[acc.Address] = acc
	return acc
}

func (r *extStorage) StoreResult(ret *cache.ResultCache) {
	for addr, cachedAcc := range ret.CachedAccounts {
		if r.Accounts[addr] == nil {
			r.Accounts[addr] = cachedAcc
		} else {
			r.Accounts[addr].Set(cachedAcc)
		}
	}

	fmt.Println("ret.DataBlockCache: ", ret.DataBlockCache)
	for addr, dataBlock := range ret.DataBlockCache {
		if r.DataBlocks[addr] == nil {
			r.DataBlocks[addr] = make(types.DataBlock)
		}

		for slot, data := range dataBlock {
			fmt.Println("addr: ", addr, "slot: ", slot, "dataBlock: ", dataBlock)

			r.DataBlocks[addr][slot] = data
		}
	}
}

func (r *extStorage) GetDataBlock(address types.Address, slot types.Slot) (types.Bytes, error) {
	if r.DataBlocks[address] == nil {
		r.DataBlocks[address] = make(types.DataBlock)
	}
	return r.DataBlocks[address][slot], nil
}

func (r *extStorage) SetDataBlock(address types.Address, slot types.Slot, data types.Bytes) {
	if r.DataBlocks[address] == nil {
		r.DataBlocks[address] = make(types.DataBlock)
	}
	r.DataBlocks[address][slot] = data
}
