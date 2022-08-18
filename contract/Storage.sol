// SPDX-License-Identifier: MIT
pragma solidity ^0.8.16;

/**
* @title Storage
* @dev store or retrieve variable value
*/
contract Storage {

	uint256 value;

	function store(uint256 number) public{
		value = number;
	}

	function retrieve() public view returns (uint256){
		return value;
	}
}
