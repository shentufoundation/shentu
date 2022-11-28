package keeper_test

import (
	"strings"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (suite *KeeperTestSuite) TestVoteWeightedReq() {
	proposer := suite.address[0]

	stakingCoins := sdk.NewCoins(sdk.NewCoin("uctk", sdk.NewInt(1e10)))
	minDeposit := suite.app.GovKeeper.GetDepositParams(suite.ctx).MinDeposit
	inactiveCoins := sdk.NewCoins(sdk.NewCoin("uctk", sdk.NewInt(1)))

	// add staking coins to depositor
	suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, proposer, stakingCoins))

	proposalContent := govtypes.NewTextProposal("title0", "description0")
	// active proposal
	msg, err := govtypes.NewMsgSubmitProposal(
		proposalContent,
		minDeposit,
		proposer,
	)
	suite.Require().NoError(err)
	res, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.ProposalId)

	// inactive proposal
	inactiveMsg, err := govtypes.NewMsgSubmitProposal(
		proposalContent,
		inactiveCoins,
		proposer,
	)
	suite.Require().NoError(err)
	inactiveRes, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), inactiveMsg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.ProposalId)

	cases := []struct {
		name       string
		proposalId uint64
		voterAddr  sdk.AccAddress
		options    govtypes.WeightedVoteOptions
		expectPass bool
		expErrMsg  string
	}{
		{
			"all good",
			res.ProposalId,
			proposer,
			govtypes.NewNonSplitVoteOption(govtypes.OptionYes),
			true,
			"",
		},
		{"voter error",
			res.ProposalId,
			sdk.AccAddress(strings.Repeat("a", 300)),
			govtypes.NewNonSplitVoteOption(govtypes.OptionYes),
			false,
			"address max length is 255",
		},
		{"vote on inactive proposal",
			inactiveRes.ProposalId,
			proposer,
			govtypes.NewNonSplitVoteOption(govtypes.OptionYes),
			false,
			"inactive proposal",
		},
	}

	for _, tc := range cases {
		suite.Run(tc.name, func() {
			voteReq := govtypes.NewMsgVoteWeighted(tc.voterAddr, tc.proposalId, tc.options)
			_, err = suite.msgServer.VoteWeighted(sdk.WrapSDKContext(suite.ctx), voteReq)
			if tc.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
			}
		})
	}
}
