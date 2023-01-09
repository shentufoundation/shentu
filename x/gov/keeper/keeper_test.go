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
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	stakingkeeper "github.com/shentufoundation/shentu/v2/x/staking/keeper"
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

	app                 *shentuapp.ShentuApp
	ctx                 sdk.Context
	keeper              keeper.Keeper
	address             []sdk.AccAddress
	queryClient         govtypes.QueryClient
	validatorAccAddress sdk.AccAddress
	msgServer           govtypes.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.GovKeeper
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	govtypes.RegisterQueryServer(queryHelper, suite.app.GovKeeper)
	suite.queryClient = govtypes.NewQueryClient(queryHelper)
	suite.address = []sdk.AccAddress{acc1, acc2, acc3, acc4}

	// suite.app.CertKeeper.SetCertifier(suite.ctx, certtypes.NewCertifier(suite.address[3], "", suite.address[3], ""))
	validatorAddress := sdk.ValAddress(suite.address[3])
	suite.validatorAccAddress = suite.address[3]
	suite.msgServer = keeper.NewMsgServerImpl(suite.app.GovKeeper)
	pks := shentuapp.CreateTestPubKeys(5)
	powers := []int64{1, 1, 1}
	cdc := sdksimapp.MakeTestEncodingConfig().Marshaler
	suite.app.StakingKeeper = stakingkeeper.NewKeeper(
		cdc,
		suite.app.GetKey(stakingtypes.StoreKey),
		suite.app.AccountKeeper,
		suite.app.BankKeeper,
		suite.app.GetSubspace(stakingtypes.ModuleName),
	)

	val1, err := stakingtypes.NewValidator(validatorAddress, pks[0], stakingtypes.Description{})
	suite.Require().NoError(err)
	val1.Status = stakingtypes.Bonded
	val1.DelegatorShares = sdk.OneDec()
	val1.Tokens = sdk.OneInt()
	suite.app.StakingKeeper.SetValidator(suite.ctx, val1)
	suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, val1)
	suite.app.StakingKeeper.SetNewValidatorByPowerIndex(suite.ctx, val1)
	suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, val1.GetOperator())
	_, _ = suite.app.StakingKeeper.Delegate(suite.ctx, suite.address[0], suite.app.StakingKeeper.TokensFromConsensusPower(suite.ctx, powers[0]), stakingtypes.Unbonded, val1, true)
}

// TODO: Add proposer in type proposal for all test cases
func (suite *KeeperTestSuite) TestKeeper_ProposeAndDeposit() {
	type proposal struct {
		title       string
		description string
	}

	tests := []struct {
		name               string
		proposal           proposal
		proposer           sdk.AccAddress
		depositor          sdk.AccAddress
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
			depositor:          suite.address[1],
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
			depositor:          suite.address[1],
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
			depositor:          suite.address[1],
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
			depositor:          suite.address[1],
			fundedCoins:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (1500)*1e6)),
			depositAmount:      sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			votingPeriodStatus: true,
			reDeposit:          true,
			err:                true,
			shouldPass:         false,
		},
		{
			name: "proposal submitted by validator is already active and doesn't require deposit",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:           suite.validatorAccAddress,
			depositor:          suite.address[1],
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
		proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent)
		suite.Require().NoError(err)
		// add staking coins to depositor
		suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, tc.depositor, tc.fundedCoins))

		// deposit staked coins to get the proposal into voting period once it has exceeded minDeposit
		votingPeriodActivated, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, tc.depositor, tc.depositAmount)

		if tc.reDeposit {
			_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, tc.depositor, tc.depositAmount)
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
		depositor            sdk.AccAddress
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
			depositor:            suite.address[1],
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
			depositor:            suite.address[1],
			fundedCoins:          sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount:        sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (600)*1e6)),
			finalAmount:          sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), 100*1e6)),
			testRefund:           false,
			setInvalidProposalId: false,
			shouldPass:           true,
		},
	}

	for _, tc := range tests {
		textProposalContent := govtypes.NewTextProposal(tc.proposal.title, tc.proposal.description)

		// create/submit a new proposal
		proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent)
		suite.Require().NoError(err)

		// add staking coins to depositor
		suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, tc.depositor, tc.fundedCoins))

		// deposit staked coins to get the proposal into voting period once it has exceeded minDeposit
		_, err = suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, tc.depositor, tc.depositAmount)
		suite.Require().NoError(err)

		if tc.setInvalidProposalId {
			proposal.ProposalId = proposal.ProposalId + 10
		}

		if tc.testRefund {
			suite.app.GovKeeper.RefundDeposits(suite.ctx, proposal.ProposalId)
		} else {
			suite.app.GovKeeper.DeleteDeposits(suite.ctx, proposal.ProposalId)
		}

		if tc.shouldPass {
			suite.Require().Equal(tc.finalAmount, suite.app.BankKeeper.GetAllBalances(suite.ctx, tc.depositor))
		}

		// emptying depositor for next set of test cases
		suite.app.BankKeeper.SendCoins(suite.ctx, tc.depositor, suite.address[2], suite.app.BankKeeper.GetAllBalances(suite.ctx, tc.depositor))
	}
}

func (suite *KeeperTestSuite) TestKeeper_Vote() {
	type proposal struct {
		title       string
		description string
	}

	tests := []struct {
		name               string
		proposal           proposal
		proposer           sdk.AccAddress
		depositor          sdk.AccAddress
		voter              sdk.AccAddress
		fundedCoins        sdk.Coins
		depositAmount      sdk.Coins
		expResults         map[govtypes.VoteOption]sdk.Dec
		votingPeriodStatus bool
		err                bool
		shouldPass         bool
	}{
		{
			name: "certifier/validator votes yes on a proposal, vote should be counted",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:      suite.address[0],
			depositor:     suite.address[1],
			voter:         suite.validatorAccAddress,
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			expResults: map[govtypes.VoteOption]sdk.Dec{
				govtypes.OptionYes:        sdk.OneDec(),
				govtypes.OptionAbstain:    sdk.ZeroDec(),
				govtypes.OptionNo:         sdk.ZeroDec(),
				govtypes.OptionNoWithVeto: sdk.ZeroDec(),
			},
			votingPeriodStatus: true,
			err:                false,
			shouldPass:         true,
		},
		{
			name: "non certifier/validator votes yes on a proposal, vote should not be counted",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:      suite.address[0],
			depositor:     suite.address[1],
			voter:         suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			expResults: map[govtypes.VoteOption]sdk.Dec{
				govtypes.OptionYes:        sdk.ZeroDec(),
				govtypes.OptionAbstain:    sdk.ZeroDec(),
				govtypes.OptionNo:         sdk.ZeroDec(),
				govtypes.OptionNoWithVeto: sdk.ZeroDec(),
			},
			votingPeriodStatus: true,
			err:                false,
			shouldPass:         true,
		},
		{
			name: "non certifier/validator vote, vote should not be counted, invalid expected results",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:      suite.address[0],
			depositor:     suite.address[1],
			voter:         suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			expResults: map[govtypes.VoteOption]sdk.Dec{
				govtypes.OptionYes:        sdk.OneDec(),
				govtypes.OptionAbstain:    sdk.ZeroDec(),
				govtypes.OptionNo:         sdk.ZeroDec(),
				govtypes.OptionNoWithVeto: sdk.ZeroDec(),
			},
			votingPeriodStatus: true,
			err:                false,
			shouldPass:         false,
		},
		{
			name: "voting period not started, so add vote should give error",
			proposal: proposal{
				title:       "title0",
				description: "description0",
			},
			proposer:      suite.address[0],
			depositor:     suite.address[1],
			voter:         suite.validatorAccAddress,
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (100)*1e6)),
			expResults: map[govtypes.VoteOption]sdk.Dec{
				govtypes.OptionYes:        sdk.OneDec(),
				govtypes.OptionAbstain:    sdk.ZeroDec(),
				govtypes.OptionNo:         sdk.ZeroDec(),
				govtypes.OptionNoWithVeto: sdk.ZeroDec(),
			},
			votingPeriodStatus: false,
			err:                true,
			shouldPass:         false,
		},
	}

	for _, tc := range tests {
		textProposalContent := govtypes.NewTextProposal(tc.proposal.title, tc.proposal.description)

		// create/submit a new proposal
		proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent)
		suite.Require().NoError(err)
		suite.Require().NotNil(proposal)

		// add staking coins to depositor
		suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, tc.depositor, tc.fundedCoins))

		// deposit staked coins to get the proposal into voting period once it has exceeded minDeposit
		votingPeriodStatus, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, tc.depositor, tc.depositAmount)
		if !tc.err {
			suite.Require().NoError(err)
			suite.Require().Equal(tc.votingPeriodStatus, votingPeriodStatus)
		}

		// vote
		options := govtypes.NewNonSplitVoteOption(govtypes.OptionYes)
		vote := govtypes.NewVote(proposal.ProposalId, tc.voter, options)
		voter, _ := sdk.AccAddressFromBech32(vote.Voter)
		err = suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, voter, options)
		if !tc.err {
			suite.Require().NoError(err)
		}

		// tally proposal
		_, _, results := suite.keeper.Tally(suite.ctx, proposal)
		if tc.shouldPass {
			suite.Require().Equal(govtypes.NewTallyResultFromMap(tc.expResults), results)
		} else {
			if tc.err {
				suite.Require().Error(err)
			} else {
				suite.Require().NotEqual(govtypes.NewTallyResultFromMap(tc.expResults), results)
			}
		}
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
