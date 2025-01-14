# SealEVM 示例

本目录包含了演示 SealEVM 使用方法的示例代码。以下是编译和运行示例的说明，以及各个文件用途的描述。

> **免责声明**：这些示例代码仅用于演示 SealEVM 的使用方法，不适用于任何商业用途。使用者需自行承担相关风险。

## 编译

要编译示例，只需运行：

```bash
./compile.sh
```

## 示例目录结构

### /example 根目录文件
- `compile.sh`：编译 Solidity 合约并生成可执行文件
- `evm.go`：EVM 核心初始化和配置
- `storage.go`：存储实现示例
- `basicExample.go`：演示基本合约部署和交互的示例
- `precompiledWithStorageExample.go`：展示带存储功能的预编译合约使用示例
- `deployHelper.go`：合约部署辅助函数
- `printer.go`：输出格式化工具函数

### /example/contracts 目录
- `basicExample.sol`：演示基本功能的 Solidity 合约
- `basicExampleCodes.go`：执行 basicExample.sol 合约所需要的字节码
- `crossTxDataShareExample.sol`：演示跨交易数据共享的 Solidity 合约
- `crossTxDataShareExampleCodes.go`：执行 crossTxDataShareExample.sol 合约所需要的字节码

### /example/precompiledWithStorage 目录
- `crossTxDataShare.go`：实现带存储功能的预编译合约示例

## 示例概述

1. 基础示例（`basicExample.go`）：
   - 演示基本的合约部署和交互
   - 使用 `contracts/basicExample.sol` 中定义的合约

2. 带存储的预编译合约示例（`precompiledWithStorageExample.go`）：
   - 展示如何使用带存储功能的预编译合约
   - 实现跨交易数据共享
   - 使用了 `contracts/crossTxDataShareExample.sol` 和 `precompiledWithStorage/crossTxDataShare.go` 中的合约

每个示例都可以在编译后独立运行。有关详细的实现和使用示例，请参考相应的 Go 文件。
