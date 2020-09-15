package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
	"github.com/certikfoundation/shentu/x/oracle"
)

// Push pushes MsgInquiryEvent to certik chain.
func Push(ctx types.Context, ctkMsgChan <-chan interface{}, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Context().Done():
			return
		case msg := <-ctkMsgChan:
			switch m := msg.(type) {
			case oracle.MsgTaskResponse:
				go PushMsgTaskResponse(ctx.WithLoggerLabels("type", "MsgTaskResponse"), m)
			}
		}
	}
}

// PushMsgTaskResponse pushes MsgTaskResponse message to CertiK Chain.
func PushMsgTaskResponse(ctx types.Context, msg oracle.MsgTaskResponse) {
	logger := ctx.Logger()
	if err := msg.ValidateBasic(); err != nil {
		ctx.Logger().Error(err.Error())
		return
	}
	receipt, err := CompleteAndBroadcastTx(ctx.ClientContext(), ctx.TxBuilder(), []sdk.Msg{msg})
	if err != nil {
		ctx.Logger().Error(err.Error())
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
	logger.Debug("Finished pushing task response back to certik-chain")
}
