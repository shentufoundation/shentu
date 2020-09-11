// Package runner runs operator.
package runner

import (
	oracleTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
)

// Name is the name of the module.
func Name() string { return "Oracle-Operator" }

// Start function starts operator runner.
func Start(ctx oracleTypes.Context) chan error {
	errorChan := make(chan error, 1000)
	ctkMsgChan := make(chan interface{}, 1000)
	respMsgChan := make(chan interface{}, 1000)

	go Listen(ctx.WithLoggerLabels("protocol", "certik", "submodule", "listener"), ctkMsgChan, errorChan)
	go Push(ctx.WithLoggerLabels("protocol", "certik", "submodule", "pusher"), ctkMsgChan, respMsgChan, errorChan)
	return errorChan
}
