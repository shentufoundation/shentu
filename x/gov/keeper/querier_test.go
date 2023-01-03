package keeper_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

const custom = "custom"

func getQueriedParams(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier) (types.DepositParams, govtypes.VotingParams, types.TallyParams) {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, govtypes.QuerierRoute, govtypes.QueryParams, govtypes.ParamDeposit}, "/"),
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{govtypes.QueryParams, govtypes.ParamDeposit}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var depositParams types.DepositParams
	require.NoError(t, cdc.UnmarshalJSON(bz, &depositParams))

	query = abci.RequestQuery{
		Path: strings.Join([]string{custom, govtypes.QuerierRoute, govtypes.QueryParams, govtypes.ParamVoting}, "/"),
		Data: []byte{},
	}

	bz, err = querier(ctx, []string{govtypes.QueryParams, govtypes.ParamVoting}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var votingParams govtypes.VotingParams
	require.NoError(t, cdc.UnmarshalJSON(bz, &votingParams))

	query = abci.RequestQuery{
		Path: strings.Join([]string{custom, govtypes.QuerierRoute, govtypes.QueryParams, govtypes.ParamTallying}, "/"),
		Data: []byte{},
	}

	bz, err = querier(ctx, []string{govtypes.QueryParams, govtypes.ParamTallying}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var tallyParams types.TallyParams
	require.NoError(t, cdc.UnmarshalJSON(bz, &tallyParams))

	return depositParams, votingParams, tallyParams
}

func getQueriedProposals(
	t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier,
	depositor, voter sdk.AccAddress, status govtypes.ProposalStatus, page, limit int,
) []govtypes.Proposal {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, govtypes.QuerierRoute, govtypes.QueryProposals}, "/"),
		Data: cdc.MustMarshalJSON(govtypes.NewQueryProposalsParams(page, limit, status, voter, depositor)),
	}

	bz, err := querier(ctx, []string{govtypes.QueryProposals}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var proposals govtypes.Proposals
	require.NoError(t, cdc.UnmarshalJSON(bz, &proposals))

	return proposals
}

func getQueriedDeposit(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, proposalID uint64, depositor sdk.AccAddress) govtypes.Deposit {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, govtypes.QuerierRoute, govtypes.QueryDeposit}, "/"),
		Data: cdc.MustMarshalJSON(govtypes.NewQueryDepositParams(proposalID, depositor)),
	}

	bz, err := querier(ctx, []string{govtypes.QueryDeposit}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var deposit govtypes.Deposit
	require.NoError(t, cdc.UnmarshalJSON(bz, &deposit))

	return deposit
}

func getQueriedDeposits(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, proposalID uint64) []govtypes.Deposit {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, govtypes.QuerierRoute, govtypes.QueryDeposits}, "/"),
		Data: cdc.MustMarshalJSON(govtypes.NewQueryProposalParams(proposalID)),
	}

	bz, err := querier(ctx, []string{govtypes.QueryDeposits}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var deposits []govtypes.Deposit
	require.NoError(t, cdc.UnmarshalJSON(bz, &deposits))

	return deposits
}

func getQueriedVote(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier, proposalID uint64, voter sdk.AccAddress) govtypes.Vote {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, govtypes.QuerierRoute, govtypes.QueryVote}, "/"),
		Data: cdc.MustMarshalJSON(govtypes.NewQueryVoteParams(proposalID, voter)),
	}

	bz, err := querier(ctx, []string{govtypes.QueryVote}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var vote govtypes.Vote
	require.NoError(t, cdc.UnmarshalJSON(bz, &vote))

	return vote
}

func getQueriedVotes(t *testing.T, ctx sdk.Context, cdc *codec.LegacyAmino, querier sdk.Querier,
	proposalID uint64, page, limit int,
) []govtypes.Vote {
	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, govtypes.QuerierRoute, govtypes.QueryVote}, "/"),
		Data: cdc.MustMarshalJSON(govtypes.NewQueryProposalVotesParams(proposalID, page, limit)),
	}

	bz, err := querier(ctx, []string{govtypes.QueryVotes}, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	var votes []govtypes.Vote
	require.NoError(t, cdc.UnmarshalJSON(bz, &votes))

	return votes
}

func TestQueries(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	legacyQuerierCdc := app.LegacyAmino()
	querier := keeper.NewQuerier(app.GovKeeper, legacyQuerierCdc)

	TestAddrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(1000000000000))

	SortAddresses(TestAddrs)

	oneCoins := sdk.NewCoins(sdk.NewInt64Coin("uctk", 1))
	consCoins := sdk.NewCoins(sdk.NewInt64Coin("uctk", 10))
	fmt.Println(TestAddrs[0].String(), TestAddrs[1].String())
	tp := TestProposal

	depositParams, _, _ := getQueriedParams(t, ctx, legacyQuerierCdc, querier)

	// TestAddrs[0] proposes (and deposits) proposals #1 and #2
	proposal1, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	deposit1 := govtypes.NewDeposit(proposal1.ProposalId, TestAddrs[0], oneCoins)
	depositer1, err := sdk.AccAddressFromBech32(deposit1.Depositor)
	require.NoError(t, err)
	_, err = app.GovKeeper.AddDeposit(ctx, deposit1.ProposalId, depositer1, deposit1.Amount)
	require.NoError(t, err)

	proposal1.TotalDeposit = proposal1.TotalDeposit.Add(deposit1.Amount...)

	proposal2, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	deposit2 := govtypes.NewDeposit(proposal2.ProposalId, TestAddrs[0], consCoins)
	depositer2, err := sdk.AccAddressFromBech32(deposit2.Depositor)
	require.NoError(t, err)
	_, err = app.GovKeeper.AddDeposit(ctx, deposit2.ProposalId, depositer2, deposit2.Amount)
	require.NoError(t, err)

	proposal2.TotalDeposit = proposal2.TotalDeposit.Add(deposit2.Amount...)

	// TestAddrs[1] proposes (and deposits) on proposal #3
	proposal3, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	deposit3 := govtypes.NewDeposit(proposal3.ProposalId, TestAddrs[1], oneCoins)
	depositer3, err := sdk.AccAddressFromBech32(deposit3.Depositor)
	require.NoError(t, err)

	_, err = app.GovKeeper.AddDeposit(ctx, deposit3.ProposalId, depositer3, deposit3.Amount)
	require.NoError(t, err)

	proposal3.TotalDeposit = proposal3.TotalDeposit.Add(deposit3.Amount...)

	// TestAddrs[1] deposits on proposals #2 & #3
	deposit4 := govtypes.NewDeposit(proposal2.ProposalId, TestAddrs[1], depositParams.MinDeposit)
	depositer4, err := sdk.AccAddressFromBech32(deposit4.Depositor)
	require.NoError(t, err)
	_, err = app.GovKeeper.AddDeposit(ctx, deposit4.ProposalId, depositer4, deposit4.Amount)
	require.NoError(t, err)

	proposal2.TotalDeposit = proposal2.TotalDeposit.Add(deposit4.Amount...)
	proposal2.Status = govtypes.StatusVotingPeriod
	proposal2.VotingEndTime = proposal2.VotingEndTime.Add(govtypes.DefaultPeriod)

	deposit5 := govtypes.NewDeposit(proposal3.ProposalId, TestAddrs[1], depositParams.MinDeposit)
	depositer5, err := sdk.AccAddressFromBech32(deposit5.Depositor)
	require.NoError(t, err)
	_, err = app.GovKeeper.AddDeposit(ctx, deposit5.ProposalId, depositer5, deposit5.Amount)
	require.NoError(t, err)

	proposal3.TotalDeposit = proposal3.TotalDeposit.Add(deposit5.Amount...)
	proposal3.Status = govtypes.StatusVotingPeriod
	proposal3.VotingEndTime = proposal3.VotingEndTime.Add(govtypes.DefaultPeriod)
	// total deposit of TestAddrs[1] on proposal #3 is worth the proposal deposit + individual deposit
	deposit5.Amount = deposit5.Amount.Add(deposit3.Amount...)

	// check deposits on proposal1 match individual deposits

	deposits := getQueriedDeposits(t, ctx, legacyQuerierCdc, querier, proposal1.ProposalId)
	require.Len(t, deposits, 1)
	require.Equal(t, deposit1, deposits[0])

	deposit := getQueriedDeposit(t, ctx, legacyQuerierCdc, querier, proposal1.ProposalId, TestAddrs[0])
	require.Equal(t, deposit1, deposit)

	// check deposits on proposal2 match individual deposits
	deposits = getQueriedDeposits(t, ctx, legacyQuerierCdc, querier, proposal2.ProposalId)
	require.Len(t, deposits, 2)
	// NOTE order of deposits is determined by the addresses
	require.Equal(t, deposit2, deposits[0])
	require.Equal(t, deposit4, deposits[1])

	// check deposits on proposal3 match individual deposits
	deposits = getQueriedDeposits(t, ctx, legacyQuerierCdc, querier, proposal3.ProposalId)
	require.Len(t, deposits, 1)
	require.Equal(t, deposit5, deposits[0])

	deposit = getQueriedDeposit(t, ctx, legacyQuerierCdc, querier, proposal3.ProposalId, TestAddrs[1])
	require.Equal(t, deposit5, deposit)

	// Only proposal #1 should be in types.Deposit Period
	proposals := getQueriedProposals(t, ctx, legacyQuerierCdc, querier, nil, nil, govtypes.StatusDepositPeriod, 1, 0)
	require.Len(t, proposals, 1)
	require.Equal(t, proposal1, proposals[0])

	// Only proposals #2 and #3 should be in Voting Period
	proposals = getQueriedProposals(t, ctx, legacyQuerierCdc, querier, nil, nil, govtypes.StatusVotingPeriod, 1, 0)
	require.Len(t, proposals, 2)
	proposal2.DepositEndTime = proposals[0].DepositEndTime
	proposal3.DepositEndTime = proposals[1].DepositEndTime
	require.Equal(t, proposal2, proposals[0])
	require.Equal(t, proposal3, proposals[1])

	// Addrs[0] votes on proposals #2 & #3
	vote1 := govtypes.NewVote(proposal2.ProposalId, TestAddrs[0], govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
	vote2 := govtypes.NewVote(proposal3.ProposalId, TestAddrs[0], govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
	app.GovKeeper.SetVote(ctx, vote1)
	app.GovKeeper.SetVote(ctx, vote2)

	// Addrs[1] votes on proposal #3
	vote3 := govtypes.NewVote(proposal3.ProposalId, TestAddrs[1], govtypes.NewNonSplitVoteOption(govtypes.OptionYes))
	app.GovKeeper.SetVote(ctx, vote3)

	// Test query voted by TestAddrs[0]
	proposals = getQueriedProposals(t, ctx, legacyQuerierCdc, querier, nil, TestAddrs[0], govtypes.StatusNil, 1, 0)
	require.Equal(t, proposal2, proposals[0])
	require.Equal(t, proposal3, proposals[1])

	// Test query votes on types.Proposal 2
	votes := getQueriedVotes(t, ctx, legacyQuerierCdc, querier, proposal2.ProposalId, 1, 0)
	require.Len(t, votes, 1)
	checkEqualVotes(t, vote1, votes[0])

	vote := getQueriedVote(t, ctx, legacyQuerierCdc, querier, proposal2.ProposalId, TestAddrs[0])
	checkEqualVotes(t, vote1, vote)

	// Test query votes on types.Proposal 3
	votes = getQueriedVotes(t, ctx, legacyQuerierCdc, querier, proposal3.ProposalId, 1, 0)
	require.Len(t, votes, 2)
	checkEqualVotes(t, vote2, votes[0])
	checkEqualVotes(t, vote3, votes[1])

	// Test query all proposals
	proposals = getQueriedProposals(t, ctx, legacyQuerierCdc, querier, nil, nil, govtypes.StatusNil, 1, 0)
	require.Equal(t, proposal1, proposals[0])
	require.Equal(t, proposal2, proposals[1])
	require.Equal(t, proposal3, proposals[2])

	// Test query voted by TestAddrs[1]
	proposals = getQueriedProposals(t, ctx, legacyQuerierCdc, querier, nil, TestAddrs[1], govtypes.StatusNil, 1, 0)
	require.Equal(t, proposal3.ProposalId, proposals[0].ProposalId)

	// Test query deposited by TestAddrs[0]
	proposals = getQueriedProposals(t, ctx, legacyQuerierCdc, querier, TestAddrs[0], nil, govtypes.StatusNil, 1, 0)
	require.Equal(t, proposal1.ProposalId, proposals[0].ProposalId)

	// Test query deposited by addr2
	proposals = getQueriedProposals(t, ctx, legacyQuerierCdc, querier, TestAddrs[1], nil, govtypes.StatusNil, 1, 0)
	require.Equal(t, proposal2.ProposalId, proposals[0].ProposalId)
	require.Equal(t, proposal3.ProposalId, proposals[1].ProposalId)

	// Test query voted AND deposited by addr1
	proposals = getQueriedProposals(t, ctx, legacyQuerierCdc, querier, TestAddrs[0], TestAddrs[0], govtypes.StatusNil, 1, 0)
	require.Equal(t, proposal2.ProposalId, proposals[0].ProposalId)
}

func TestPaginatedVotesQuery(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	legacyQuerierCdc := app.LegacyAmino()

	proposal := govtypes.Proposal{
		ProposalId: 100,
		Status:     govtypes.StatusVotingPeriod,
	}

	app.GovKeeper.SetProposal(ctx, proposal)

	votes := make([]govtypes.Vote, 20)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	addrMap := make(map[string]struct{})
	genAddr := func() string {
		addr := make(sdk.AccAddress, 20)
		for {
			random.Read(addr)
			addrStr := addr.String()
			if _, ok := addrMap[addrStr]; !ok {
				addrMap[addrStr] = struct{}{}
				return addrStr
			}
		}
	}
	for i := range votes {
		vote := govtypes.Vote{
			ProposalId: proposal.ProposalId,
			Voter:      genAddr(),
			Options:    govtypes.NewNonSplitVoteOption(govtypes.OptionYes),
		}
		votes[i] = vote
		app.GovKeeper.SetVote(ctx, vote)
	}

	querier := keeper.NewQuerier(app.GovKeeper, legacyQuerierCdc)

	// keeper preserves consistent order for each query, but this is not the insertion order
	all := getQueriedVotes(t, ctx, legacyQuerierCdc, querier, proposal.ProposalId, 1, 0)
	require.Equal(t, len(all), len(votes))

	type testCase struct {
		description string
		page        int
		limit       int
		votes       []govtypes.Vote
	}
	for _, tc := range []testCase{
		{
			description: "SkipAll",
			page:        2,
			limit:       len(all),
		},
		{
			description: "GetFirstChunk",
			page:        1,
			limit:       10,
			votes:       all[:10],
		},
		{
			description: "GetSecondsChunk",
			page:        2,
			limit:       10,
			votes:       all[10:],
		},
		{
			description: "InvalidPage",
			page:        -1,
		},
	} {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			votes := getQueriedVotes(t, ctx, legacyQuerierCdc, querier, proposal.ProposalId, tc.page, tc.limit)
			require.Equal(t, len(tc.votes), len(votes))
			for i := range votes {
				require.Equal(t, tc.votes[i], votes[i])
			}
		})
	}
}

// checkEqualVotes checks that two votes are equal, without taking into account
// graceful fallback for `Option`.
// When querying, the keeper populates the `vote.Option` field when there's
// only 1 vote, this function checks equality of structs while skipping that
// field.
func checkEqualVotes(t *testing.T, vote1, vote2 govtypes.Vote) {
	require.Equal(t, vote1.Options, vote2.Options)
	require.Equal(t, vote1.Voter, vote2.Voter)
	require.Equal(t, vote1.ProposalId, vote2.ProposalId)
}
