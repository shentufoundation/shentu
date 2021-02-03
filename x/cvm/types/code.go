package types

import "github.com/hyperledger/burrow/acm"

// CVMCodeType is the type for code in CVM.
type CVMCodeType byte

// CVM code types
const (
	CVMCodeTypeEVMCode CVMCodeType = iota
	CVMCodeTypeEWASMCode
)

// CVMCode defines the data structure of code in CVM.
type CVMCode struct {
	CodeType CVMCodeType
	Code     acm.Bytecode
}

// NewCVMCode returns a new CVM code instance.
func NewCVMCode(codeType CVMCodeType, code []byte) CVMCode {
	return CVMCode{
		CodeType: codeType,
		Code:     code,
	}
}
