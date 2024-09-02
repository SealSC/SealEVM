# SealEVM

SealEVM是一个独立的EVM实现，它通过接口和缓存方式，实现与存储系统的解耦，可以轻松的移植到任意使用golang实现的区块链系统中，为其增加EVM支持。

---

- [English](https://github.com/SealSC/SealEVM/tree/master#readme)

##

### 示例说明

[example](https://github.com/SealSC/SealEVM/tree/master/example)目录下，提供了一个简单的SealEVM的使用参考示例。该示例使用了内存作为外部存储，展示了简单的合约部署、调用、变量读取等功能。

**⚠️注意：example目录下的示例仅做代码使用的简单展示，请勿用于任何实际商业和生产环境中**

##

### 主要结构体与接口

**⚠️注意，源码中的namespace，是address的别名，其意义与以太坊中的address一致。**

##

>#### 创建EVM实例配置参数
```go
type EVMParam struct {
	MaxStackDepth  int //最大栈深度
	ExternalStore  storage.IExternalStorage //外部存储接口，后续章节会详细说明
	ResultCallback EVMResultCallback //EVM执行完成后的回调函数
	Context        *environment.Context //EVM执行时的环境上下文，内部字段含义请阅读源码
	GasSetting     *instructions.GasSetting //OpCode的自定义gas费用设置
}
```

##

>#### 外部存储接口
SealEVM将通过该接口，与外部存储进行交互，来实现必要的状态读取、地址创建、新合约存储等功能。
```go
type IExternalStorage interface {
    //从外部存储获取指定地址的账户余额
    GetBalance(address *evmInt256.Int) (*evmInt256.Int, error)
    
    //从外部存储获取指定地址的合约代码
    GetCode(address *evmInt256.Int) ([]byte, error)
    
    //从外部存储获取指定地址的合约代码大小
    GetCodeSize(address *evmInt256.Int) (*evmInt256.Int, error)
    
    //从外部存储获取指定地址合约代码的哈希
    GetCodeHash(address *evmInt256.Int) (*evmInt256.Int, error)
    
    //从外部存储获取指定区块的哈希
    GetBlockHash(block *evmInt256.Int) (*evmInt256.Int, error)
    
    //在执行opcode CREAT(0xF0)时，将调用该方法来获取创建的合约的地址
    CreateAddress(caller *evmInt256.Int, tx environment.Transaction) *evmInt256.Int
    
    //在执行opcode CREAT2(0xF5)时，将调用该方法来获取创建的合约的地址
    CreateFixedAddress(caller *evmInt256.Int, salt *evmInt256.Int, tx environment.Transaction) *evmInt256.Int
    
    //从外部存储获取是否可以进行账户余额的转移
    CanTransfer(from *evmInt256.Int, to *evmInt256.Int, amount *evmInt256.Int) bool
    
    //在执行opcode SLOAD(0x54) 时，从外部存储获取指定位置的256位数据
    //注意：参数n是当前执行的合约的地址，参数k是执行opcode SLOAD(0x54)时，给出的存储位置的key
    Load(n *evmInt256.Int, k *evmInt256.Int) (*evmInt256.Int, error)

    //执行CREATE或CREATE2成功后的外部存储回调，创建的新合约的地址和代码，会通过该接口提供给外部存储
    //注意：n就是新合约的地址，code就是新合约的bytecode
    NewContract(n *evmInt256.Int, code []byte) error
}
```

##

>#### 执行结果
SealEVM在执行合约时，会将除新合约部署外的，所有有变动的数据，放入缓存中，不会通知给外部存储。
```go
type ExecuteResult struct {
    ResultData   []byte //合约执行返回的数据
    GasLeft      uint64 //剩余gas
    StorageCache storage.ResultCache //外部状态变化的缓存，外部数据需要根据该缓存更新存储数据，下面会详细说明
    ExitOpCode   opcodes.OpCode //执行完毕时，最后一个执行的opcode
}
```

##

>#### 执行结果的缓存
下面是关于这些缓存变量的作用的描述，详细的结构请参考源码
```go
type ResultCache struct {
    OriginalData CacheUnderNamespace //缓存执行过程中，从外部存储读取的信息
    CachedData   CacheUnderNamespace //缓存执行过程中，所有的状态更新的最终结果

    TOriginalData CacheUnderNamespace //缓存执行过程中，从外部瞬时存储读取的信息
    TCachedData   CacheUnderNamespace //缓存执行过程中，所有的瞬时存储状态更新的最终结果
	
    Balance   BalanceCache //缓存执行过程中，所有的账户余额变化
    Logs      LogCache //缓存执行过程中，合约产生的所有日志
    Destructs Cache //缓存执行了opcode SELFDESTRUCT(0xff)的合约地址缓存
}
```

---

# License

[Apache License 2.0](https://raw.githubusercontent.com/SealSC/SealEVM/master/LICENSE)

