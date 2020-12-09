package keeper

import (
	"strconv"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/certikfoundation/shentu/x/cert/internal/types"
)

type QueryResRequestContent struct {
	RequestContentType     types.RequestContentType `json:"request_content_type"`
	RequestContentTypeName string                   `json:"request_content_type_name"`
	RequestContent         string                   `json:"request_content"`
}

func NewQueryResRequestContent(
	requestContentType types.RequestContentType,
	requestContentTypeName string,
	requestContent string,
) QueryResRequestContent {
	return QueryResRequestContent{
		RequestContentType:     requestContentType,
		RequestContentTypeName: requestContentTypeName,
		RequestContent:         requestContent,
	}
}

type QueryResCertificate struct {
	CertificateID      string                 `json:"certificate_id"`
	CertificateType    string                 `json:"certificate_type"`
	RequestContent     QueryResRequestContent `json:"request_content"`
	CertificateContent []types.KVPair         `json:"certificate_content"`
	Description        string                 `json:"description"`
	Certifier          string                 `json:"certifier"`
	TxHash             string                 `json:"txhash"`
}

func NewQueryResCertificate(
	certificateID string,
	certificateType string,
	requestContent types.RequestContent,
	certificateContent []types.KVPair,
	description string,
	certifier string,
	txhash string,
) QueryResCertificate {
	resRequestContent := NewQueryResRequestContent(
		requestContent.RequestContentType,
		requestContent.RequestContentType.String(),
		requestContent.RequestContent,
	)
	return QueryResCertificate{
		CertificateID:      certificateID,
		CertificateType:    certificateType,
		RequestContent:     resRequestContent,
		CertificateContent: certificateContent,
		Description:        description,
		Certifier:          certifier,
		TxHash:             txhash,
	}
}

func queryCertificate(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	if err := validatePathLength(path, 1); err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(path[0], 10, 64)
	if err != nil {
		return nil, err
	}
	certificate, err := keeper.GetCertificateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resCertificate := NewQueryResCertificate(
		strconv.FormatUint(certificate.ID(), 10),
		certificate.Type().String(),
		certificate.RequestContent(),
		certificate.FormattedCertificateContent(),
		certificate.Description(),
		certificate.Certifier().String(),
		certificate.TxHash(),
	)
	res, err := codec.MarshalJSONIndent(keeper.cdc, resCertificate)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, err
}

type QueryResCertificates struct {
	Total        uint64                `json:"total"`
	Certificates []QueryResCertificate `json:"certificates"`
}

func queryCertificates(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	if err := validatePathLength(path, 0); err != nil {
		return nil, err
	}
	var params types.QueryCertificatesParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	total, certificates, err := keeper.GetCertificatesFiltered(ctx, params)
	if err != nil {
		return nil, err
	}
	resCertificates := []QueryResCertificate{}
	for _, certificate := range certificates {
		resCertificate := NewQueryResCertificate(
			strconv.FormatUint(certificate.ID(), 10),
			certificate.Type().String(),
			certificate.RequestContent(),
			certificate.FormattedCertificateContent(),
			certificate.Description(),
			certificate.Certifier().String(),
			certificate.TxHash(),
		)
		resCertificates = append(resCertificates, resCertificate)
	}
	res, err := codec.MarshalJSONIndent(keeper.cdc, QueryResCertificates{Total: total, Certificates: resCertificates})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, err
}
