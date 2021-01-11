package keeper

import (
	"strings"

	"github.com/tmthrgd/go-hex"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/txs/payload"

	"github.com/certikfoundation/shentu/x/cvm/types"
)

// NewQuerier is the module level router for state queries.
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryCode:
			return queryCode(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case types.QueryStorage:
			return queryStorage(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case types.QueryAbi:
			return queryAbi(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case types.QueryAddrMeta:
			return queryAddrMeta(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case types.QueryMeta:
			return queryMeta(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case types.QueryView:
			return queryView(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case types.QueryAccount:
			return queryAccount(ctx, path[1:], req, keeper, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown cvm query endpoint "+strings.Join(path, "/"))
		}
	}
}

func queryView(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if len(path) != 2 {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting 2 args. Found %d.", len(path))
	}

	// caller is an empty address for a view command
	callerString := path[0]
	var caller sdk.AccAddress
	if callerString != "" {
		caller, err = sdk.AccAddressFromBech32(path[0])
		if err != nil {
			panic("could not parse caller address " + path[0])
		}
	} else {
		caller = make([]byte, 20)
	}

	callee, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		panic("could not parse address " + path[1])
	}

	value, err := keeper.Tx(ctx, caller, callee, 0, req.Data, []*payload.ContractMeta{}, true, false, false)

	if err != nil {
		panic("failed to get storage at address " + path[0])
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryResView{Ret: value})
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryCode(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if len(path) != 1 {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting 1 args. Found %d.", len(path))
	}

	addr, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		panic("Could not parse address " + path[0])
	}

	account := keeper.getAccount(ctx, crypto.MustAddressFromBytes(addr))
	if account == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, path[0])
	}

	if len(account.EVMCode) != 0 {
		res, err = codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryResCode{Code: account.EVMCode})
		if err != nil {
			panic("could not marshal result to JSON")
		}
	} else {
		res, err = codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryResCode{Code: account.WASMCode})
		if err != nil {
			panic("could not marshal result to JSON")
		}
	}

	return res, nil
}

func queryStorage(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if len(path) != 2 {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting 2 args. Found %d.", len(path))
	}

	addr, err1 := sdk.AccAddressFromBech32(path[0])
	if err1 != nil {
		panic("Could not parse address " + path[0])
	}

	i, ok := sdk.NewIntFromString(path[1])
	if !ok {
		panic("Could not parse key " + path[1])
	}
	var key binary.Word256
	bytes := i.BigInt().Bytes()
	if len(bytes) > len(key) {
		panic("key size is too large " + path[1])
	}
	copy(key[:], bytes)

	value, err2 := keeper.GetStorage(ctx, crypto.MustAddressFromBytes(addr), key)
	if err2 != nil {
		panic("failed to get storage at address " + path[0] + " key " + path[1])
	}

	res, err3 := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryResStorage{Value: value})
	if err3 != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryAbi(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if len(path) != 1 {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting 1 args. Found %d.", len(path))
	}

	addr, err1 := sdk.AccAddressFromBech32(path[0])
	if err1 != nil {
		panic("Could not parse address " + path[0])
	}

	abi := keeper.GetAbi(ctx, crypto.MustAddressFromBytes(addr))
	res, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryResAbi{Abi: abi})
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryAddrMeta(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if len(path) != 1 {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting 1 args. Found %d.", len(path))
	}

	addr, err1 := sdk.AccAddressFromBech32(path[0])
	if err1 != nil {
		panic("Could not parse address " + path[0])
	}

	contMeta, err := keeper.getAddrMeta(ctx, crypto.MustAddressFromBytes(addr))
	if err != nil {
		panic("could not get metadata for the address")
	}

	var resString string
	for i := range contMeta {
		contMeta[i].MetadataHash = []byte(contMeta[i].MetadataHash.String())
		contMeta[i].CodeHash = []byte(contMeta[i].CodeHash.String())
		resString += contMeta[i].String()
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryResAddrMeta{Metahash: resString})
	if err != nil {
		panic("could not marshal result to JSON")
	}
	return
}

func queryMeta(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if len(path) != 1 {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting 1 args. Found %d.", len(path))
	}

	hash, err := hex.DecodeString(path[0])
	if err != nil {
		panic(err)
	}
	var metahash acmstate.MetadataHash
	copy(metahash[:], hash)

	meta, err := keeper.getMeta(ctx, metahash)
	if err != nil {
		panic("could not get metadata for the address")
	}
	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, types.QueryResMeta{Meta: meta})
	return
}

func queryAccount(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	if len(path) != 1 {
		return []byte{}, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Expecting 1 args. Found %d.", len(path))
	}

	addr, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		panic("Could not parse address " + path[0])
	}

	account := keeper.GetAccount(ctx, addr)
	if account == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, path[0])
	}

	res, err = codec.MarshalJSONIndent(legacyQuerierCdc, account)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
