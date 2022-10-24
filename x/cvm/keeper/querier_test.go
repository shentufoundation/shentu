package keeper_test

import (
	"encoding/hex"
	"testing"
	"time"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/txs/payload"

	. "github.com/shentufoundation/shentu/v2/x/cvm/keeper"
	"github.com/shentufoundation/shentu/v2/x/cvm/types"
)

func TestNewQuerier(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cvmk := app.CVMKeeper
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(cvmk, app.LegacyAmino())

	bz, err := querier(ctx, []string{"other"}, query)
	require.Error(t, err)
	require.Nil(t, bz)

	path := []string{"code", addrs[0].String()}

	bz, err = querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	path = []string{"abi", addrs[0].String()}

	bz, err = querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	path = []string{"account", addrs[0].String()}

	bz, err = querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)
}

func TestViewQuery(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	cvmk := app.CVMKeeper

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := NewQuerier(cvmk, app.LegacyAmino())

	code, err := hex.DecodeString(Hello55BytecodeString)
	require.Nil(t, err)

	newContractAddress, err2 := cvmk.Tx(ctx, addrs[0], nil, 0, code, []*payload.ContractMeta{}, false, false, false)
	require.Nil(t, err2)
	require.NotNil(t, newContractAddress)

	contAddr, err := sdk.AccAddressFromHex(hex.EncodeToString(newContractAddress))
	require.Nil(t, err)
	path := []string{"code", contAddr.String()}

	bz, err := querier(ctx, path, query)
	require.NoError(t, err)
	require.NotNil(t, bz)

	getMyFavoriteNumberCall, _, err := abi.EncodeFunctionCall(
		Hello55AbiJsonString,
		"sayHi",
		WrapLogger(ctx.Logger()),
	)
	require.Nil(t, err)

	path = []string{"view", addrs[0].String(), contAddr.String()}
	query.Data = getMyFavoriteNumberCall
	bz, err = querier(ctx, path, query)

	var res types.QueryResView
	err = app.LegacyAmino().UnmarshalJSON(bz, &res)
	require.Nil(t, err)
	out, err := abi.DecodeFunctionReturn(Hello55AbiJsonString, "sayHi", res.Ret)
	require.Equal(t, "55", out[0].Value)
}

func TestQueryMeta(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	addrs := shentuapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(80000*1e6))
	cvmk := app.CVMKeeper

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}
	querier := NewQuerier(cvmk, app.LegacyAmino())

	state := cvmk.NewState(ctx)

	callframe := engine.NewCallFrame(state, acmstate.Named("TxCache"))
	cache := callframe.Cache

	runtime, err := hex.DecodeString(Hello55BytecodeString)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(runtime)
	var codehash acmstate.CodeHash
	copy(codehash[:], hash.Sum(nil))

	metadata := Hello55MetadataJsonString
	payloadMeta := payload.ContractMeta{CodeHash: codehash.Bytes(), Meta: metadata}
	require.Nil(t, err)
	addr, err := crypto.AddressFromBytes(addrs[0].Bytes())
	require.Nil(t, err)
	err = engine.UpdateContractMeta(cache, state, addr, []*payload.ContractMeta{&payloadMeta})
	require.Nil(t, err)
	err = cache.Sync(state)
	require.Nil(t, err)

	metaAddr, err := sdk.AccAddressFromHex(hex.EncodeToString(addr.Bytes()))
	require.Nil(t, err)
	path := []string{"address-meta", metaAddr.String()}
	bz, err := querier(ctx, path, query)
	require.NotNil(t, bz)
	require.Nil(t, err)

	path = []string{"meta", "3D6C2B3049DCD9E34EDE8EB0DA5F2FB3E5667A12942AFF5E6F57435231526AAE"}
	bz, err = querier(ctx, path, query)
	require.NotNil(t, bz)
	require.Nil(t, err)

	path = []string{"meta", "3B6C2B3049DCD9E34EDE8EB0DA5F2FB3E5667A12942AFF5E6F57435231526AAE"}
	bz, err = querier(ctx, path, query)
	require.Nil(t, err)
	var meta types.QueryResMeta
	err = app.LegacyAmino().UnmarshalJSON(bz, &meta)
	require.Nil(t, err)
	require.Equal(t, meta.Meta, "")
}
