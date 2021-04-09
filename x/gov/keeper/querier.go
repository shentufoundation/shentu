package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/gov/types"
)

// NewQuerier creates a new gov Querier instance.
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case govtypes.QueryParams:
			return queryParams(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case govtypes.QueryProposals:
			return queryProposals(ctx, path[1:], req, keeper, legacyQuerierCdc)

		case govtypes.QueryProposal:
			return queryProposal(ctx, req, keeper, legacyQuerierCdc)

		case govtypes.QueryDeposits:
			return queryDeposits(ctx, req, keeper, legacyQuerierCdc)

		case govtypes.QueryDeposit:
			return queryDeposit(ctx, req, keeper, legacyQuerierCdc)

		case govtypes.QueryVotes:
			return queryVotes(ctx, req, keeper, legacyQuerierCdc)

		case govtypes.QueryVote:
			return queryVote(ctx, req, keeper, legacyQuerierCdc)

		case govtypes.QueryTally:
			return queryTally(ctx, req, keeper, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	switch path[0] {
	case govtypes.ParamDeposit:
		bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, keeper.GetDepositParams(ctx))
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil

	case govtypes.ParamVoting:
		bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, keeper.GetVotingParams(ctx))
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil

	case govtypes.ParamTallying:
		bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, keeper.GetTallyParams(ctx))
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil

	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "%s is not a valid query request path", req.Path)
	}
}

func queryProposal(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params govtypes.QueryProposalParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	proposal, ok := keeper.GetProposal(ctx, params.ProposalID)
	if !ok {
		return nil, sdkerrors.Wrapf(govtypes.ErrUnknownProposal, "%d", params.ProposalID)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, proposal)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDeposit(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params govtypes.QueryDepositParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	deposit, _ := keeper.GetDeposit(ctx, params.ProposalID, params.Depositor)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, deposit)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params govtypes.QueryVoteParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	vote, _ := keeper.GetVote(ctx, params.ProposalID, params.Voter)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, vote)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDeposits(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params govtypes.QueryProposalParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	deposits := keeper.GetDeposits(ctx, params.ProposalID)
	if deposits == nil {
		deposits = types.Deposits{}
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, deposits)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTally(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params govtypes.QueryProposalParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	proposalID := params.ProposalID

	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return nil, sdkerrors.Wrapf(govtypes.ErrUnknownProposal, "%d", proposalID)
	}

	var tallyResult govtypes.TallyResult

	switch {
	case proposal.Status == types.StatusDepositPeriod:
		tallyResult = govtypes.EmptyTallyResult()

	case proposal.Status == types.StatusPassed || proposal.Status == types.StatusRejected:
		tallyResult = proposal.FinalTallyResult

	case proposal.Status == types.StatusCertifierVotingPeriod:
		_, _, tallyResult = SecurityTally(ctx, keeper, proposal)

	default:
		// proposal is in voting period
		_, _, tallyResult = Tally(ctx, keeper, proposal)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, tallyResult)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params govtypes.QueryProposalVotesParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	votes := keeper.GetVotesPaginated(ctx, params.ProposalID, uint(params.Page), uint(params.Limit))

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, votes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

type QueryResProposals struct {
	Total     int              `json:"total"`
	Proposals []types.Proposal `json:"proposals"`
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (q QueryResProposals) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, x := range q.Proposals {
		err := x.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}

func queryProposals(ctx sdk.Context, _ []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryProposalsParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	proposals := keeper.GetProposalsFiltered(ctx, params)
	if proposals == nil {
		proposals = types.Proposals{}
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, QueryResProposals{Total: len(proposals), Proposals: proposals})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
