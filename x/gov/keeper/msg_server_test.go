package keeper_test

//
//import (
//	"strings"
//
//	"github.com/cosmos/cosmos-sdk/testutil/testdata"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
//	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
//
//	"github.com/shentufoundation/shentu/v2/common"
//)
//
//func (suite *KeeperTestSuite) TestSubmitProposalReq() {
//	govAcct := suite.app.GovKeeper.GetGovernanceAccount(suite.ctx).GetAddress()
//	addrs := suite.addrs
//	proposer := addrs[0]
//
//	coins := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(100)))
//	initialDeposit := coins
//	minDeposit := suite.app.GovKeeper.GetDepositParams(suite.ctx).MinDeposit
//	bankMsg := &banktypes.MsgSend{
//		FromAddress: govAcct.String(),
//		ToAddress:   proposer.String(),
//		Amount:      coins,
//	}
//
//	cases := map[string]struct {
//		preRun    func() (*v1.MsgSubmitProposal, error)
//		expErr    bool
//		expErrMsg string
//	}{
//		"metadata too long": {
//			preRun: func() (*v1.MsgSubmitProposal, error) {
//				return v1.NewMsgSubmitProposal(
//					[]sdk.Msg{bankMsg},
//					initialDeposit,
//					proposer.String(),
//					strings.Repeat("1", 300),
//				)
//			},
//			expErr:    true,
//			expErrMsg: "metadata too long",
//		},
//		"many signers": {
//			preRun: func() (*v1.MsgSubmitProposal, error) {
//				return v1.NewMsgSubmitProposal(
//					[]sdk.Msg{testdata.NewTestMsg(govAcct, addrs[0])},
//					initialDeposit,
//					proposer.String(),
//					"",
//				)
//			},
//			expErr:    true,
//			expErrMsg: "expected gov account as only signer for proposal message",
//		},
//		"signer isn't gov account": {
//			preRun: func() (*v1.MsgSubmitProposal, error) {
//				return v1.NewMsgSubmitProposal(
//					[]sdk.Msg{testdata.NewTestMsg(addrs[0])},
//					initialDeposit,
//					proposer.String(),
//					"",
//				)
//			},
//			expErr:    true,
//			expErrMsg: "expected gov account as only signer for proposal message",
//		},
//		"invalid msg handler": {
//			preRun: func() (*v1.MsgSubmitProposal, error) {
//				return v1.NewMsgSubmitProposal(
//					[]sdk.Msg{testdata.NewTestMsg(govAcct)},
//					initialDeposit,
//					proposer.String(),
//					"",
//				)
//			},
//			expErr:    true,
//			expErrMsg: "proposal message not recognized by router",
//		},
//		"all good": {
//			preRun: func() (*v1.MsgSubmitProposal, error) {
//				return v1.NewMsgSubmitProposal(
//					[]sdk.Msg{bankMsg},
//					initialDeposit,
//					proposer.String(),
//					"",
//				)
//			},
//			expErr: false,
//		},
//		"all good with min deposit": {
//			preRun: func() (*v1.MsgSubmitProposal, error) {
//				return v1.NewMsgSubmitProposal(
//					[]sdk.Msg{bankMsg},
//					minDeposit,
//					proposer.String(),
//					"",
//				)
//			},
//			expErr: false,
//		},
//	}
//
//	for name, tc := range cases {
//		suite.Run(name, func() {
//			msg, err := tc.preRun()
//			suite.Require().NoError(err)
//			res, err := suite.msgSrvr.SubmitProposal(suite.ctx, msg)
//			if tc.expErr {
//				suite.Require().Error(err)
//				suite.Require().Contains(err.Error(), tc.expErrMsg)
//			} else {
//				suite.Require().NoError(err)
//				suite.Require().NotNil(res.ProposalId)
//			}
//		})
//	}
//}
