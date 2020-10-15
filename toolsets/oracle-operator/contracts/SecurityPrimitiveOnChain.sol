pragma solidity ^0.6.0;

interface ISecurityPrimitive {
    function getInsight(string calldata contractAddress, string calldata functionSignature) external view returns (bool, string memory);
}

contract SecurityPrimitiveOnChain is ISecurityPrimitive {
    address _owner;

    event OwnershipTransferred(address previousOwner, address newOwner);

    modifier onlyOwner() {
        require(_owner == msg.sender, "caller is not the owner");
        _;
    }

    constructor () public {
        _owner = msg.sender;
    }

    function renounceOwnership() public onlyOwner {
        emit OwnershipTransferred(_owner, address(0));
        _owner = address(0);
    }

    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "new owner is the zero address");
        emit OwnershipTransferred(_owner, newOwner);
        _owner = newOwner;
    }

    function getInsight(string memory contractAddress, string memory functionSignature) public view override returns (bool, string memory) {
        bytes memory input = bytes(contractAddress);
        bytes memory result = new bytes(1);
        assembly {
            let len := mload(input)
            let success := staticcall(0, 0x09, add(input, 0x20), len, add(result, 0x20), 0x01)
        }
        if (result[0] == 0x01) {
            return (false, "100");
        } else {
            return (false, "50");
        }
    }
}
