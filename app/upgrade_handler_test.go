package app

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
	
	"io/ioutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/common"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/test-go/testify/require"
	sktypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	fgtypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/codec"
)

func TestMigrateStore(t *testing.T) {
	genesisState := loadState(t)

	app := Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(common.Bech32PrefixValAddr, common.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(common.Bech32PrefixConsAddr, common.Bech32PrefixConsPub)
	cfg.Seal()

	for _, m := range []string{"auth", "bank", "staking", "feegrant", "gov"} {
		app.mm.Modules[m].InitGenesis(ctx, app.appCodec, genesisState[m])
	}
	checkStaking(t, ctx, app, true)
	checkFeegrant(t, ctx, app, true)
	checkGov(t, ctx, app, true)
	transAddrPrefix(ctx, *app)
	checkStaking(t, ctx, app, false)
	checkFeegrant(t, ctx, app, false)
	checkGov(t, ctx, app, false)
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
	t *testing.T
	app *ShentuApp
	store sdk.KVStore
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
		jsonstr := MustMarshalJSON(v)
		c.checkStr(jsonstr)
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

func checkFeegrant(t *testing.T, ctx sdk.Context, app *ShentuApp, old bool) {
	fgKeeper := app.FeegrantKeeper
	store := ctx.KVStore(app.keys[fgtypes.StoreKey])
	ck := NewChecker(t, app, store, old)
	fgKeeper.IterateAllFeeAllowances(ctx, func(grant fgtypes.Grant) bool {
		ck.checkStr(grant.Grantee + grant.Granter)
		return false
	})
}
