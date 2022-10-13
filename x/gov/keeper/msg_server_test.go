package keeper_test

import (
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (suite *KeeperTestSuite) TestVoteWeightedReq() {
	proposer := suite.address[0]
	coins := sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 1e10))
	// add staking coins to depositor
	suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, proposer, coins))

	msg, err := govtypes.NewMsgSubmitProposal(
		govtypes.NewTextProposal("title0", "description0"),
		coins,
		proposer,
	)
	suite.Require().NoError(err)

	res, err := suite.msgServer.SubmitProposal(sdk.WrapSDKContext(suite.ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res.ProposalId)

	tests := []struct {
		proposalID uint64
		voterAddr  sdk.AccAddress
		options    govtypes.WeightedVoteOptions
		expectPass bool
	}{
		{res.ProposalId, proposer, govtypes.NewNonSplitVoteOption(govtypes.OptionYes), true},
	}
	for _, tc := range tests {
		voteReq := govtypes.NewMsgVoteWeighted(tc.voterAddr, tc.proposalID, tc.options)
		_, _ = suite.msgServer.VoteWeighted(sdk.WrapSDKContext(suite.ctx), voteReq)
	}

}
