package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	govtyepsv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtyepsv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var (
	_, _, addr   = testdata.KeyTestPubAddr()
	govAcct      = authtypes.NewModuleAddress(types.ModuleName)
	TestProposal = getTestProposal()
)

func getTestProposal() []sdk.Msg {
	legacyProposalMsg, err := govtyepsv1.NewLegacyContent(govtyepsv1beta1.NewTextProposal("Title", "description"), authtypes.NewModuleAddress(types.ModuleName).String())
	if err != nil {
		panic(err)
	}

	return []sdk.Msg{
		banktypes.NewMsgSend(govAcct, addr, sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1000)))),
		legacyProposalMsg,
	}
}
