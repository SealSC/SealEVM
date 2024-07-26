# SealEVM

SealEVM is an independent EVM implementation that decouples from the storage system through interfaces and caching, and can be easily ported to any blockchain system implemented in golang, adding EVM support for it.

---

- [中文](https://github.com/SealSC/SealEVM/blob/master/README_zh.md)

##

### Example

In the [example](https://github.com/SealSC/SealEVM/tree/master/example) directory, a simple usage reference example of SealEVM is provided. This example uses memory as external storage and shows simple functions such as contract deployment, invocation, variable reading, etc.

**⚠️Note: The codes under [example](https://github.com/SealSC/SealEVM/tree/master/example) directory is only for a simple demonstration of code usage, please do not use it in any actual commercial and production environment**

##

### Main structures and interfaces

**⚠️Note, the namespace in the source code is an alias for address, which has the same meaning as address in Ethereum.**

##

>#### Parameters for creating EVM instances
```go
type EVMParam struct {
    MaxStackDepth  int //Maximum stack depth
    ExternalStore  storage.IExternalStorage //External storage interface, which will be explained in detail in later chapters
    ResultCallback EVMResultCallback //Callback function after EVM execution
    Context        *environment.Context //Environment context for EVM execution, please read the source code for the meaning of internal fields
    GasSetting     *instructions.GasSetting //Custom gas cost setting for OpCode
}
```

##

>#### External storage interface
SealEVM will interact with external storage through this interface to implement necessary functions such as state reading, address creation, new contract storage, etc.
```go
type IExternalStorage interface {
    //Get the account balance of the specified address from external storage
    GetBalance(address *evmInt256.Int) (*evmInt256.Int, error)
    
    //Get the contract code of the specified address from external storage
    GetCode(address *evmInt256.Int) ([]byte, error)
    
    //Get the contract code size of the specified address from external storage
    GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error)
    
    //Get the contract code hash of the specified address from external storage
    GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error)
    
    //Get the hash of the specified block from external storage
    GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)
    
    //When executing opcode CREAT(0xF0), this method will be called to get the address of the created contract
    CreateAddress(caller *evmInt256.Int, tx environment.Transaction) *evmInt256.Int
    
    //When executing opcode CREAT2(0xF5), this method will be called to get the address of the created contract
    CreateFixedAddress(caller *evmInt256.Int, salt *evmInt256.Int, tx environment.Transaction) *evmInt256.Int
    
    //Get from external storage whether account balance transfer can be performed
    CanTransfer(from *evmInt256.Int, to *evmInt256.Int, amount *evmInt256.Int) bool
    
    //When executing opcode SLOAD(0x54), get 256-bit data from external storage at the specified location
    //Note: The parameter n is the address of the current executing contract, and the parameter k is the key of the storage location given when executing opcode SLOAD(0x54)
    Load(n string, k string) (*evmInt256.Int, error)

	//When executing opcode TLOAD(0x5C), get 256-bit data from external transient storage at the specified location
    //Note: The parameter n is the address of the current executing contract, and the parameter k is the key of the storage location given when executing opcode TLOAD(0x5C)
    TLoad(n string, k string) (*evmInt256.Int, error)

    //External storage callback after successful execution of CREATE or CREATE2. The address and code of the newly created contract will be provided to external storage through this interface
    //Note: n is the address of the new contract, and code is the bytecode of the new contract
    NewContract(n string, code []byte) error
}
```

##

>#### Execution result
SealEVM will put all data that has changed during contract execution, except for new contract deployment, into cache and will not notify external storage.
```go
type ExecuteResult struct {
    ResultData   []byte //Data returned by contract execution
    GasLeft      uint64 //Remaining gas
    StorageCache storage.ResultCache //Cache of external state changes. External data needs to be updated according to this cache. This will be explained in detail below.
    ExitOpCode   opcodes.OpCode //The last executed opcode when execution is completed
}
```

##

>#### Cache of execution results
The following is a description of the functions of these cache variables. Please refer to the source code for detailed structures.
```go
type ResultCache struct {
    OriginalData CacheUnderNamespace //Cache information read from external storage during execution
    CachedData   CacheUnderNamespace //Cache all state updates at the end of execution

	TOriginalData CacheUnderNamespace //Cache information read from external transient storage during execution
	TCachedData   CacheUnderNamespace //Cache all transient state updates at the end of execution

    Balance   BalanceCache //Cache all account balance changes during execution
    Logs      LogCache //Cache all logs generated by contracts during execution
    Destructs Cache //Cache addresses of contracts that executed opcode SELFDESTRUCT(0xff)
}
```

---

# License

[Apache License 2.0](https://raw.githubusercontent.com/SealSC/SealEVM/master/LICENSE)
