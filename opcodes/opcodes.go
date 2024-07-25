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

package opcodes

type OpCode byte

const (
	MaxOpCodesCount = 256
)

const (
	STOP OpCode = iota + 0x00
	ADD
	MUL
	SUB
	DIV
	SDIV
	MOD
	SMOD
	ADDMOD
	MULMOD
	EXP
	SIGNEXTEND
)

const (
	unusedX0C OpCode = iota + 0x0C
	unusedX0D
	unusedX0E
	unusedX0F
)

const (
	LT OpCode = iota + 0x10
	GT
	SLT
	SGT
	EQ
	ISZERO
)

const (
	AND OpCode = iota + 0x16
	OR
	XOR
	NOT
	BYTE
	SHL
	SHR
	SAR
)

const (
	unusedX1E OpCode = iota + 0x1E
	unusedX1F
)

const (
	SHA3 OpCode = iota + 0x20
)

const (
	unusedX21 OpCode = iota + 0x21
	unusedX22
	unusedX23
	unusedX24
	unusedX25
	unusedX26
	unusedX27
	unusedX28
	unusedX29
	unusedX2A
	unusedX2B
	unusedX2C
	unusedX2D
	unusedX2E
	unusedX2F
)

const (
	ADDRESS OpCode = iota + 0x30
	BALANCE
	ORIGIN
	CALLER
	CALLVALUE
	CALLDATALOAD
	CALLDATASIZE
	CALLDATACOPY
	CODESIZE
	CODECOPY
	GASPRICE
	EXTCODESIZE
	EXTCODECOPY
	RETURNDATASIZE
	RETURNDATACOPY
	EXTCODEHASH
	BLOCKHASH
	COINBASE
	TIMESTAMP
	NUMBER
	DIFFICULTY
	GASLIMIT
)

const (
	CHAINID OpCode = iota + 0x46
	SELFBALANCE
)

const (
	BASEFEE OpCode = iota + 0x48
	BLOBHASH
	BLOBBASEFEE
	unusedX4B
	unusedX4C
	unusedX4D
	unusedX4E
	unusedX4F
)

const (
	POP OpCode = iota + 0x50
	MLOAD
	MSTORE
	MSTORE8
	SLOAD
	SSTORE
	JUMP
	JUMPI
	PC
	MSIZE
	GAS
	JUMPDEST
)

const (
	TLOAD OpCode = iota + 0x5C
	TSTORE
	unusedX5E
)

const (
	PUSH0 OpCode = iota + 0x5F
	PUSH1
	PUSH2
	PUSH3
	PUSH4
	PUSH5
	PUSH6
	PUSH7
	PUSH8
	PUSH9
	PUSH10
	PUSH11
	PUSH12
	PUSH13
	PUSH14
	PUSH15
	PUSH16
	PUSH17
	PUSH18
	PUSH19
	PUSH20
	PUSH21
	PUSH22
	PUSH23
	PUSH24
	PUSH25
	PUSH26
	PUSH27
	PUSH28
	PUSH29
	PUSH30
	PUSH31
	PUSH32
)

const (
	DUP1 OpCode = iota + 0x80
	DUP2
	DUP3
	DUP4
	DUP5
	DUP6
	DUP7
	DUP8
	DUP9
	DUP10
	DUP11
	DUP12
	DUP13
	DUP14
	DUP15
	DUP16
)

const (
	SWAP1 OpCode = iota + 0x90
	SWAP2
	SWAP3
	SWAP4
	SWAP5
	SWAP6
	SWAP7
	SWAP8
	SWAP9
	SWAP10
	SWAP11
	SWAP12
	SWAP13
	SWAP14
	SWAP15
	SWAP16
)

const (
	LOG0 OpCode = iota + 0xA0
	LOG1
	LOG2
	LOG3
	LOG4
)

const (
	unusedXA5 OpCode = iota + 0xA5
	unusedXA6
	unusedXA7
	unusedXA8
	unusedXA9
	unusedXAA
	unusedXAB
	unusedXAC
	unusedXAD
	unusedXAE
	unusedXAF
	unusedXB0 //unofficial PUSH
	unusedXB1 //unofficial DUP
	unusedXB2 //unofficial SWAP
	unusedXB3
	unusedXB4
	unusedXB5
	unusedXB6
	unusedXB7
	unusedXB8
	unusedXB9
	unusedXBA
	unusedXBB
	unusedXBC
	unusedXBD
	unusedXBE
	unusedXBF
	unusedXC0
	unusedXC1
	unusedXC2
	unusedXC3
	unusedXC4
	unusedXC5
	unusedXC6
	unusedXC7
	unusedXC8
	unusedXC9
	unusedXCA
	unusedXCB
	unusedXCC
	unusedXCD
	unusedXCE
	unusedXCF
	unusedXD0
	unusedXD1
	unusedXD2
	unusedXD3
	unusedXD4
	unusedXD5
	unusedXD6
	unusedXD7
	unusedXD8
	unusedXD9
	unusedXDA
	unusedXDB
	unusedXDC
	unusedXDD
	unusedXDE
	unusedXDF
	unusedXE0
	unusedXE1
	unusedXE2
	unusedXE3
	unusedXE4
	unusedXE5
	unusedXE6
	unusedXE7
	unusedXE8
	unusedXE9
	unusedXEA
	unusedXEB
	unusedXEC
	unusedXED
	unusedXEE
	unusedXEF
)

const (
	CREATE OpCode = iota + 0xF0
	CALL
	CALLCODE
	RETURN
	DELEGATECALL
	CREATE2
)

const (
	unusedXF6 OpCode = iota + 0xF6
	unusedXF7
	unusedXF8
	unusedXF9
)

const (
	STATICCALL OpCode = iota + 0xFA
)

const (
	unusedXFB OpCode = iota + 0xFB
	unusedXFC
)

const (
	REVERT OpCode = iota + 0xFD
)

const (
	unusedXFE OpCode = iota + 0xFE
)

const (
	SELFDESTRUCT OpCode = iota + 0xFF
)
