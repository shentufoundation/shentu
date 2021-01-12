package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/certikfoundation/shentu/x/gov/types"
)

// NewQuerier creates a new gov Querier instance.
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case govTypes.QueryParams:
			return queryParams(ctx, path[1:], req, keeper)

		case govTypes.QueryProposals:
			return queryProposals(ctx, path[1:], req, keeper)

		case govTypes.QueryProposal:
			return queryProposal(ctx, req, keeper)

		case govTypes.QueryDeposits:
			return queryDeposits(ctx, req, keeper)

		case govTypes.QueryDeposit:
			return queryDeposit(ctx, req, keeper)

		case govTypes.QueryVotes:
			return queryVotes(ctx, req, keeper)

		case govTypes.QueryVote:
			return queryVote(ctx, req, keeper)

		case govTypes.QueryTally:
			return queryTally(ctx, req, keeper)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	switch path[0] {
	case govTypes.ParamDeposit:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetDepositParams(ctx))
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil

	case govTypes.ParamVoting:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetVotingParams(ctx))
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil

	case govTypes.ParamTallying:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetTallyParams(ctx))
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil

	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "%s is not a valid query request path", req.Path)
	}
}

func queryProposal(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params govTypes.QueryProposalParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	proposal, ok := keeper.GetProposal(ctx, params.ProposalID)
	if !ok {
		return nil, sdkerrors.Wrapf(govTypes.ErrUnknownProposal, "%d", params.ProposalID)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, proposal)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDeposit(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params govTypes.QueryDepositParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	deposit, _ := keeper.GetDeposit(ctx, params.ProposalID, params.Depositor)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, deposit)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params govTypes.QueryVoteParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	vote, _ := keeper.GetVote(ctx, params.ProposalID, params.Voter)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, vote)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryDeposits(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params govTypes.QueryProposalParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	deposits := keeper.GetDeposits(ctx, params.ProposalID)
	if deposits == nil {
		deposits = types.Deposits{}
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, deposits)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTally(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params govTypes.QueryProposalParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	proposalID := params.ProposalID

	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return nil, sdkerrors.Wrapf(govTypes.ErrUnknownProposal, "%d", proposalID)
	}

	var tallyResult govTypes.TallyResult

	switch {
	case proposal.Status == types.StatusDepositPeriod:
		tallyResult = govTypes.EmptyTallyResult()

	case proposal.Status == types.StatusPassed || proposal.Status == types.StatusRejected:
		tallyResult = proposal.FinalTallyResult

	case proposal.Status == types.StatusCertifierVotingPeriod:
		_, _, tallyResult = SecurityTally(ctx, keeper, proposal)

	default:
		// proposal is in voting period
		_, _, tallyResult = Tally(ctx, keeper, proposal)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, tallyResult)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryVotes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params govTypes.QueryProposalVotesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	votes := keeper.GetVotesPaginated(ctx, params.ProposalID, uint(params.Page), uint(params.Limit))

	bz, err := codec.MarshalJSONIndent(keeper.cdc, votes)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryProposals(ctx sdk.Context, _ []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryProposalsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	proposals := keeper.GetProposalsFiltered(ctx, params)
	if proposals == nil {
		proposals = types.Proposals{}
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, proposals)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
