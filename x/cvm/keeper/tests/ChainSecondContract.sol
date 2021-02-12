pragma solidity >0.7.0;
import "./ChainThirdContract.sol";

contract ChainSecondContract {
    event SecondEvent(string s, string s2);
    ChainThirdContract tc;
    constructor() public {
        emit SecondEvent("Triggering a second","event");
        tc = new ChainThirdContract();

    }
}