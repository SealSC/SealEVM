# SealEVM

SealEVM是一个独立的EVM执行器，目标是实现一个完全与存储系统解耦的EVM执行环境，来为任意的区块链系统增加EVM支持。  
当前版本已经实现了通过接口和缓存的方式实现了与存储系统的解耦，支持为任意使用golang实现的区块链系统增加EVM支持。

##

⚠️ **Note: SealEVM 注重解耦和独立操作，确保操作码行为与 EVM 匹配，但 GAS 消耗和预编译合约可能与以太坊不同。**

##

- [English](https://github.com/SealSC/SealEVM/blob/master/README.md) | 中文

---

## 示例代码
[example](https://github.com/SealSC/SealEVM/tree/master/example)目录下，提供了一个简单的SealEVM的使用参考示例。该示例使用了内存作为外部存储，展示了简单的合约部署、调用、变量读取等功能。

**⚠️注意：example目录下的示例仅做代码使用的简单展示，请勿用于任何实际商业和生产环境中**

## 主要结构体与接口

>#### 创建EVM实例配置参数
```go
type EVMParam struct {
    MaxStackDepth  int //最大执行栈深度
    ExternalStore  storage.IExternalStorage //外部存储接口，说明见后续章节
    ResultCallback EVMResultCallback //EVM执行完成后的回调函数
    Context        *environment.Context //EVM执行时的环境上下文，内部字段含义
    GasSetting     *gasSetting.Setting //Gas费用设置，nil时使用默认设置
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

>#### 执行结果
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

>#### 执行结果的缓存
下面是关于这些缓存变量的作用的描述，详细的结构请参考源码。
```go
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

## Some Usage Scenarios
SealEVM是一个独立的，灵活可配置，结构良好的EVM执行环境，因此如果您有以下需求，那么基于SealEVM进行开发，会是一个不错的选择：
- 模块化区块链系统中EVM环境
- Layer2、Layer3中的EVM环境
- 定制化GAS费用、预编译合约的EVM环境

**使用者案例**

[长安链](https://git.chainmaker.org.cn/chainmaker/vm-evm)

---

# License

[Apache License 2.0](https://raw.githubusercontent.com/SealSC/SealEVM/master/LICENSE)
