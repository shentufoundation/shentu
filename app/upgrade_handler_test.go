package app

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/test-go/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkauthtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	fgtypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	sktypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/shentufoundation/shentu/v2/common"
	authtypes "github.com/shentufoundation/shentu/v2/x/auth/types"
)

func setConfig(prefix string) {
	stoc := func(txt string) string {
		if prefix == "certik" {
			return strings.Replace(txt, "shentu", "certik", 1)
		}
		return txt
	}
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(stoc(common.Bech32PrefixAccAddr), stoc(common.Bech32PrefixAccPub))
	cfg.SetBech32PrefixForValidator(stoc(common.Bech32PrefixValAddr), stoc(common.Bech32PrefixValPub))
	cfg.SetBech32PrefixForConsensusNode(stoc(common.Bech32PrefixConsAddr), stoc(common.Bech32PrefixConsPub))
}

func TestMigrateStore(t *testing.T) {
	genesisState := loadState(t)

	app := Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	setConfig("certik")

	for _, m := range []string{"auth", "authz", "bank", "staking", "gov", "slashing"} {
		app.mm.Modules[m].InitGenesis(ctx, app.appCodec, genesisState[m])
	}
	//it has to be independently set store for feegrant to avoid affect accAddrCache
	//setStoreForFeegrant(ctx, app, genesisState["feegrant"])

	checkStaking(t, ctx, app, true)
	checkFeegrant(t, ctx, app, true)
	checkGov(t, ctx, app, true)
	checkSlashing(t, ctx, app, true)
	//checkAuth(t, ctx, app, true)
	//checkAuthz(t, ctx, app, true)

	setConfig("shentu")
	transAddrPrefix(ctx, *app)

	checkStaking(t, ctx, app, false)
	checkFeegrant(t, ctx, app, false)
	checkGov(t, ctx, app, false)
	checkSlashing(t, ctx, app, false)
	//checkAuth(t, ctx, app, false)
	//checkAuthz(t, ctx, app, false)
}

func loadState(t *testing.T) GenesisState {
	data, err := ioutil.ReadFile("../tests/pruned-state.json")
	if err != nil {
		t.Fatal("failed to read in json")
	}
	var genesisState GenesisState
	if err = json.Unmarshal(data, &genesisState); err != nil {
		t.Fatal("failed to parse the json")
	}
	return genesisState
}

func MustMarshalJSON(v any) string {
	bs, err := json.Marshal(v)
	if err != nil {
		panic("failed to do json marshal")
	}
	return string(bs)
}

func NewChecker(t *testing.T, app *ShentuApp, store sdk.KVStore, old bool) Checker {
	prefixPos, prefixNeg := "shentu", "certik"
	if old {
		prefixPos, prefixNeg = "certik", "shentu"
	}
	return Checker{t, app, store, prefixPos, prefixNeg}
}

type Checker struct {
	t         *testing.T
	app       *ShentuApp
	store     sdk.KVStore
	prefixPos string
	prefixNeg string
}

func (c Checker) checkForOneKey(keyPrefix []byte, v any) {
	iter := sdk.KVStorePrefixIterator(c.store, keyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		iv, ok := v.(codec.ProtoMarshaler)
		if !ok {
			panic("failed to cast to codec.ProtoMarshaler")
		}
		c.app.appCodec.MustUnmarshal(iter.Value(), iv)
		jsonStr := MustMarshalJSON(v)
		c.checkStr(jsonStr)
	}
}

func (c Checker) checkForAuth(keyPrefix []byte) {
	ak := c.app.AccountKeeper
	iter := sdk.KVStorePrefixIterator(c.store, keyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		acc, err := ak.UnmarshalAccount(iter.Value())
		if err != nil {
			panic(err)
		}

		switch account := acc.(type) {
		case *sdkauthtypes.BaseAccount:
			c.checkStr(account.String())
		case *sdkauthtypes.ModuleAccount:
			c.checkStr(account.String())
		case *authtypes.ManualVestingAccount:
			c.checkStr(account.String())
		default:
			panic("unknown account type")
		}
	}
}

func (c Checker) checkStr(str string) {
	require.True(c.t, strings.Contains(str, c.prefixPos))
	require.False(c.t, strings.Contains(str, c.prefixNeg))
}

func checkStaking(t *testing.T, ctx sdk.Context, app *ShentuApp, old bool) {
	skKeeper := app.StakingKeeper.Keeper
	store := ctx.KVStore(app.keys[sktypes.StoreKey])
	ck := NewChecker(t, app, store, old)

	for _, v := range skKeeper.GetAllValidators(ctx) {
		ck.checkStr(v.OperatorAddress)
	}
	ck.checkForOneKey(sktypes.DelegationKey, &sktypes.Delegation{})
	ck.checkForOneKey(sktypes.UnbondingDelegationKey, &sktypes.UnbondingDelegation{})
	ck.checkForOneKey(sktypes.RedelegationKey, &sktypes.Redelegation{})
	ck.checkForOneKey(sktypes.UnbondingQueueKey, &sktypes.DVPairs{})
	ck.checkForOneKey(sktypes.RedelegationQueueKey, &sktypes.DVVTriplets{})
	ck.checkForOneKey(sktypes.ValidatorQueueKey, &sktypes.ValAddresses{})
	ck.checkForOneKey(sktypes.HistoricalInfoKey, &sktypes.HistoricalInfo{})

}

func checkGov(t *testing.T, ctx sdk.Context, app *ShentuApp, old bool) {
	store := ctx.KVStore(app.keys[govtypes.StoreKey])
	ck := NewChecker(t, app, store, old)
	ck.checkForOneKey(govtypes.DepositsKeyPrefix, &govtypes.Deposit{})
	ck.checkForOneKey(govtypes.VotesKeyPrefix, &govtypes.Vote{})
}

func setStoreForFeegrant(ctx sdk.Context, app *ShentuApp, jraw json.RawMessage) {
	store := ctx.KVStore(app.keys[fgtypes.StoreKey])
	var fggs fgtypes.GenesisState
	app.appCodec.MustUnmarshalJSON(jraw, &fggs)
	for _, one := range fggs.Allowances {
		granter := sdk.MustAccAddressFromBech32(one.Granter)
		grantee := sdk.MustAccAddressFromBech32(one.Grantee)
		key := fgtypes.FeeAllowanceKey(granter, grantee)
		bz := app.appCodec.MustMarshal(&one)
		store.Set(key, bz)
	}
}

func checkFeegrant(t *testing.T, ctx sdk.Context, app *ShentuApp, old bool) {
	fgKeeper := app.FeegrantKeeper
	store := ctx.KVStore(app.keys[fgtypes.StoreKey])
	ck := NewChecker(t, app, store, old)
	fgKeeper.IterateAllFeeAllowances(ctx, func(grant fgtypes.Grant) bool {
		ck.checkStr(grant.Grantee + grant.Granter)
		return false
	})
}

func checkAuth(t *testing.T, ctx sdk.Context, app *ShentuApp, old bool) {
	store := ctx.KVStore(app.keys[sdkauthtypes.StoreKey])
	ck := NewChecker(t, app, store, old)
	ck.checkForAuth(sdkauthtypes.AddressStoreKeyPrefix)
}

func checkSlashing(t *testing.T, ctx sdk.Context, app *ShentuApp, old bool) {
	store := ctx.KVStore(app.keys[slashingtypes.StoreKey])
	ck := NewChecker(t, app, store, old)
	ck.checkForOneKey(slashingtypes.ValidatorSigningInfoKeyPrefix, &slashingtypes.ValidatorSigningInfo{})
}

func checkAuthz(t *testing.T, ctx sdk.Context, app *ShentuApp, old bool) {
	store := ctx.KVStore(app.keys[authzkeeper.StoreKey])
	ck := NewChecker(t, app, store, old)
	ck.checkForOneKey(authzkeeper.GrantKey, &authz.Grant{})
}
