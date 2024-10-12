package environment

import (
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/types"
)

type Account struct {
	Address  types.Address
	Balance  *evmInt256.Int
	Contract *Contract
	Slots    map[types.Slot]*evmInt256.Int
}

func NewAccount(address types.Address, balance *evmInt256.Int, contract *Contract) *Account {
	if balance == nil {
		balance = evmInt256.New(0)
	}
	return &Account{
		Address:  address,
		Balance:  balance,
		Contract: contract,
		Slots:    map[types.Slot]*evmInt256.Int{},
	}
}

func (a Account) Clone() *Account {
	replica := &Account{
		Address: a.Address,
		Slots:   map[types.Slot]*evmInt256.Int{},
	}

	if a.Balance == nil {
		replica.Balance = evmInt256.New(0)
	} else {
		replica.Balance = a.Balance.Clone()
	}

	if a.Contract != nil {
		replica.Contract = a.Contract.Clone()
	}

	for s, v := range a.Slots {
		replica.Slots[s] = v.Clone()
	}

	return replica
}

func (a *Account) Set(newAcc *Account) {
	if newAcc == nil {
		return
	}

	*a = *newAcc
}
