package v260_test

import (
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/common"
	v260 "github.com/shentufoundation/shentu/v2/x/gov/legacy/v260"
	"github.com/shentufoundation/shentu/v2/x/gov/types"
)

func Test_MigrateProposalStore(t *testing.T) {
	var (
		depositParams govtypes.DepositParams
		tallyParams   govtypes.TallyParams
		customParams  types.CustomParams
	)

	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	govSubspace := app.GetSubspace(govtypes.ModuleName)
	tableField := reflect.ValueOf(&govSubspace).Elem().FieldByName("table")
	tableFieldPtr := reflect.NewAt(tableField.Type(), unsafe.Pointer(tableField.UnsafeAddr()))
	tableFieldPtr.Elem().Set(reflect.ValueOf(v260.ParamKeyTable()))

	minInitialDepositTokens := sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	minDepositTokens := sdk.TokensFromConsensusPower(5, sdk.DefaultPowerReduction)
	defaultTally := govtypes.NewTallyParams(sdk.NewDecWithPrec(335, 3), sdk.NewDecWithPrec(6, 1), sdk.NewDecWithPrec(335, 3))
	certifierUpdateSecurityVoteTally := govtypes.NewTallyParams(sdk.NewDecWithPrec(335, 3), sdk.NewDecWithPrec(668, 3), sdk.NewDecWithPrec(335, 3))
	certifierUpdateStakeVoteTally := govtypes.NewTallyParams(sdk.NewDecWithPrec(335, 3), sdk.NewDecWithPrec(8, 1), sdk.NewDecWithPrec(335, 3))

	oldDepositParams := v260.DepositParams{
		MinInitialDeposit: sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, minInitialDepositTokens)},
		MinDeposit:        sdk.Coins{sdk.NewCoin(common.MicroCTKDenom, minDepositTokens)},
		MaxDepositPeriod:  govtypes.DefaultPeriod,
	}
	oldTallyParams := v260.TallyParams{
		DefaultTally:                     &defaultTally,
		CertifierUpdateSecurityVoteTally: &certifierUpdateSecurityVoteTally,
		CertifierUpdateStakeVoteTally:    &certifierUpdateStakeVoteTally,
	}
	// set old data
	govSubspace.Set(ctx, govtypes.ParamStoreKeyDepositParams, &oldDepositParams)
	govSubspace.Set(ctx, govtypes.ParamStoreKeyTallyParams, &oldTallyParams)

	tableFieldPtr.Elem().Set(reflect.ValueOf(types.ParamKeyTable()))
	err := v260.MigrateParams(ctx, govSubspace)
	require.NoError(t, err)
	// get migrate params
	govSubspace.Get(ctx, govtypes.ParamStoreKeyDepositParams, &depositParams)
	govSubspace.Get(ctx, govtypes.ParamStoreKeyTallyParams, &tallyParams)
	govSubspace.Get(ctx, types.ParamStoreKeyCustomParams, &customParams)

	require.Equal(t, depositParams.MinDeposit, oldDepositParams.MinDeposit)
	require.Equal(t, depositParams.MaxDepositPeriod, oldDepositParams.MaxDepositPeriod)
	require.Equal(t, tallyParams, defaultTally)
	require.Equal(t, customParams.CertifierUpdateSecurityVoteTally, &certifierUpdateSecurityVoteTally)
	require.Equal(t, customParams.CertifierUpdateStakeVoteTally, &certifierUpdateStakeVoteTally)
}
