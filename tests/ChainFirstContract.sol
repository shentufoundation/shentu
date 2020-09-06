pragma solidity >=0.4.22 <0.7.0;
import "./ChainSecondContract.sol";
contract ChainFirstContract {
    event FirstEvent(string s);
    ChainSecondContract sc;

    constructor() public {
        emit FirstEvent("Does event work in the constructor ?");
        sc = new ChainSecondContract();
    }
}