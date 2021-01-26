pragma solidity >=0.4.26 <0.7.0;

contract StringTest {
    string a;

    function gets() public returns(string memory) {
      return a;
    }

    function getl() public returns(uint) {
        return bytes(a).length;
    }

    function changeString(string memory _c) public returns(string memory){
        a = _c;
        return (a);
    }

    function changeGiven(string memory _s) public returns(string memory) {
        bytes memory byteString = bytes(_s);
        byteString[0] = 'A';
        byteString[1] = 'b';
        byteString[2] = 'c';

        return string(byteString);
    }

    function testStuff() public returns(uint) {
        string memory a = "asdfdd000";
        string memory b = "asdfdd000";
        string memory c = "asdfdd001";
        require(keccak256(abi.encodePacked(a))
              ==keccak256(abi.encodePacked(b)));
        require(keccak256(abi.encodePacked(b))
              !=keccak256(abi.encodePacked(c)));

        string memory t = "€bÁnç!";
        string memory t2 = "AbcbÁnç!";
        bytes memory bt2 = bytes(t2);
        bytes memory byteString = bytes(t);
        byteString[0] = "A";
        byteString[1] = "b";
        byteString[2] = "c";
        require(keccak256(abi.encodePacked(bt2))
              ==keccak256(abi.encodePacked(byteString)));

        return 123123;
    }
}
