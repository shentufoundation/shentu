package gov_test

//
//import (
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	"cosmossdk.io/log"
//	dbm "github.com/cometbft/cometbft-db"
//	abcitypes "github.com/cometbft/cometbft/abci/types"
//	tmjson "github.com/cometbft/cometbft/libs/json"
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//
//	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
//	"github.com/cosmos/cosmos-sdk/x/gov/types"
//
//	shentuapp "github.com/shentufoundation/shentu/v2/app"
//)
//
//func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
//	db := dbm.NewMemDB()
//	encCdc := shentuapp.MakeEncodingConfig()
//
//	app := shentuapp.NewShentuApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, shentuapp.DefaultNodeHome, 5, encCdc, shentuapp.EmptyAppOptions{})
//
//	genesisState := shentuapp.GenesisStateWithSingleValidator(t, app)
//	stateBytes, err := tmjson.Marshal(genesisState)
//	require.NoError(t, err)
//
//	app.InitChain(
//		abcitypes.RequestInitChain{
//			AppStateBytes: stateBytes,
//			ChainId:       "test-chain-id",
//		},
//	)
//
//	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
//	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
//	require.NotNil(t, acc)
//}
