package keeper

import (
	"context"
	"encoding/hex"

	"github.com/hyperledger/burrow/execution/evm/abi"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/x/cvm/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

func (q Querier) Code(c context.Context, request *types.QueryCodeRequest) (*types.QueryCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}
	vmAddr := crypto.MustAddressFromBytes(addr)

	code, err := q.GetCode(ctx, vmAddr)
	codeString := hex.EncodeToString(code)
	if err != nil {
		return nil, err
	}
	return &types.QueryCodeResponse{
		Code: codeString,
	}, nil
}

func (q Querier) Abi(c context.Context, request *types.QueryAbiRequest) (*types.QueryAbiResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}
	vmAddr, _ := crypto.AddressFromBytes(addr)
	abiString := string(q.GetAbi(ctx, vmAddr))
	return &types.QueryAbiResponse{
		Abi: abiString,
	}, nil
}

func (q Querier) Storage(c context.Context, request *types.QueryStorageRequest) (*types.QueryStorageResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}
	vmAddr, _ := crypto.AddressFromBytes(addr)

	key, err := hex.DecodeString(request.Key)
	if err != nil {
		return nil, err
	}

	word256Key := binary.LeftPadWord256(key)

	storage, err := q.GetStorage(ctx, vmAddr, word256Key)
	if err != nil {
		return nil, err
	}

	return &types.QueryStorageResponse{
		Value: storage,
	}, nil
}

func (q Querier) AddressMeta(c context.Context, request *types.QueryAddressMetaRequest) (*types.QueryAddressMetaResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}
	vmAddr, _ := crypto.AddressFromBytes(addr)
	state := q.NewState(ctx)
	addrMeta, err := state.GetAddressMeta(vmAddr)
	if err != nil {
		return nil, err
	}

	var resString string
	for i := range addrMeta {
		addrMeta[i].MetadataHash = []byte(addrMeta[i].MetadataHash.String())
		addrMeta[i].CodeHash = []byte(addrMeta[i].CodeHash.String())
		resString += addrMeta[i].String()
	}

	return &types.QueryAddressMetaResponse{
		MetaHash: resString,
	}, nil
}

func (q Querier) Meta(c context.Context, request *types.QueryMetaRequest) (*types.QueryMetaResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	state := q.NewState(ctx)
	hash, err := hex.DecodeString(request.Hash)
	if err != nil {
		return nil, err
	}
	var metahash acmstate.MetadataHash
	copy(metahash[:], hash)

	meta, err := state.GetMetadata(metahash)
	if err != nil {
		return nil, err
	}
	return &types.QueryMetaResponse{
		Meta: meta,
	}, nil
}

func (q Querier) Account(c context.Context, request *types.QueryAccountRequest) (*acm.Account, error) {
	ctx := sdk.UnwrapSDKContext(c)
	state := q.NewState(ctx)
	addr, err := sdk.AccAddressFromBech32(request.Address)
	if err != nil {
		return nil, err
	}
	vmAddr, _ := crypto.AddressFromBytes(addr)
	account, err := state.GetAccount(vmAddr)

	return account, nil
}

func (q Querier) View(c context.Context, request *types.QueryViewRequest) (*types.QueryViewResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var caller sdk.AccAddress
	var err error
	if request.Caller == "" {
		caller = crypto.ZeroAddress.Bytes()
	} else {
		caller, err = sdk.AccAddressFromBech32(request.Caller)
	}
	if err != nil {
		return nil, err
	}

	callee, err := sdk.AccAddressFromBech32(request.Callee)
	if err != nil {
		return nil, err
	}
	ret, err := q.Tx(ctx, caller, callee, 0, request.Data, nil, true, false, false)

	out, err := abi.DecodeFunctionReturn(string(request.AbiSpec), request.FunctionName, ret)
	if err != nil {
		return nil, err
	}
	result := []*types.ReturnVars{}
	for _, v := range out {
		result = append(result, &types.ReturnVars{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return &types.QueryViewResponse{
		ReturnVars: result,
	}, nil
}

var _ types.QueryServer = Querier{}
