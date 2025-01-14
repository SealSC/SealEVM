// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface ICrossTxDataShare {
    function share(bytes32 slot, bytes calldata data) external returns (bytes32);
    function read(bytes32 slot) external view returns (bytes memory);
}

contract CrossTxDataShareExample {
    // Address of the CrossTxDataShare precompiled contract
    ICrossTxDataShare constant CROSS_TX_DATA_SHARE = ICrossTxDataShare(address(0x20001));

    // Events for tracking operations
    event DataShared(bytes32 slot, bytes data);
    event DataRead(bytes32 slot, bytes data);

    /**
     * @dev Share data that can be accessed across transactions
     * @param slot The storage slot where the data will be stored
     * @param data The data to be shared
     */
    function shareData(bytes32 slot, bytes calldata data) external returns (bytes32) {
        bytes32 result = CROSS_TX_DATA_SHARE.share(slot, data);
        emit DataShared(slot, data);
        return result;
    }

    /**
     * @dev Read previously shared data
     * @param slot The storage slot (sharing ID) to read from
     * @return data The stored data
     */
    function readSharedData(bytes32 slot) external returns (bytes memory data) {
        // data = CROSS_TX_DATA_SHARE.read(slot);
        emit DataRead(slot, data);
        return data;
    }
} 