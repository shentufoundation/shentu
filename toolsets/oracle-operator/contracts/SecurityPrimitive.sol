pragma solidity ^0.6.0;

interface ISecurityPrimitive {
    function getInsight(string calldata contractAddress, string calldata functionSignature) external view returns (bool, string memory);
}

contract SecurityPrimitive is ISecurityPrimitive {
    address _owner;
    string _endpoint;

    event OwnershipTransferred(address previousOwner, address newOwner);
    event EndpointChanged(string previousOwner, string newOwner);

    modifier onlyOwner() {
        require(_owner == msg.sender, "caller is not the owner");
        _;
    }

    constructor (string memory endpoint) public {
        _owner = msg.sender;
        _endpoint = endpoint;
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

    function changeEndpoint(string memory newEndpoint) public onlyOwner {
        emit EndpointChanged(_endpoint, newEndpoint);
        _endpoint = newEndpoint;
    }

    function getInsight(string memory contractAddress, string memory functionSignature) public view override returns (bool, string memory) {
        return (true, getEndpointUrl(contractAddress, functionSignature));
    }

    function getEndpointUrl(string memory contractAddress, string memory functionSignature) public view returns (string memory) {
        // TODO: support flexible url pattern construction
        return string(abi.encodePacked(_endpoint, "?address=", contractAddress, "&functionSignature=", functionSignature));
    }
}