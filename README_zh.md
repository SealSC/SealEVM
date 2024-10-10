# SealEVM

SealEVM是一个独立的EVM执行器，目标是实现一个完全与存储系统解耦的EVM执行环境，来为任意的区块链系统增加EVM支持。
当前版本已经实现了通过接口和缓存的方式实现了与存储系统的解耦，支持为任意使用golang实现的区块链系统增加EVM支持。

**[English](https://github.com/SealSC/SealEVM/blob/master/README.md) | 中文**

##

- [如何使用](#如何使用)
- [主要结构体与接口](#主要结构体与接口)
  - [创建EVM实例配置参数](#创建evm实例配置参数)
  - [外部存储接口](#外部存储接口)
  - [执行环境结构体](#执行环境结构体)
  - [执行结果结构体](#执行结果结构体)
  - [执行结果缓存结构体](#执行结果缓存结构体)
- [Gas配置](#gas设置)
- [预编译合约](#预编译合约)
- [使用场景](#使用场景)
  - [使用案例](#使用案例)

------

## 如何使用

```go
func main() {
    //载入SealEVM模块，SealEVM没有使用golang的init特性，必须显式的在使用之前，全局执行一次
    SealEVM.Load()
    
    //根据需要，配置SealEVM实例参数，交易输入等参数都在该结构体中提供，参数结构体说明见后续章节
    evmParam := SealEVM.EVMParam{}
    
    //提供该参数来初始化一个SealEVM实例
    evm := SealEVM.New(evmParam)
    
    //evm执行，返回值中的result是一个ExecuteResult结构体
    //该结构体存储了数据的原始状态、最终状态、合约Log、内部创建合约等信息，该结构体说明见后续章节
    result, err := evm.ExecuteContract()
}
```
[**example**](https://github.com/SealSC/SealEVM/tree/master/example)目录下，提供了一个简单的SealEVM的使用参考示例。该示例使用了内存作为外部存储，展示了简单的合约部署、调用、变量读取等功能。

## 主要结构体与接口

>#### 创建EVM实例配置参数
```go
//执行结果回调函数定义，接口中的result和err，与evm执行返回值是相同的
type EVMResultCallback func(result ExecuteResult, err error)

type EVMParam struct {
    MaxStackDepth  int //最大执行栈深度，注意，不是存储栈深度
    ExternalStore  storage.IExternalStorage //外部存储接口，说明见后续章节
    ResultCallback EVMResultCallback //EVM执行完成后的回调函数，定义见本代码段开头
    Context        *environment.Context //EVM执行时的环境上下文结构体，说明见后续章节
    GasSetting     *gasSetting.Setting //Gas费用设置，nil时使用默认设置，说明见后续章节
}
```

##

>#### 外部存储接口
SealEVM将通过该接口，与外部存储进行交互，来实现必要的合约读取、状态读取、地址创建等功能。

```go
type IExternalStorage interface {
    //获取合约已存储的合约
    GetContract(address types.Address) (*environment.Contract, error)
    
    //获取指定高度的区块哈希
    GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)
    
    //检查合约是否存在
    ContractExist(address types.Address) bool
    
    //检查地址是否为空，空的定义请参与EIP-161
    ContractEmpty(address types.Address) bool
    
    //返回给定合约代码的哈希值
    HashOfCode(code []byte) types.Hash
    
    //根据参数，返回创建的合约地址，操作码 CREATE（0xF0）创建合约时使用
    CreateAddress(caller types.Address, tx environment.Transaction) types.Address
    
    //根据参数，返回创建的合约地址，操作码 CREATE2（0xF5）创建合约时使用
    CreateFixedAddress(caller types.Address, salt types.Hash, code []byte, tx environment.Transaction) types.Address
    
    //在执行opcode SLOAD(0x54) 时，从外部存储获取指定位置的256位数据
    Load(address types.Address, slot types.Slot) (*evmInt256.Int, error)
}
```

##

>#### 执行环境结构体
该结构体在environment包内，是SealEVM执行时的执行上下文，包括区块、交易、消息、合约等参数。

```go

//执行环境汇总结构体，Contract为nil时，视为创建合约交易
//执行创建合约交易时，会使用本结构体中的Message.Data字段作为合约代码，生成合约结构体实例赋值给Contract字段
type Context struct {
    Block       Block       //区块环境结构体，详细信息见本代码段下文
    Contract    *Contract   //本次执行的合约结构体，详细信息见本代码段下文，
    Transaction Transaction //交易结构体，详细信息见本代码段下文
    Message     Message     //消息结构体，详细信息见本代码段下文
}

//区块环境结构体
type Block struct {
    ChainID     *evmInt256.Int //区块ID
    Coinbase    types.Address  //出块人地址
    Timestamp   uint64         //区块的秒级时间戳
    Number      uint64         //区块高度，使用的是64位无符号整数
    Difficulty  *evmInt256.Int //难度
    GasLimit    *evmInt256.Int //区块gas限制
    Hash        types.Hash     //区块哈希
    BaseFee     *evmInt256.Int //EIP-1559中定义的base-fee
    BlobBaseFee *evmInt256.Int //EIP-7516中定义的blob base-fee
}

//合约环境结构体，
type Contract struct {
    Address  types.Address  //要执行的合约的地址
    Code     []byte         //合约代码
    CodeHash types.Hash     //合约的哈希
    CodeSize uint64         //合约代码字节大小
    Balance  *evmInt256.Int //合约地址的账户余额
}

//交易环境结构体
type Transaction struct {
    TxHash   types.Hash     //交易哈希
    Origin   types.Address  //发起本次交易的地址，操作码ORIGIN(0x32)获取到的值
    GasPrice *evmInt256.Int //交易的gas价格
    GasLimit *evmInt256.Int //交易的gas限制
    
    BlobHashes []types.Hash //EIP-4844中的tx.blob_versioned_hashes
}

//消息结构体
type Message struct {
    Caller types.Address //合约调用者地址，操作码CALLER(0x33)获取到的值
    Value  *evmInt256.Int //合约调用时，发送的ETH数量
    Data   []byte //合约调用时，传递给合约的参数
}
```

##

>#### 执行结果结构体
SealEVM在执行合约时，会将除新合约部署外的，所有有变动的数据，放入缓存中，不会通知给外部存储。
```go
type ExecuteResult struct {
    ResultData   []byte //合约执行返回的数据
    GasLeft      uint64 //剩余gas
    StorageCache storage.ResultCache //缓存结构体，说明见后续章节
    ExitOpCode   opcodes.OpCode //执行完毕时，最后一个执行的opcode
}
```

##

>#### 执行结果缓存结构体
SealEVM会将执行过程中以及执行完毕后，从外部获取的原始数据，以及需要最终存储结果数据，
放入一个统一缓存结构体实例，在合约执行完毕后作为结果返回给调用者。
这种设计让外部存储不需要关心执行中的状态变更，只需要提供初始数据，就能够得到执行后的结果。

```go
//数据缓存类型，是Slot缓存的最下层map
type Cache map[types.Slot]*evmInt256.Int

//Slot缓存类型，通过[address][slot]来寻址访问数据
type SlotCache map[types.Address]Cache

//日志缓存类型，会顺序的存放合约执行过程中，依次产生的Log数据
type LogCache []*types.Log

//销毁合约地址缓存类型，存放执行了SELFDESTRUCT(0xFF)的合约地址
type DestructCache map[types.Address]types.Address

//合约缓存类型，存放内部交易创建的合约实例
type ContractCache map[types.Address]*environment.Contract

//缓存汇总结构体
type ResultCache struct {
    OriginalData SlotCache //从外部存储通过SLOAD载入的原始数据
    CachedData   SlotCache //合约执行后，SSTORE存入的数据
    
    //TOriginalData和TCachedData是Transient storage的缓存，
    //该类型缓存是EIP-1153引入的，是合约执行过程中的临时存储空间
    TOriginalData SlotCache
    TCachedData   SlotCache
    
    Logs         *LogCache //操作码LOG0(0xA0)~LOG4(0xA4)产生的日志缓存
    Destructs    DestructCache //执行了SELFDESTRUCT(0xFF)的合约的缓存
    NewContracts ContractCache //执行过程中，内部交易创建的合约的缓存
}
```

## Gas设置
SealEVM通过[gasSetting](./gasSetting)包来实现灵活的Gas设置，并且提供了一个尽可能与以太坊Gas系统一致的默认配置。

```go
//交易固有Gas费用计算函数类型定义，data为交易的输入数据，也就是参数，to为合约结构体指针，to为nil时表示本次交易为创建合约交易
//返回值为针对交易的固有Gas消耗量(gasCost)
type intrinsicGasSetting.IntrinsicGas func(data []byte, to *environment.Contract) (gasCost uint64)

//通用的动态Gas消耗计算函数类型定义
//需要返回要扩展的内存大小(memExpSize)、gas消耗量(gasCost)
type dynamicGasSetting.CommonCalculator func(
    contract *environment.Contract, //操作码执行的合约环境变量
    stx *stack.Stack,               //操作码执行开始时的堆栈环境
    mem *memory.Memory,             //操作码执行开始时的内存环境
    store *storage.Storage,         //操作码开始执行时的存储环境
) (memExpSize uint64, gasCost uint64, err error)

//为CALL、CALLCODE、STATICCALL、DELEGATECALL设计的Gas消耗计算函数类型定义
//需要返回要扩展的内存大小(memExpSize)、gas消耗量(gasCost)、内部调用发送的gas量(sendGas)
type CallGas func(
    code opcodes.OpCode,    //当前的操作码
    availableGas uint64,    //当前可使用的Gas数量
    stx *stack.Stack,       //操作码执行开始时的堆栈环境
    mem *memory.Memory,     //操作码执行开始时的内存环境
    store *storage.Storage, //操作码开始执行时的存储环境
) (expSize uint64, gasCost uint64, sendGas uint64, err error)

//为存储新合约设计的Gas消耗计算函数类型定义
//需要返回gas消耗量(gasCost)
type ContractStoreGas func(code []byte, gasRemaining uint64) (gasCost uint64, err error)

type Setting struct {
    //每个交易的固定Gas费用计算函数，会在执行开始前调用一次，以扣除固定的交易费用
    IntrinsicCost intrinsicGasSetting.IntrinsicGas

    //动态Gas消耗计算函数配置，CALL、CALLCODE、STATICCALL、DELEGATECALL、SSTORE操作码将忽略该配置
    CommonDynamicCost [opcodes.MaxOpCodesCount]dynamicGasSetting.CommonCalculator

    //内部调用CALL、CALLCODE、STATICCALL、DELEGATECALL时的gas计算配置
    CallCost [opcodes.MaxOpCodesCount]dynamicGasSetting.CallGas

    //Create、Create2以及创建合约的交易，在存储合约代码时的Gas计算配置
    ContractStoreCost dynamicGasSetting.ContractStoreGas

    //操作码固定消耗配置，如果CommonDynamicCost、CallCost、SStoreCost都未被使用，则使用该操作码的固定消耗配置
    ConstCost [opcodes.MaxOpCodesCount]uint64
}
```

[gasSetting](./gasSetting)包，为支持使用者自定义默认Gas配置，或者基于该默认Gas设置进行修改，以满足不同的需求场景。

```go
//设置自定义的默认Gas配置
func Set(s *Setting)

//获取当前默认Gas配置
func Get() *Setting
```

## 预编译合约
SealEVM在保留地址空间内，提供了自定义预编译合约注册接口，来为不同系统需求提供更好的扩展性。  

保留地址空间:  
0x0000000000000000000000000000000100000000 ~ 0x00000000000000000000000000000001FFFFFFFF。

>SealEVM已经实现了与以太坊一致的预编译合约，来提供良好的应用兼容性。

```go
//预编译合约接口定义，所有注册的预编译合约必须实现该接口
type PrecompiledContract interface {
    GasCost(input []byte) uint64 //返回预编译合约的Gas消耗
    Execute(input []byte) ([]byte, error) //预编译合约执行函数
}

//预编译合约注册函数
//addr为预编译合约注册的地址，范围必须在SealEVM保留地址空间内，否则会注册失败
func RegisterContracts(addr types.Address, c PrecompiledContract) error

```

## 使用场景
SealEVM是一个独立的，灵活可配置，结构良好的EVM执行环境，因此如果您有以下需求，那么基于SealEVM进行开发，会是一个不错的选择：
- 模块化区块链系统中EVM环境
- Layer2、Layer3中的EVM环境
- 定制化GAS费用、预编译合约的EVM环境

>#### 使用案例

[长安链](https://git.chainmaker.org.cn/chainmaker/vm-evm)

---

# License

[Apache License 2.0](https://raw.githubusercontent.com/SealSC/SealEVM/master/LICENSE)
