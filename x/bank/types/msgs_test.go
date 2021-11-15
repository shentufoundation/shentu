package types

import (
	//"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
}

func (suite *IntegrationTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, types.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.BankKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = queryClient
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestMsgSendRoute() {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	UnlockerAddress := sdk.AccAddress([]byte("unlocker"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("uctk", 10))
	var msg = NewMsgLockedSend(addr1, addr2, UnlockerAddress.String(), coins)
	suite.Require().Equal(msg.Route(), bankTypes.RouterKey)
	suite.Require().Equal(msg.Type(), "locked_send")
}


func (suite *IntegrationTestSuite)  TestMsgSendValidation() {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	UnlockerAddress := sdk.AccAddress([]byte("unlocker"))
	CTK123 := sdk.NewCoins(sdk.NewInt64Coin("ctk", 123))
	CTK0 := sdk.NewCoins(sdk.NewInt64Coin("ctk", 0))
	CTK123eth123 := sdk.NewCoins(sdk.NewInt64Coin("ctk", 123), sdk.NewInt64Coin("eth", 123))
	CTK123eth0 := sdk.Coins{sdk.NewInt64Coin("ctk", 123), sdk.NewInt64Coin("eth", 0)}

	var emptyAddr sdk.AccAddress

	cases := []struct {
		tx    *MsgLockedSend
		valid bool
	}{
		{NewMsgLockedSend(addr1, addr2, UnlockerAddress.String(), CTK123), true},       // valid send
		{NewMsgLockedSend(addr1, addr2, UnlockerAddress.String(), CTK123eth123), true}, // valid send with multiple coins
		{NewMsgLockedSend(addr1, addr2, UnlockerAddress.String(), CTK0), false},        // non positive coin
		{NewMsgLockedSend(addr1, addr2, UnlockerAddress.String(), CTK123eth0), false},  // non positive coin in multicoins
		{NewMsgLockedSend(emptyAddr, addr2, UnlockerAddress.String(), CTK123), false},  // empty from addr
		{NewMsgLockedSend(addr1, emptyAddr, UnlockerAddress.String(), CTK123), false},  // empty to addr
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			suite.Require().Nil(err)
		}
	}
}
