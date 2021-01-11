pragma solidity >=0.4.0;

contract CertifyValidator {
    string valConsPub = "cosmosvalconspub1zcjduepqxhy6865hf90lwmckjuegfdvqmyznhd6a4dkjr90pq0a82fxxg2qqcpfqat";

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
