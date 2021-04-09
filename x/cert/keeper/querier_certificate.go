package keeper

import (
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/cert/types"
)

func queryCertificate(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	certificateID, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, err
	}

	certificate, err := keeper.GetCertificateByID(ctx, certificateID)
	if err != nil {
		return nil, err
	}
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

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (q QueryResCertificates) UnpackInterfaces(unpacker codecTypes.AnyUnpacker) error {
	for _, x := range q.Certificates {
		err := x.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
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
