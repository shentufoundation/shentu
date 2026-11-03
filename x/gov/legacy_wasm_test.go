package gov_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	legacywasm "github.com/shentufoundation/shentu/v2/app/legacy/wasm"
	govmodule "github.com/shentufoundation/shentu/v2/x/gov"
	govtypes "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

func TestGenesisRoundTripHistoricalWasmParams(t *testing.T) {
	// Independent wire fixture using the wasmd v0.55.1 schema, rather than
	// packing a local message (which would prepopulate Any's decoded cache).
	// authority="gov", upload=AnyOfAddresses["alice","bob"], instantiate=Everybody.
	wire, err := hex.DecodeString("0a03676f7612120a0e08041a05616c6963651a03626f621003")
	require.NoError(t, err)
	proposal := govv1.Proposal{
		Id: 42, Status: govv1.StatusPassed,
		Messages: []*codectypes.Any{{TypeUrl: "/cosmwasm.wasm.v1.MsgUpdateParams", Value: wire}},
	}
	stored, err := proposal.Marshal()
	require.NoError(t, err)

	// Reproduce the original failure without historical type registration.
	registry := codectypes.NewInterfaceRegistry()
	govv1.RegisterInterfaces(registry)
	var decoded govv1.Proposal
	err = codec.NewProtoCodec(registry).Unmarshal(stored, &decoded)
	require.ErrorContains(t, err, "no concrete type registered for type URL /cosmwasm.wasm.v1.MsgUpdateParams")

	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	require.NoError(t, app.GovKeeper.ProposalID.Set(ctx, 43))
	// Persist the proposal without decoding its Any; export must decode from storage.
	require.NoError(t, app.GovKeeper.Proposals.Set(ctx, 42, proposal))
	module := govmodule.NewAppModule(app.AppCodec(), app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace("gov"))
	jsonState := module.ExportGenesis(ctx, app.AppCodec())
	require.Contains(t, string(jsonState), `"@type":"/cosmwasm.wasm.v1.MsgUpdateParams"`)

	app2 := shentuapp.Setup(t, false)
	ctx2 := app2.BaseApp.NewContext(false)
	module2 := govmodule.NewAppModule(app2.AppCodec(), app2.GovKeeper, app2.AccountKeeper, app2.BankKeeper, app2.GetSubspace("gov"))
	module2.InitGenesis(ctx2, app2.AppCodec(), jsonState)
	jsonState2 := module2.ExportGenesis(ctx2, app2.AppCodec())
	require.JSONEq(t, string(jsonState), string(jsonState2))

	var genesis govtypes.GenesisState
	require.NoError(t, app2.AppCodec().UnmarshalJSON(jsonState2, &genesis))
	require.Len(t, genesis.Proposals, 1)
	require.Equal(t, wire, genesis.Proposals[0].Messages[0].Value)
	msgs, err := genesis.Proposals[0].GetMsgs()
	require.NoError(t, err)
	msg := msgs[0].(*legacywasm.MsgUpdateParams)
	require.Equal(t, "gov", msg.Authority)
	require.Equal(t, legacywasm.AccessTypeAnyOfAddresses, msg.Params.CodeUploadAccess.Permission)
	require.Equal(t, []string{"alice", "bob"}, msg.Params.CodeUploadAccess.Addresses)
	require.Equal(t, legacywasm.AccessTypeEverybody, msg.Params.InstantiateDefaultPermission)
	require.Nil(t, app2.MsgServiceRouter().Handler(msg), "historical decoding must not enable wasm execution")
}
