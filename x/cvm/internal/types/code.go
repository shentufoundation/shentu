package types

import "github.com/hyperledger/burrow/acm"

// CVM code types
const (
	CVMCodeTypeEVMCode   = 0
	CVMCodeTypeEWASMCode = 1
)

type ContractMetas = []acm.ContractMeta

// NewCVMCode returns a new CVM code instance.
func NewCVMCode(codeType int64, code []byte) CVMCode {
	return CVMCode{
		CodeType: codeType,
		Code:     code,
	}
}
