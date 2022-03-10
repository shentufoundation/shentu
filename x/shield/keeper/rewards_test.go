package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/certikfoundation/shentu/v2/app"
	shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"
)

// TestBlockRewardRatio tests the calculation of shield block reward ratio
func TestGetShieldBlockRewardRatio(t *testing.T) {
	app := shentuapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})

	tests := []struct {
		name            string
		totalShield     sdk.Int
		totalCollateral sdk.Int
		paramA          sdk.Dec
		paramB          sdk.Dec
		paramL          sdk.Dec
		expRatio        sdk.Dec
	}{
		{
			name:            "Get Shield Block Reward Ratio with All Zeros",
			totalShield:     sdk.ZeroInt(),
			totalCollateral: sdk.ZeroInt(),
			paramA:          sdk.ZeroDec(),
			paramB:          sdk.ZeroDec(),
			paramL:          sdk.ZeroDec(),
			expRatio:        sdk.ZeroDec(),
		},
		{
			name:            "Get Shield Block Reward Ratio with Zero Total Shield and Zero Total Collateral",
			totalShield:     sdk.ZeroInt(),
			totalCollateral: sdk.ZeroInt(),
			paramA:          sdk.NewDecWithPrec(30, 2), // 0.3
			paramB:          sdk.NewDecWithPrec(50, 2), // 0.5
			paramL:          sdk.NewDec(5),
			expRatio:        sdk.NewDecWithPrec(30, 2), // 0.3
		},
		{
			name:            "Get Shield Block Reward Ratio with Zero Total Collateral",
			totalShield:     sdk.NewInt(1e4),
			totalCollateral: sdk.ZeroInt(),
			paramA:          sdk.NewDecWithPrec(30, 2), // 0.3
			paramB:          sdk.NewDecWithPrec(50, 2), // 0.5
			paramL:          sdk.NewDec(5),
			expRatio:        sdk.NewDecWithPrec(30, 2), // 0.3
		},
		{
			name:            "Get Shield Block Reward Ratio with Zero Total Shield",
			totalShield:     sdk.ZeroInt(),
			totalCollateral: sdk.NewInt(1e4),
			paramA:          sdk.NewDecWithPrec(30, 2), // 0.3
			paramB:          sdk.NewDecWithPrec(50, 2), // 0.5
			paramL:          sdk.NewDec(5),
			expRatio:        sdk.NewDecWithPrec(30, 2), // 0.3
		},
		{
			name:            "Get Shield Block Reward Ratio with Total Shield Less Than Total Collateral",
			totalShield:     sdk.NewInt(1e3),
			totalCollateral: sdk.NewInt(1e4),
			paramA:          sdk.NewDecWithPrec(30, 2), // 0.3
			paramB:          sdk.NewDecWithPrec(50, 2), // 0.5
			paramL:          sdk.NewDec(5),
			expRatio:        sdk.NewDecWithPrec(307843137254901961, 18), // 0.3078
		},
		{
			name:            "Get Shield Block Reward Ratio with Total Shield Equal to Total Collateral",
			totalShield:     sdk.NewInt(1e4),
			totalCollateral: sdk.NewInt(1e4),
			paramA:          sdk.NewDecWithPrec(30, 2), // 0.3
			paramB:          sdk.NewDecWithPrec(50, 2), // 0.5
			paramL:          sdk.NewDec(5),
			expRatio:        sdk.NewDecWithPrec(366666666666666667, 18), // 0.3667
		},
		{
			name:            "Get Shield Block Reward Ratio with Total Shield Greater Than Total Collateral",
			totalShield:     sdk.NewInt(1e5),
			totalCollateral: sdk.NewInt(1e4),
			paramA:          sdk.NewDecWithPrec(30, 2), // 0.3
			paramB:          sdk.NewDecWithPrec(50, 2), // 0.5
			paramL:          sdk.NewDec(5),
			expRatio:        sdk.NewDecWithPrec(566666666666666667, 18), // 0.5667
		},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		app.ShieldKeeper.SetTotalShield(ctx, tc.totalShield)
		app.ShieldKeeper.SetTotalCollateral(ctx, tc.totalCollateral)
		app.ShieldKeeper.SetBlockRewardParams(ctx, shieldtypes.NewBlockRewardParams(tc.paramA, tc.paramB, tc.paramL))
		ratio := app.ShieldKeeper.GetShieldBlockRewardRatio(ctx)
		require.True(t, ratio.Equal(tc.expRatio))
	}
}
