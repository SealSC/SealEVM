// SPDX-License-Identifier: MIT
pragma solidity 0.8.18;

contract subContract {
    event Info(string info);
    constructor(){}

    function call(string memory info) external {
        emit Info(info);
    }
}

contract example {
    constructor(){}
    uint256 public counter;

    event Count(string reason, uint256 indexed value);

    function increaseFor(string calldata reason) external returns(uint256) {
        counter += 1;
        emit Count(reason, counter);

        subContract sub = new subContract();
        sub.call(reason);

        return counter;
    }
}
