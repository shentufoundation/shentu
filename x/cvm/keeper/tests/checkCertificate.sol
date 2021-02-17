pragma solidity >0.7.0;

contract CheckCertificate {
    string not_certified = "certik1udzp23twf4a6pf47er26ujgczzh3tcst7c5aze";
    string auditing = "cosmos1xxkueklal9vejv9unqu80w9vptyepfa95pd53u";
    string proof = "cosmos1r60hj2xaxn79qth4pkjm9t27l985xfsmnz9paw";
    string compilation = "certik1udzp23twf4a6pf47er26ujgczzh3tcst7c5aze";
    string everything = "cosmos1xxkueklal9vejv9unqu80w9vptyepfa95pd53u";
    string sourceCodeHash = "dummysourcecodehash";

    function callCheck() public returns (bytes memory) {
        bytes memory input = bytes(auditing);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x65, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }

    function callCheckNotCertified() public returns (bytes memory) {
        bytes memory input = bytes(not_certified);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x65, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }

    function proofCheck() public returns (bytes memory) {
        bytes memory input = bytes(proof);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x66, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }

    function compilationCheck() public returns (bytes memory) {
        bytes memory input = bytes(sourceCodeHash);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x67, add(input, 0x20), len, out, 0x01)
            return (out,0x01)
        }
    }

    function proofAndAuditingCheck() public returns (bytes memory) {
        bytes memory input = bytes(everything);
        assembly {
            let out := 0x01
            let len := mload(input)
            let success := staticcall(50000, 0x65, add(input, 0x20), len, out, 0x01)
            let success2 := staticcall(50000, 0x66, add(input, 0x20), len, out, 0x01)
            if eq(success, 0x01) {
                if eq(success2, 0x01) {
                    return (0x01,0x01)
                }
            }
            return (0x00,0x01)
        }
    }
}
