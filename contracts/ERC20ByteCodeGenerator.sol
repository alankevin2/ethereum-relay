// SPDX-License-Identifier: UNLICENSED 

pragma solidity ^0.8.11;

import "./openzeppelin/tokens/ERC20/ERC20.sol";
import "./openzeppelin/access/Ownable.sol";

abstract contract ERC20ByteCodeGenerator is ERC20, Ownable {
}