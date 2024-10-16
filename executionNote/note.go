package executionNote

import (
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/opcodes"
	"github.com/SealSC/SealEVM/storage/cache"
	"github.com/SealSC/SealEVM/types"
)

type ExecutionType opcodes.OpCode

const (
	ExternalCall ExecutionType = 0

	Call         = ExecutionType(opcodes.CALL)
	StaticCall   = ExecutionType(opcodes.STATICCALL)
	DelegateCall = ExecutionType(opcodes.DELEGATECALL)
	CallCode     = ExecutionType(opcodes.CALLCODE)
	Create       = ExecutionType(opcodes.CREATE)
	Create2      = ExecutionType(opcodes.CREATE2)
)

var executionTypeNames = map[ExecutionType]string{
	ExternalCall: "ExternalCall",
	Call:         "Call",
	StaticCall:   "StaticCall",
	DelegateCall: "DelegateCall",
	CallCode:     "CallCode",
	Create:       "Create",
	Create2:      "Create2",
}

func (e ExecutionType) String() string {
	return executionTypeNames[e]
}

type NoteConfig struct {
	RecordCache bool
}

type Note struct {
	Type  ExecutionType
	From  types.Address
	To    *types.Address
	Gas   uint64
	Val   *evmInt256.Int
	Input []byte

	ExecutionError error
	ReturnData     []byte
	StorageCache   *cache.ResultCache

	SubNotes []*Note

	config *NoteConfig
}

type MeetNote func(note *Note, depth uint64)

func New(cfg *NoteConfig, execType ExecutionType, tx *environment.Transaction, msg *environment.Message) *Note {
	return &Note{
		Type:   execType,
		From:   msg.Caller,
		To:     tx.To,
		Gas:    tx.GasLimit.Uint64(),
		Val:    msg.Value,
		Input:  msg.Data,
		config: cfg,
	}
}

func (n *Note) SetResult(retData []byte, retErr error, storageCache cache.ResultCache) {
	n.ReturnData = retData
	n.ExecutionError = retErr

	if n.config.RecordCache {
		sCache := storageCache.Clone()
		n.StorageCache = &sCache
	}
}

func (n *Note) GenSubNote(execType ExecutionType, tx *environment.Transaction, msg *environment.Message) *Note {
	newNote := New(n.config, execType, tx, msg)
	n.SubNotes = append(n.SubNotes, newNote)
	return newNote
}

func (n *Note) walk(meetNote MeetNote, depth uint64) {
	if meetNote != nil {
		meetNote(n, depth)
	}
	for _, subNote := range n.SubNotes {
		subNote.walk(meetNote, depth+1)
	}
}

func (n *Note) Walk(meetNote MeetNote) {
	n.walk(meetNote, 0)
}
