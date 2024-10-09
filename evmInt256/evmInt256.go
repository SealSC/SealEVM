/*
 * Copyright 2020 The SealEVM Authors
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package evmInt256

import (
	"github.com/SealSC/SealEVM/utils"
	"math/big"
)

func maxUint(bits uint) *big.Int {
	maxA1 := big.NewInt(0).Lsh(big.NewInt(1), bits)
	return maxA1.Sub(maxA1, big.NewInt(1))
}

func bit(n uint) *big.Int {
	i := big.NewInt(1)
	i.Lsh(i, n)
	return i
}

func pow(x, y int64) *big.Int {
	bx := big.NewInt(x)
	by := big.NewInt(y)
	return bx.Exp(bx, by, nil)
}

func exp256Modulus() *big.Int {
	m := &big.Int{}
	return m.Exp(big.NewInt(2), big.NewInt(256), nil)
}

const (
	maxBits  = 256
	maxBytes = 32
)

var (
	uint256MAX = maxUint(maxBits)
	int256MAX  = pow(2, maxBits-1)
	bit256     = bit(maxBits)
	one        = big.NewInt(1)
	expModulus = exp256Modulus()
)

type Int struct {
	*big.Int
}

func (i *Int) AsStringKey() string {
	return string(utils.LeftPaddingSlice(i.Bytes(), utils.HashLength))
}

func (i Int) Clone() *Int {
	return FromBigInt(i.Int)
}

func New(i uint64) *Int {
	ret := &Int{big.NewInt(0)}
	ret.SetUint64(i)
	return ret
}

func FromBytes(b []byte) *Int {
	ret := &Int{big.NewInt(0)}
	ret.SetBytes(b)
	return ret
}

func FromBigInt(i *big.Int) *Int {
	ret := &Int{big.NewInt(0)}
	if i != nil {
		ret.Set(i)
	}
	return ret
}

func (i *Int) Set(x *big.Int) *Int {
	if x == nil {
		i.Int = big.NewInt(0)
		return i
	}

	i.Int.Set(x)
	return i
}

func (i *Int) SetBytes(b []byte) *Int {
	i.Int.SetBytes(b)
	i.toI256()
	return i
}

func (i *Int) SetUint64(x uint64) *Int {
	i.Int.SetUint64(x)
	return i
}

func (i *Int) toI256() *Int {
	i.Int.And(i.Int, uint256MAX)
	return i
}

func (i *Int) GetSigned() *Int {
	ensure256 := big.NewInt(0).And(i.Int, uint256MAX)

	if ensure256.Cmp(int256MAX) >= 0 {
		ensure256.Sub(ensure256, bit256)
	}

	return &Int{ensure256}
}

func (i *Int) ExtendedAlign(unit uint64) *Int {
	if unit == 0 {
		return i
	}

	v := i.Uint64()
	if v == 0 {
		return i
	}

	org := i.Clone()
	base := New(unit)
	ext := org.Mod(base)

	if ext.IsZero() {
		return i
	}

	ext.SetUint64(unit - ext.Uint64())
	i.Add(ext)
	return i
}

func (i *Int) Add(y *Int) *Int {
	i.Int.Add(i.Int, y.Int)
	return i.toI256()
}

func (i *Int) Mul(y *Int) *Int {
	i.Int.Mul(i.Int, y.Int)
	return i.toI256()
}

func (i *Int) Sub(y *Int) *Int {
	i.Int.Sub(i.Int, y.Int)
	return i.toI256()
}

func (i *Int) Div(y *Int) *Int {
	if y.Sign() == 0 {
		i.Int.SetUint64(0)
		return i
	}

	i.Int.Div(i.Int, y.Int)
	return i.toI256()
}

func (i *Int) SDiv(y *Int) *Int {
	if y.Sign() == 0 {
		i.SetUint64(0)
		return i
	}

	x := i.GetSigned()
	y = y.GetSigned()

	needNeg := x.Sign() != y.Sign()
	i.Int.Div(x.Int.Abs(x.Int), y.Int.Abs(y.Int))

	if needNeg {
		i.Neg(i.Int)
	}
	return i.toI256()
}

func (i *Int) Mod(m *Int) *Int {
	if m.Sign() == 0 || i.Sign() == 0 || m.Int.Cmp(one) == 0 {
		i.SetUint64(0)
		return i
	}

	i.Int.Mod(i.Int, m.Int)
	return i.toI256()
}

func (i *Int) SMod(m *Int) *Int {
	m = m.GetSigned()
	n := i.GetSigned()

	mAbs := big.NewInt(0)
	mAbs.Abs(m.Int)

	if m.Sign() == 0 || n.Sign() == 0 || m.Int.Cmp(one) == 0 {
		i.SetUint64(0)
		return i
	}

	needNeg := n.Sign() < 0
	n.Int.Mod(n.Abs(n.Int), mAbs)
	if needNeg {
		n.Neg(n.Int)
	}

	i.Set(n.Int)
	return i.toI256()
}

func (i *Int) AddMod(y *Int, m *Int) *Int {
	if m.Sign() <= 0 || m.Int.Cmp(one) == 0 {
		i.SetUint64(0)
		return i
	}

	i.Int.Add(i.Int, y.Int)
	return i.Mod(m)
}

func (i *Int) MulMod(y *Int, m *Int) *Int {
	if m.Sign() <= 0 || m.Int.Cmp(one) == 0 {
		i.SetUint64(0)
		return i
	}

	i.Int.Mul(i.Int, y.Int)
	return i.Mod(m)
}

func (i *Int) Exp(e *Int) *Int {
	i.Int.Exp(i.Int, e.Int, expModulus)
	return i.toI256()
}

func (i *Int) SignExtend(baseBytes *Int) *Int {
	if baseBytes.Cmp(big.NewInt(31)) < 0 {
		bits := uint(baseBytes.Uint64()*8 + 7)
		mask := maxUint(bits)

		if i.Int.Bit(int(bits)) > 0 {
			i.Int.Or(i.Int, mask.Not(mask))
		} else {
			i.Int.And(i.Int, mask)
		}
	}

	return i.toI256()
}

func (i *Int) LT(y *Int) bool {
	return i.Int.Cmp(y.Int) < 0
}

func (i *Int) GT(y *Int) bool {
	return i.Int.Cmp(y.Int) > 0
}

func (i *Int) SLT(y *Int) bool {
	si := i.GetSigned()
	sy := y.GetSigned()

	return si.Int.Cmp(sy.Int) < 0
}

func (i *Int) SGT(y *Int) bool {
	si := i.GetSigned()
	sy := y.GetSigned()

	return si.Int.Cmp(sy.Int) > 0
}

func (i *Int) EQ(y *Int) bool {
	return i.Int.Cmp(y.Int) == 0
}

func (i *Int) IsZero() bool {
	return i.Sign() <= 0
}

func (i *Int) And(y *Int) *Int {
	i.Int.And(i.Int, y.Int)
	return i.toI256()
}

func (i *Int) Or(y *Int) *Int {
	i.Int.Or(i.Int, y.Int)
	return i.toI256()
}

func (i *Int) XOr(y *Int) *Int {
	i.Int.Xor(i.Int, y.Int)
	return i.toI256()
}

func (i *Int) Not(y *Int) *Int {
	i.Int.Not(i.Int)
	return i.toI256()
}

func (i *Int) ByteAt(n *Int) byte {
	if n.GT(New(maxBytes - 1)) {
		return 0
	}

	fullBytes := make([]byte, maxBytes, maxBytes)
	bnBytes := i.Int.Bytes()
	bnLen := len(bnBytes)

	if bnLen < maxBytes {
		copy(fullBytes[maxBytes-bnLen:], bnBytes)
	} else {
		copy(fullBytes[:maxBytes], bnBytes[:maxBytes])
	}

	return fullBytes[n.Uint64()]
}

func (i *Int) SHL(n *Int) *Int {
	if n.LT(New(maxBits)) {
		i.Int.Lsh(i.Int, uint(n.Uint64()))
	} else {
		i.Int.SetUint64(0)
	}

	return i.toI256()
}

func (i *Int) SHR(n *Int) *Int {
	if n.LT(New(maxBits)) {
		i.Int.Rsh(i.Int, uint(n.Uint64()))
	} else {
		i.Int.SetUint64(0)
	}

	return i.toI256()
}

func (i *Int) SAR(n *Int) *Int {
	si := i.GetSigned()

	if n.GT(New(maxBits)) {
		if si.Sign() >= 0 {
			si.SetUint64(0)
		} else {
			si.SetInt64(-1)
		}
	} else {
		si.Rsh(si.Int, uint(n.Uint64()))
	}

	i.Set(si.Int)
	return i.toI256()
}
