package mint_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/mint"
	"github.com/shentufoundation/shentu/v2/x/mint/types"
)

func TestBeginBlocker(t *testing.T) {
	app := shentuapp.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)
	k := app.MintKeeper

	p := types.DefaultGenesisState().GetParams()
	err := k.Params.Set(ctx, p)
	require.NoError(t, err)
	type args struct {
		minter minttypes.Minter
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"normal", args{
				minttypes.Minter{
					Inflation:        math.LegacyNewDecWithPrec(12, 2),
					AnnualProvisions: math.LegacyNewDecWithPrec(7, 2)},
			},
		},
		{
			"zero inflation", args{
				minttypes.Minter{
					Inflation:        math.LegacyNewDecWithPrec(0, 2),
					AnnualProvisions: math.LegacyNewDecWithPrec(0, 2)},
			},
		},
		{
			"hundred inflation", args{
				minttypes.Minter{
					Inflation:        math.LegacyNewDecWithPrec(100, 2),
					AnnualProvisions: math.LegacyNewDecWithPrec(100, 2)},
			},
		},
	}
	for _, tt := range tests {
		err := k.Minter.Set(ctx, tt.args.minter)
		require.NoError(t, err)
		t.Run(tt.name, func(t *testing.T) {
			err = mint.BeginBlocker(ctx, k, minttypes.DefaultInflationCalculationFn)
			require.NoError(t, err)
		})
	}
}
