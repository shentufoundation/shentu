package keeper

import (
	"github.com/irisnet/irismod/modules/nft/keeper"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/nft/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryCertificate:
			return queryCertificate(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryCertificates:
			return queryCertificates(ctx, path[1:], req, k, legacyQuerierCdc)
		default:
			return keeper.NewQuerier(k.Keeper, legacyQuerierCdc)(ctx, path, req)
		}
	}
}

func validatePathLength(path []string, length int) error {
	if len(path) != length {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting %d args. Found %d.", length, len(path))
	}
	return nil
}

func queryCertificate(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if err := validatePathLength(path, 2); err != nil {
		return nil, err
	}

	certNFT, err := keeper.GetNFT(ctx, path[0], path[1])
	if err != nil {
		return nil, err
	}
	certificate := keeper.UnmarshalCertificate(ctx, certNFT.GetData())

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, certificate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, err
}

type QueryResCertificates struct {
	Total        uint64              `json:"total"`
	Certificates []types.Certificate `json:"certificates"`
}

func queryCertificates(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}
	var params types.QueryCertificatesParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	total, certificates, err := keeper.GetCertificatesFiltered(ctx, params)
	if err != nil {
		return nil, err
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, QueryResCertificates{Total: total, Certificates: certificates})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, err
}
