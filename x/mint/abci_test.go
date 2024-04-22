package mint_test

import (
	"github.com/shentufoundation/shentu/v2/x/mint/types"
	"testing"
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/mint"
)

func TestBeginBlocker(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	k := app.MintKeeper

	p := types.DefaultGenesisState().GetParams()
	k.SetParams(ctx, p)
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
					Inflation:        sdk.NewDecWithPrec(12, 2),
					AnnualProvisions: sdk.NewDecWithPrec(7, 2)},
			},
		},
		{
			"zero inflation", args{
				minttypes.Minter{
					Inflation:        sdk.NewDecWithPrec(0, 2),
					AnnualProvisions: sdk.NewDecWithPrec(0, 2)},
			},
		},
		{
			"hundred inflation", args{
				minttypes.Minter{
					Inflation:        sdk.NewDecWithPrec(100, 2),
					AnnualProvisions: sdk.NewDecWithPrec(100, 2)},
			},
		},
	}
	for _, tt := range tests {
		k.SetMinter(ctx, tt.args.minter)
		t.Run(tt.name, func(t *testing.T) {
			mint.BeginBlocker(ctx, k)
		})
	}
}
