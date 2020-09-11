package runner

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	oracleTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
	"github.com/certikfoundation/shentu/x/oracle"
)

// Push pushes MsgInquiryEvent to certik chain.
func Push(ctx oracleTypes.Context, ctkMsgChan <-chan interface{}, msgRespChan chan<- interface{}, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Context().Done():
			return
		case msg := <-ctkMsgChan:
			switch m := msg.(type) {
			case oracleTypes.MsgRespondToTask:
				go PushMsgRespondToTask(
					ctx.WithLoggerLabels("type", "msgRespondToTask"),
					m,
					msgRespChan,
					errorChan,
				)
			}
		}
	}
}

// PushMsgRespondToTask pushes MsgInquiryEvent message to CertiK Chain.
func PushMsgRespondToTask(
	ctx oracleTypes.Context,
	msgRespondToTask oracleTypes.MsgRespondToTask,
	msgRespChan chan<- interface{},
	errorChan chan<- error,
) {
	logger := ctx.Logger()
	// TODO: handle certik response error
	msgRespondToTaskTx := oracle.NewMsgTaskResponse(
		msgRespondToTask.Contract,
		msgRespondToTask.Function,
		int64(msgRespondToTask.Score),
		msgRespondToTask.Operator,
	)

	if err := msgRespondToTaskTx.ValidateBasic(); err != nil {
		ctx.Logger().Error(err.Error())
		errorChan <- err
		return
	}

	receipt, err := CompleteAndBroadcastTx(ctx.ClientContext(), ctx.TxBuilder(), []sdk.Msg{msgRespondToTaskTx})
	if err != nil {
		ctx.Logger().Error(err.Error())
		errorChan <- err
		return
	}
	ctx.Logger().Debug(
		"CertiK MsgTaskResponse Receipt",
		"height", receipt.Height,
		"txHash", receipt.TxHash,
		"data", receipt.Data,
		"rawLog", receipt.RawLog,
		"info", receipt.Info,
		"gasUsed", receipt.GasUsed,
		"timestamp", receipt.Timestamp,
	)
	msgRespChan <- receipt
	logger.Info("Finished pushing task response back to certik-chain")
}
