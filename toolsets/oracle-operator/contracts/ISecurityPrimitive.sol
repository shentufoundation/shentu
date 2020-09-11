pragma solidity ^0.6.0;

interface ISecurityPrimitive {
    function getInsight(string calldata contractAddress, string calldata functionSignature) external view returns (bool, string memory);
}
