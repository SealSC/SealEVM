# SealEVM Examples

This directory contains examples demonstrating the usage of SealEVM. Below are instructions on how to compile and run the examples, as well as a description of each file's purpose.

> **Disclaimer**: These example codes are for demonstration purposes only to show the usage of SealEVM. They are not intended for any commercial use. Use at your own risk.

## Compilation

To compile the examples, simply run:

```bash
./compile.sh
```

## Example Directory Structure

### /example Root Directory Files
- `compile.sh`: Script for compiling Solidity contracts and generating executable files
- `evm.go`: Core EVM initialization and configuration
- `storage.go`: Storage implementation example
- `basicExample.go`: Basic example demonstrating simple contract deployment and interaction
- `precompiledWithStorageExample.go`: Example showing precompiled contract usage with storage
- `deployHelper.go`: Helper functions for contract deployment
- `printer.go`: Utility functions for output formatting

### /example/contracts Directory
- `basicExample.sol`: Simple Solidity contract for basic functionality demonstration
- `basicExampleCodes.go`: Bytecodes for basicExample.sol
- `crossTxDataShareExample.sol`: Solidity contract demonstrating cross-transaction data sharing
- `crossTxDataShareExampleCodes.go`: Bytecodes for crossTxDataShareExample.sol

### /example/precompiledWithStorage Directory
- `crossTxDataShare.go`: Implementing an example of a precompiled contract with storage functionality

## Examples Overview

1. Basic Example (`basicExample.go`):
   - Demonstrates basic contract deployment and interaction
   - Uses the contract defined in `contracts/basicExample.sol`

2. Precompiled With Storage Example (`precompiledWithStorageExample.go`):
   - Shows how to use precompiled contracts with storage functionality
   - Implements cross-transaction data sharing
   - Uses contracts from both `contracts/crossTxDataShareExample.sol` and `precompiledWithStorage/crossTxDataShare.go`

Each example can be run independently after compilation. For detailed implementation and usage examples, please refer to the respective Go files.
