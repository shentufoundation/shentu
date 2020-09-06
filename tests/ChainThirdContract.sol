contract ChainThirdContract {
    event ThirdEvent(string s, string s2, string s3);

    constructor() public {
        emit ThirdEvent("Third","time's","the charm");
    }
}