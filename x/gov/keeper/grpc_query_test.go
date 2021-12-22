package keeper_test

import (
	"fmt"

	"github.com/certikfoundation/shentu/v2/x/gov/types"
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (suite *KeeperTestSuite) TestQueryProposal() {
	ctx, queryClient := suite.ctx, suite.queryClient
	type proposal struct {
		title       string
		description string
		proposer    sdk.AccAddress
	}
	tests := []struct {
		name       string
		proposal   proposal
		proposalId uint64
		shouldPass bool
	}{
		{
			name: "Proposal ID exists",
			proposal: proposal{
				title:       "title0",
				description: "description0",
				proposer:    suite.address[0],
			},
			proposalId: 1,
			shouldPass: true,
		},
		{
			name: "Proposal ID does not exist",
			proposal: proposal{
				title:       "title1",
				description: "description1",
				proposer:    suite.address[0],
			},
			proposalId: 10,
			shouldPass: false,
		},
		{
			name: "Proposal ID can't be 0",
			proposal: proposal{
				title:       "title2",
				description: "description2",
				proposer:    suite.address[0],
			},
			proposalId: 0,
			shouldPass: false,
		},
	}

	for _, tc := range tests {
		textProposalContent := govtypes.NewTextProposal(tc.proposal.title, tc.proposal.description)

		// submit a new proposal
		_, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent, tc.proposal.proposer)
		suite.Require().NoError(err)
	}

	for _, tc := range tests {
		queryResponse, err := queryClient.Proposal(ctx.Context(), &types.QueryProposalRequest{ProposalId: tc.proposalId})
		if tc.shouldPass {
			suite.Require().NoError(err)
			suite.Require().Equal(tc.proposalId, queryResponse.Proposal.ProposalId)
		} else {
			suite.Require().Error(err)
		}
	}
}

func (suite *KeeperTestSuite) TestQueryProposals() {
	ctx, queryClient := suite.ctx, suite.queryClient
	type proposal struct {
		title       string
		description string
		proposer    sdk.AccAddress
	}
	events := []struct {
		proposal      proposal
		depositor     sdk.AccAddress
		voter         sdk.AccAddress
		fundedCoins   sdk.Coins
		depositAmount sdk.Coins
		deposit       bool
		vote          bool
	}{
		{
			proposal: proposal{
				// StatusValidatorVotingPeriod 3
				title:       "title",
				description: "description",
				proposer:    suite.validatorAccAddress,
			},
		},
		{
			proposal: proposal{
				// StatusDepositPeriod 1, insufficient deposit
				title:       "title0",
				description: "description0",
				proposer:    suite.address[0],
			},
			deposit:       true,
			depositor:     suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (100)*1e6)),
		},
		{
			proposal: proposal{
				// StatusDepositPeriod 1, insufficient funds
				title:       "title1",
				description: "description1",
				proposer:    suite.address[0],
			},
			deposit:       true,
			depositor:     suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (500)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (500)*1e6)),
		},
		{
			proposal: proposal{
				// StatusCertifierVotingPeriod 3
				title:       "title2",
				description: "description2",
				proposer:    suite.address[0],
			},
			deposit:       true,
			depositor:     suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
		},
		{
			proposal: proposal{
				// StatusCertifierVotingPeriod 3
				title:       "title3",
				description: "description3",
				proposer:    suite.address[0],
			},
			deposit:       true,
			depositor:     suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			vote:          true,
			voter:         suite.validatorAccAddress,
		},
		{
			proposal: proposal{
				// StatusCertifierVotingPeriod 3
				title:       "title4",
				description: "description4",
				proposer:    suite.address[0],
			},
			deposit:       true,
			depositor:     suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			vote:          true,
			// non validator/certifier
			voter: suite.address[1],
		},
	}

	for _, event := range events {
		textProposalContent := govtypes.NewTextProposal(event.proposal.title, event.proposal.description)
		// create/submit a new proposal
		proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent, event.proposal.proposer)
		suite.Require().NoError(err)

		// add staking coins to depositor
		suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, event.depositor, event.fundedCoins))

		if event.deposit {
			// deposit staked coins to get the proposal into voting period once it has exceeded minDeposit
			_, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, event.depositor, event.depositAmount)
			suite.Require().NoError(err)

		}

		if event.vote {
			// vote
			options := govtypes.NewNonSplitVoteOption(govtypes.OptionYes)
			vote := govtypes.NewVote(proposal.ProposalId, event.voter, options)
			voter, _ := sdk.AccAddressFromBech32(vote.Voter)
			err = suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, voter, options)
			suite.Require().NoError(err)

		}

		// emptying depositor for next set of events
		suite.app.BankKeeper.SendCoins(suite.ctx, event.depositor, suite.address[2], suite.app.BankKeeper.GetAllBalances(suite.ctx, event.depositor))

	}

	tests := []struct {
		name                    string
		voter                   string
		depositor               string
		filteredProposalsLength int
		proposalStatus          int32
		shouldPass              bool
	}{
		{
			name:                    "proposals with status 1",
			proposalStatus:          1,
			filteredProposalsLength: 2,
			shouldPass:              true,
		},
		{
			name:                    "proposals with status 3",
			proposalStatus:          3,
			filteredProposalsLength: 4,
			shouldPass:              true,
		},
		{
			name:                    "proposals with status 4",
			proposalStatus:          4,
			filteredProposalsLength: 3,
			shouldPass:              false,
		},
		{
			// none of the proposal has to go through security (certifier) voting
			name:                    "proposals with status 2",
			proposalStatus:          2,
			filteredProposalsLength: 0,
			shouldPass:              true,
		},
		{
			name:                    "proposals with specified voter(validator)",
			voter:                   suite.validatorAccAddress.String(),
			filteredProposalsLength: 1,
			shouldPass:              true,
		},
		{
			name:                    "proposals with specified voter(non-validator)",
			voter:                   suite.address[1].String(),
			filteredProposalsLength: 1,
			shouldPass:              true,
		},
		{
			name:                    "proposals with specified depositor",
			depositor:               suite.address[0].String(),
			filteredProposalsLength: 5,
			shouldPass:              true,
		},
		{
			name:                    "proposals with specified voter(non-validator)",
			voter:                   suite.address[2].String(),
			filteredProposalsLength: 1,
			shouldPass:              false,
		},
	}

	for _, tc := range tests {
		queryResponse, err := queryClient.Proposals(ctx.Context(), &types.QueryProposalsRequest{ProposalStatus: types.ProposalStatus(tc.proposalStatus), Voter: tc.voter, Depositor: tc.depositor})
		suite.Require().NoError(err)
		fmt.Println(tc.name)
		if tc.shouldPass {
			suite.Require().Equal(tc.filteredProposalsLength, len(queryResponse.Proposals))
		} else {
			suite.Require().NotEqual(tc.filteredProposalsLength, len(queryResponse.Proposals))
		}
	}
}

func (suite *KeeperTestSuite) TestQueryVote() {
	ctx, queryClient := suite.ctx, suite.queryClient
	type proposal struct {
		title       string
		description string
		proposer    sdk.AccAddress
		proposalId  int
	}
	tests := []struct {
		proposal      proposal
		voter         string
		name          string
		depositor     sdk.AccAddress
		fundedCoins   sdk.Coins
		depositAmount sdk.Coins
		voteOption    govtypes.VoteOption
		deposit       bool
		shouldPass    bool
	}{
		{
			name: "Proposal submitted by validator, vote yes",
			proposal: proposal{
				title:       "title",
				description: "description",
				proposer:    suite.validatorAccAddress,
				proposalId:  1,
			},
			voter:      suite.validatorAccAddress.String(),
			voteOption: govtypes.OptionYes,
			deposit:    false,
			shouldPass: true,
		},
		{
			name: "Proposal submitted by non-validator, vote yes",
			proposal: proposal{
				title:       "title0",
				description: "description0",
				proposer:    suite.address[0],
				proposalId:  2,
			},
			depositor:     suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			voter:         suite.validatorAccAddress.String(),
			voteOption:    govtypes.OptionYes,
			deposit:       true,
			shouldPass:    true,
		},
		{
			name: "Invalid voter address",
			proposal: proposal{
				title:       "title1",
				description: "description1",
				proposer:    suite.address[0],
				proposalId:  3,
			},
			depositor:     suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			voter:         "",
			voteOption:    govtypes.OptionYes,
			deposit:       true,
			shouldPass:    false,
		},
		{
			name: "Invalid proposal ID",
			proposal: proposal{
				title:       "title1",
				description: "description1",
				proposer:    suite.address[0],
				proposalId:  10,
			},
			depositor:     suite.address[0],
			fundedCoins:   sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			depositAmount: sdk.NewCoins(sdk.NewInt64Coin(suite.app.StakingKeeper.BondDenom(suite.ctx), (700)*1e6)),
			voter:         suite.validatorAccAddress.String(),
			voteOption:    govtypes.OptionYes,
			deposit:       true,
			shouldPass:    false,
		},
	}

	for _, tc := range tests {
		textProposalContent := govtypes.NewTextProposal(tc.proposal.title, tc.proposal.description)
		// create/submit a new proposal
		proposal, err := suite.app.GovKeeper.SubmitProposal(suite.ctx, textProposalContent, tc.proposal.proposer)
		suite.Require().NoError(err)

		// add staking coins to depositor
		suite.Require().NoError(sdksimapp.FundAccount(suite.app.BankKeeper, suite.ctx, tc.depositor, tc.fundedCoins))

		if tc.deposit {
			// deposit staked coins to get the proposal into voting period once it has exceeded minDeposit
			_, err := suite.app.GovKeeper.AddDeposit(suite.ctx, proposal.ProposalId, tc.depositor, tc.depositAmount)
			suite.Require().NoError(err)

		}
		voter, _ := sdk.AccAddressFromBech32(tc.voter)
		options := govtypes.NewNonSplitVoteOption(tc.voteOption)
		_ = suite.app.GovKeeper.AddVote(suite.ctx, proposal.ProposalId, voter, options)

		queryResponse, err := queryClient.Vote(ctx.Context(), &types.QueryVoteRequest{ProposalId: uint64(tc.proposal.proposalId), Voter: tc.voter})
		if tc.shouldPass {
			suite.Require().NoError(err)
			suite.Require().Equal(tc.voter, queryResponse.Vote.Voter)

		} else {
			suite.Require().Error(err)
		}

		// emptying depositor for next set of events
		suite.app.BankKeeper.SendCoins(suite.ctx, tc.depositor, suite.address[2], suite.app.BankKeeper.GetAllBalances(suite.ctx, tc.depositor))

	}
}
