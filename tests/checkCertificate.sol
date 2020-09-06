pragma solidity >=0.4.0 <0.7.0;

contract CheckCertificate {
    string not_certified = "certik1j9thnw367d72txzrmsu2wsk5qr7ugd86hpt8qg";
    string auditing = "certik1q7h5e2gwpykq27etnd5k5fcvc2azpueujqcvap";
    string proof = "certik1m95kfvajw5dmnu9e6h365tqcnazm9auddngklv";
    string compilation = "certik1udzp23twf4a6pf47er26ujgczzh3tcst7c5aze";
    string everything = "certik1nzgvd4k34zzf6vk3qevuh5xx8xshg6uy0l8rd5";
    string sourceCodeHash = "dummysourcecodehash";

    function callCheck() public returns (bytes memory) {
        bytes memory input = bytes(auditing);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x09, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }

    function callCheckNotCertified() public returns (bytes memory) {
        bytes memory input = bytes(not_certified);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x09, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }

    function proofCheck() public returns (bytes memory) {
        bytes memory input = bytes(proof);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x0a, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }

    function compilationCheck() public returns (bytes memory) {
        bytes memory input = bytes(sourceCodeHash);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x0b, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }

    function proofAndAuditingCheck() public returns (bytes memory) {
        bytes memory input = bytes(everything);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x09, add(input, 0x20), len, out, 0x01)
            let success2 := staticcall(50000, 0x0a, add(input, 0x20), len, out, 0x01)
            if eq(success, 0x01) {
                if eq(success2, 0x01) {
                    return (0x01,0x01)
                }
            }
            return (0x00,0x01)
        }
    }
}
