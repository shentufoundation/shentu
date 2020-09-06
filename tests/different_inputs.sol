pragma solidity >=0.4.22 <0.7.0;

contract DifferentInputs {
    function iNeedAnInt(uint a) public {
    }

    function iNeedAString(string memory s) public {
    }

    function iNeedAnAddress(address payable a) public payable {
        a.transfer(msg.value);
	}
}