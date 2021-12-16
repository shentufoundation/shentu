package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	"github.com/certikfoundation/shentu/v2/x/gov/keeper"
	"github.com/certikfoundation/shentu/v2/x/gov/types"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	address     []sdk.AccAddress
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.GovKeeper
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.GovKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.address = []sdk.AccAddress{acc1, acc2, acc3, acc4}
	// suite.keeper.SetCertifier(suite.ctx, types.NewCertifier(suite.address[0], "address1", suite.address[0], ""))
}

func (suite *KeeperTestSuite) TestKeeper_ProposeAndDeposit() {
	type proposal struct {
		title       string
		description string
	}

	tests := []struct {
		name               string
		proposal           proposal
		proposer           sdk.AccAddress
		fundedCoins        sdk.Coins
		depositAmount      sdk.Coins
		votingPeriodStatus bool
		reDeposit          bool
		err                bool
		shouldPass         bool
	}{
		{
			name: "New proposal, sufficient coins to start voting",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:           suite.address[0],
			fundedCoins:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount:      sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			votingPeriodStatus: true,
			reDeposit:          false,
			err:                false,
			shouldPass:         true,
		},
		{
			name: "New proposal, insufficient coins to start voting",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:           suite.address[0],
			fundedCoins:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (10)*1e6)),
			depositAmount:      sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (10)*1e6)),
			votingPeriodStatus: false,
			reDeposit:          false,
			err:                false,
			shouldPass:         false,
		},
		{
			name: "New proposal, deposit amount greater than funded coins",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:           suite.address[0],
			fundedCoins:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (600)*1e6)),
			depositAmount:      sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			votingPeriodStatus: false,
			reDeposit:          false,
			err:                true,
			shouldPass:         false,
		},
		{
			name: "New proposal, add more deposit after votingPeriod starts",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:           suite.address[0],
			fundedCoins:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (1500)*1e6)),
			depositAmount:      sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			votingPeriodStatus: true,
			reDeposit:          true,
			err:                true,
			shouldPass:         false,
		},
	}

	for _, tc := range tests {
		textProposalContent := govtypes.NewTextProposal(tc.proposal.title, tc.proposal.description)

		// create/submit a new proposal
		proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent, tc.proposer)
		suite.Require().NoError(err)

		// add staking coins to address 1
		suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, suite.address[1], tc.fundedCoins))

		// deposit staked coins to get the proposal into voting period once it has exceeded minDeposit
		votingPeriodActivated, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.address[1], tc.depositAmount)

		if tc.reDeposit {
			_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.address[1], tc.depositAmount)
		}

		if tc.shouldPass {
			suite.Require().NoError(err)
			suite.Require().Equal(tc.votingPeriodStatus, votingPeriodActivated)
		} else {
			if tc.err {
				suite.Require().Error(err)
			}
			suite.Require().Equal(tc.votingPeriodStatus, votingPeriodActivated)
		}
	}

}

func (suite *KeeperTestSuite) TestKeeper_DepositOperations() {
	type proposal struct {
		title       string
		description string
	}

	tests := []struct {
		name                 string
		proposal             proposal
		proposer             sdk.AccAddress
		fundedCoins          sdk.Coins
		depositAmount        sdk.Coins
		finalAmount          sdk.Coins
		testRefund           bool
		setInvalidProposalId bool
		shouldPass           bool
	}{
		{
			name: "Refund all deposits in a specific proposal",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:             suite.address[0],
			fundedCoins:          sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			finalAmount:          sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			testRefund:           true,
			setInvalidProposalId: false,
			shouldPass:           true,
		},
		{
			name: "Delete all deposits in a specific proposal",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:             suite.address[0],
			fundedCoins:          sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (600)*1e6)),
			finalAmount:          sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 100*1e6)),
			testRefund:           false,
			setInvalidProposalId: false,
			shouldPass:           true,
		},
		{
			name: "Refund: Invalid proposal id",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:             suite.address[0],
			fundedCoins:          sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			testRefund:           true,
			setInvalidProposalId: true,
			shouldPass:           false,
		},
		{
			name: "Delete: Invalid proposal id",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:             suite.address[0],
			fundedCoins:          sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			testRefund:           false,
			setInvalidProposalId: true,
			shouldPass:           false,
		},
	}

	for _, tc := range tests {
		textProposalContent := govtypes.NewTextProposal(tc.proposal.title, tc.proposal.description)

		// create/submit a new proposal
		proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent, tc.proposer)
		suite.Require().NoError(err)

		// add staking coins to address 1
		suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, suite.address[1], tc.fundedCoins))

		// deposit staked coins to get the proposal into voting period once it has exceeded minDeposit
		_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, suite.address[1], tc.depositAmount)
		suite.Require().NoError(err)
		if tc.setInvalidProposalId {
			proposal.ProposalId = proposal.ProposalId + 10
		}

		if tc.testRefund {
			err = suite.app.GovKeeper.RefundDepositsByProposalID(suite.ctx, proposal.ProposalId)
		} else {
			err = suite.app.GovKeeper.DeleteDepositsByProposalID(suite.ctx, proposal.ProposalId)
		}

		if tc.shouldPass {
			suite.Require().NoError(err)
			if tc.testRefund {
				suite.Require().Equal(tc.finalAmount, suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.address[1]))
			} else {
				suite.Require().Equal(tc.finalAmount, suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.address[1]))
			}
		} else {
			suite.Require().Error(err)
		}

		// emptying address 1 for next set of test cases
		suite.app.BankKeeper.SendCoins(suite.ctx, suite.address[1], suite.address[2], suite.app.BankKeeper.GetAllBalances(suite.ctx, suite.address[1]))
	}
}

func (suite *KeeperTestSuite) TestKeeper_Vote() {
	type proposal struct {
		title       string
		description string
	}

	// tests := []struct {
	// 	name               string
	// 	proposal           proposal
	// 	proposer           sdk.AccAddress
	// 	fundedCoins        sdk.Coins
	// 	depositAmount      sdk.Coins
	// 	votingPeriodStatus bool
	// 	err                bool
	// 	shouldPass         bool
	// }{}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
