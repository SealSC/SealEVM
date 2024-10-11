# SealEVM

SealEVM is a standalone EVM executor, aiming to create a completely decoupled EVM execution environment from the storage system, to add EVM support for any blockchain system. 
The current version has achieved decoupling from the storage system through interfaces and caching, supporting the addition of EVM support to any blockchain system implemented in Golang.

**[中文](https://github.com/SealSC/SealEVM/blob/master/README_zh.md) | English**

##

- [Usage](#usage)
- [Main Structures and Interfaces](#main-structures-and-interfaces)
  - [EVM Instance Configuration Parameters](#evm-instance-configuration-parameters)
  - [External Storage Interface](#external-storage-interface)
  - [Execution Environment Struct](#execution-environment-struct)
  - [Execution Results Struct](#execution-results-struct)
  - [Execution Result Cache Structure](#execution-result-cache-structure)
- [Gas Setting](#gas-setting)
- [Precompiled Contracts](#precompiled-contracts)
- [Precompiled Contracts](#precompiled-contracts)
- [Usage Scenarios](#usage-scenarios)
  - [User Case](#user-case)

## Usage

```go
func main() {
    // Load the SealEVM module. SealEVM does not use Golang's init feature and must be explicitly executed globally before use.
    SealEVM.Load()
    
    // Configure the SealEVM instance parameters as needed. 
    // Transaction input and other parameters are provided in this struct. 
    // The structure of this parameter is explained in the following sections.
    evmParam := SealEVM.EVMParam{}
    
    // Use this parameter to initialize a SealEVM instance.
    evm := SealEVM.New(evmParam)
    
    // Execute the EVM, the result returned is an ExecuteResult struct.
    // This struct stores information about the original state of the data, the final state, contract logs, internally created contracts, etc. 
    // The structure of this struct is explained in the following sections.
    result, err := evm.ExecuteContract()
}
```
The [**example**](https://github.com/SealSC/SealEVM/tree/master/example) directory provides a simple reference example of SealEVM usage. 
This example uses memory as external storage, demonstrating simple functions like contract deployment, invocation, and variable reading.

## Main Structures and Interfaces

>#### EVM Instance Configuration Parameters
```go
// Definition of the execution result callback function. 
// The result and err parameters in the interface are the same as those returned by the EVM execution.
type EVMResultCallback func(result ExecuteResult, err error)

type EVMParam struct {
    MaxStackDepth  int // Maximum execution stack depth (note: not storage stack depth)
    ExternalStore  storage.IExternalStorage // External storage interface, explained in the following sections
    ResultCallback EVMResultCallback // Callback function after EVM execution, defined at the beginning of this code block
    Context        *environment.Context // Context structure during EVM execution, explained in the following sections
    GasSetting     *gasSetting.Setting // Gas fee settings, use default if nil, explained in the following sections
}

```

##

>#### External Storage Interface
SealEVM interacts with external storage through this interface to achieve necessary 
contract reading, state reading, address creation, and other functions.

```go
type IExternalStorage interface {
    // Retrieve stored contracts
    GetContract(address types.Address) (*environment.Contract, error)
    
    // Get block hash at a specified height
    GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)
    
    // Check if a contract exists
    ContractExist(address types.Address) bool
    
    // Check if an address is empty (see EIP161 for the definition of empty)
    ContractEmpty(address types.Address) bool
    
    // Return the hash value of the given contract code
    HashOfCode(code []byte) types.Hash
    
    // Return the created contract address based on parameters, used by opcode CREATE (0xF0)
    CreateAddress(caller types.Address, tx environment.Transaction) types.Address
    
    // Return the created contract address based on parameters, used by opcode CREATE2 (0xF5)
    CreateFixedAddress(caller types.Address, salt types.Hash, code []byte, tx environment.Transaction) types.Address
    
    // Retrieve 256-bit data from external storage at a specified slot during the execution of opcode SLOAD (0x54)
    Load(address types.Address, slot types.Slot) (*evmInt256.Int, error)
}
```

##

>#### Execution Environment Struct
Execution Environment Structure This structure is in the [environment](./environment) package
and represents the execution context during SealEVM execution,
including parameters like block, transaction, message, and contract.

```go
// Execution environment summary structure. When Contract is nil, it is considered a contract creation transaction.
// When creating a contract, the Message.Data field in this structure is used as the contract code, 
// and a contract instance is created and assigned to the Contract field.
type Context struct {
    Block       Block       // Block environment structure, details are explained below
    Contract    *Contract   // The contract structure being executed, details are explained below
    Transaction Transaction // Transaction structure, details are explained below
    Message     Message     // Message structure, details are explained below
}

// Block environment structure
type Block struct {
    ChainID     *evmInt256.Int // Block ID
    Coinbase    types.Address  // Block miner's address
    Timestamp   uint64         // Block's timestamp in seconds
    Number      uint64         // Block height, using a 64-bit unsigned integer
    Difficulty  *evmInt256.Int // Difficulty
    GasLimit    *evmInt256.Int // Block gas limit
    Hash        types.Hash     // Block hash
    BaseFee     *evmInt256.Int // Base fee as defined in EIP-1559
    BlobBaseFee *evmInt256.Int // Blob base fee as defined in EIP-7516
}

// Contract environment structure
type Contract struct {
    Address  types.Address  // Address of the contract being executed
    Code     []byte         // Contract code
    CodeHash types.Hash     // Hash of the contract code
    CodeSize uint64         // Byte size of the contract code
    Balance  *evmInt256.Int // Account balance of the contract address
}

// Transaction environment structure
type Transaction struct {
    TxHash    types.Hash     // Transaction hash
    Origin    types.Address  // Address that initiated this transaction, value obtained using the ORIGIN (0x32) opcode
    GasPrice  *evmInt256.Int // Gas price of the transaction
    GasLimit  *evmInt256.Int // Gas limit of the transaction
    BlobHashes []types.Hash  // tx.blob_versioned_hashes in EIP-4844
}

// Message structure
type Message struct {
    Caller types.Address  // Address of the contract caller, value obtained using the CALLER (0x33) opcode
    Value  *evmInt256.Int // Amount of ETH sent during the contract call
    Data   []byte         // Parameters passed to the contract during the call
}
```

##

>#### Execution Results Struct
When executing contracts, SealEVM places all data changes, except for new contract deployments, into a cache and does not notify external storage.
```go
type ExecuteResult struct {
    ResultData   []byte //Data returned by contract execution
    GasLeft      uint64 //Remaining gas
    StorageCache storage.ResultCache //Cache of external state changes. External data needs to be updated according to this cache. This will be explained in detail below.
    ExitOpCode   opcodes.OpCode //The last executed opcode when execution is completed
}
```

##

>#### Execution Result Cache Structure
SealEVM will store the original data obtained from external sources during and after execution, 
as well as the final result data, into a unified cache structure instance. 
After the contract execution is completed, this instance is returned to the caller as the result. 
This design allows external storage to focus only on providing the initial data 
without worrying about state changes during execution, and still obtain the result after execution.

```go
// Slot cache type, stores data of specified slots under the account, 
// corresponding to the storage space of SSTORE and SLOAD opcodes
type SlotCache map[types.Slot]*evmInt256.Int

// Account cache unit structure
type AccountCacheUnit struct {
    Contract *environment.Contract // Contract information under the account, see the contract environment structure section for field descriptions
    Slots    SlotCache             // Slot cache under the account, corresponding to the storage space of SSTORE and SLOAD opcodes
}

// Account cache type, contracts and Slot values loaded from external storage, 
// summarized and managed through the account cache
type AccountCache map[types.Address]*AccountCacheUnit

// Log cache type, sequentially stores Log data generated during contract execution
type LogCache []*types.Log

// Destructed contract address cache type, stores addresses of contracts that executed SELFDESTRUCT (0xFF)
type DestructCache map[types.Address]types.Address

// Contract cache type, stores instances of contracts created by internal transactions
type ContractCache map[types.Address]*environment.Contract

// Summary cache structure
type ResultCache struct {
    OriginalAccounts AccountCache // Original state cache of accounts loaded from external storage during execution
    CachedAccounts   AccountCache // Final state cache of accounts after execution

    Logs         *LogCache // Log cache generated by opcodes LOG0 (0xA0) ~ LOG4 (0xA4)
    Destructs    DestructCache // Cache for contracts that executed SELFDESTRUCT (0xFF)
    NewContracts ContractCache // Cache for contracts created by internal transactions during execution
}
```

## Gas Setting
SealEVM achieves flexible Gas settings through the [gasSettings](./gasSetting) package and provides a 
default settings instance that aligns as closely as possible with the Ethereum Gas system.

```go
// Definition of the intrinsic Gas fee calculation function type for transactions. 
// 'data' is the input data for the transaction, i.e., parameters. 
// 'to' is a pointer to the contract structure; if 'to' is nil, it indicates that this transaction is a contract creation transaction.
// The return value is the intrinsic Gas consumption for the transaction (gasCost).
type intrinsicGasSetting.IntrinsicGas func(data []byte, to *environment.Contract) (gasCost uint64)

// General dynamic Gas consumption calculation function type definition.
// Needs to return the memory expansion size (memExpSize) and Gas consumption (gasCost).
type dynamicGasSetting.CommonCalculator func(
    contract *environment.Contract, // Contract environment variable for opcode execution
    stx *stack.Stack,               // Stack environment at the start of opcode execution
    mem *memory.Memory,             // Memory environment at the start of opcode execution
    store *storage.Storage,         // Storage environment at the start of opcode execution
) (memExpSize uint64, gasCost uint64, err error)

// Gas consumption calculation function type definition for CALL, CALLCODE, STATICCALL, DELEGATECALL.
// Needs to return the memory expansion size (memExpSize), Gas consumption (gasCost), 
// and the amount of gas sent for the internal call (sendGas).
type CallGas func(
    code opcodes.OpCode,     // Current opcode
    availableGas uint64,     // Currently available Gas amount
    stx *stack.Stack,        // Stack environment at the start of opcode execution
    mem *memory.Memory,      // Memory environment at the start of opcode execution
    store *storage.Storage,  // Storage environment at the start of opcode execution
) (expSize uint64, gasCost uint64, sendGas uint64, err error)

// Gas consumption calculation function type definition for storing new contracts.
// Needs to return the gas consumption (gasCost).
type ContractStoreGas func(code []byte, gasRemaining uint64) (gasCost uint64, err error)

type Setting struct {
    // Fixed Gas fee calculation function for each transaction, 
    // called once before execution begins to deduct fixed transaction fees.
    IntrinsicCost intrinsicGasSetting.IntrinsicGas
    
    // Dynamic Gas consumption calculation function configuration. 
    // CALL, CALLCODE, STATICCALL, DELEGATECALL, SSTORE opcodes will ignore this configuration.
    CommonDynamicCost [opcodes.MaxOpCodesCount]dynamicGasSetting.CommonCalculator
    
    // Gas calculation configuration for internal calls to CALL, CALLCODE, STATICCALL, DELEGATECALL.
    CallCost [opcodes.MaxOpCodesCount]dynamicGasSetting.CallGas
    
    // Gas calculation configuration for storing contract code when creating contracts 
    // through Create, Create2, and contract creation transactions.
    ContractStoreCost dynamicGasSetting.ContractStoreGas
    
    // Fixed consumption configuration for opcodes. If CommonDynamicCost, CallCost, SStoreCost are not used, 
    // this fixed consumption configuration for the opcode will be used.
    ConstCost [opcodes.MaxOpCodesCount]uint64
}

```

The [gasSetting](./gasSetting) package also provides Get and Set package-level methods to support users in customizing 
default Gas configurations or modifying them based on the default Gas settings to meet different needs.

```go
// Set custom default Gas configuration.
func Set(s *Setting)

// Get the current default Gas configuration.
func Get() *Setting
```

## Precompiled Contracts
SealEVM provides a custom precompiled contract registration interface within the reserved address space, 
offering better extensibility for different system requirements.  

Reserved address space: 0x0000000000000000000000000000000100000000 ~ 0x00000000000000000000000000000001FFFFFFFF.

>SealEVM has implemented precompiled contracts consistent with Ethereum to provide good application compatibility.

```go
// Definition of the precompiled contract interface, all registered precompiled contracts must implement this interface
type PrecompiledContract interface {
    GasCost(input []byte) uint64          //return the Gas consumption of the contract 
    Execute(input []byte) ([]byte, error) //execute function
}

// Precompiled contract registration function
// 'addr' is the address where the precompiled contract is registered, it must be within the 
// reserved address space of SealEVM, otherwise the registration will fail
func RegisterContracts(addr types.Address, c PrecompiledContract) error

```

## Usage Scenarios
SealEVM is an independent, flexible, configurable, and well-structured EVM execution environment. If you have the following needs, developing based on SealEVM would be a good choice:
- EVM environment in modular blockchain systems
- EVM environment in Layer 2 and Layer 3
- EVM environment with customizable GAS fees and precompiled contracts

>#### User Case

[ChainMaker](https://git.chainmaker.org.cn/chainmaker/vm-evm)

---

# License

[Apache License 2.0](https://raw.githubusercontent.com/SealSC/SealEVM/master/LICENSE)
