pragma solidity >=0.4.0 <0.7.0;

contract Simpleevent {
    uint storedData;
    event MyEvent(uint myNumber);
    event MySecondEvent(uint mySecondNumber, string s);
    event MyThirdEvent(uint myThirdNumber, uint8[3] arr, address a );

    function set(uint x) public {
        storedData = x;
        emit MyEvent(storedData + 1);
        emit MySecondEvent(2 * storedData, "secondEvent!!!");
        uint8[3] memory arr = [1,2,3];
        emit MyThirdEvent(storedData + 3, arr, address(this));
    }

    function get() public view returns (uint) {
        return storedData;
    }
}
