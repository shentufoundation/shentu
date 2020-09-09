pragma solidity >=0.4.0 <0.7.0;

contract CertifyValidator {
    string valConsPub = "certikvalconspub1zcjduepq32v65eegk2yvgzdya5dqnlnc063u7mt3dh66z2xyv9rddgm6t94s4pjeat";

    function certifyValidator() public returns (bytes memory) {
        bytes memory input = bytes(valConsPub);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := call(50000, 0x0c, 0, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }
}
