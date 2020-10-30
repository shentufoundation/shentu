package utils

import "errors"

var (
	// ErrBaseAccount is the error for BaseAccount assertion
	ErrBaseAccount = errors.New("The account is not a BaseAccount")
	// ErrEmptyCVMCode is the error when the query of the CVM code is empty
	ErrEmptyCVMCode = errors.New("The cvm code is empty")
	// ErrEmptyCVMAbi is the error when the query of the CVM abi is empty
	ErrEmptyCVMAbi = errors.New("The cvm abi is empty")
)
