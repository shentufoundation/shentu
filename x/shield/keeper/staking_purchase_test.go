package keeper_test

import (
	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"testing"
)

func TestKeeper_AddStaking(t *testing.T) {
	type args struct {
		ctx        sdk.Context
		poolID     uint64
		purchaser  sdk.AccAddress
		purchaseID uint64
		stakingAmt sdk.Int
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if err := k.AddStaking(tt.args.ctx, tt.args.poolID, tt.args.purchaser, tt.args.purchaseID, tt.args.stakingAmt); (err != nil) != tt.wantErr {
				t.Errorf("AddStaking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeeper_FundShieldBlockRewards(t *testing.T) {
	type args struct {
		ctx    sdk.Context
		amount sdk.Coins
		sender sdk.AccAddress
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if err := k.FundShieldBlockRewards(tt.args.ctx, tt.args.amount, tt.args.sender); (err != nil) != tt.wantErr {
				t.Errorf("FundShieldBlockRewards() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeeper_GetAllOriginalStakings(t *testing.T) {
	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name                 string
		keeper               keeper.Keeper
		args                 args
		wantOriginalStakings []types.OriginalStaking
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if gotOriginalStakings := k.GetAllOriginalStakings(tt.args.ctx); !reflect.DeepEqual(gotOriginalStakings, tt.wantOriginalStakings) {
				t.Errorf("GetAllOriginalStakings() = %v, want %v", gotOriginalStakings, tt.wantOriginalStakings)
			}
		})
	}
}

func TestKeeper_GetAllStakeForShields(t *testing.T) {
	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name          string
		keeper        keeper.Keeper
		args          args
		wantPurchases []types.ShieldStaking
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if gotPurchases := k.GetAllStakeForShields(tt.args.ctx); !reflect.DeepEqual(gotPurchases, tt.wantPurchases) {
				t.Errorf("GetAllStakeForShields() = %v, want %v", gotPurchases, tt.wantPurchases)
			}
		})
	}
}

func TestKeeper_GetGlobalShieldStakingPool(t *testing.T) {
	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name     string
		keeper   keeper.Keeper
		args     args
		wantPool sdk.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if gotPool := k.GetGlobalShieldStakingPool(tt.args.ctx); !reflect.DeepEqual(gotPool, tt.wantPool) {
				t.Errorf("GetGlobalShieldStakingPool() = %v, want %v", gotPool, tt.wantPool)
			}
		})
	}
}

func TestKeeper_GetOriginalStaking(t *testing.T) {
	type args struct {
		ctx        sdk.Context
		purchaseID uint64
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
		want   sdk.Int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if got := k.GetOriginalStaking(tt.args.ctx, tt.args.purchaseID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOriginalStaking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetStakeForShield(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		poolID    uint64
		purchaser sdk.AccAddress
	}
	tests := []struct {
		name         string
		keeper       keeper.Keeper
		args         args
		wantPurchase types.ShieldStaking
		wantFound    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			gotPurchase, gotFound := k.GetStakeForShield(tt.args.ctx, tt.args.poolID, tt.args.purchaser)
			if !reflect.DeepEqual(gotPurchase, tt.wantPurchase) {
				t.Errorf("GetStakeForShield() gotPurchase = %v, want %v", gotPurchase, tt.wantPurchase)
			}
			if gotFound != tt.wantFound {
				t.Errorf("GetStakeForShield() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestKeeper_IterateOriginalStakings(t *testing.T) {
	type args struct {
		ctx      sdk.Context
		callback func(original types.OriginalStaking) (stop bool)
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.keeper
		})
	}
}

func TestKeeper_IterateStakeForShields(t *testing.T) {
	type args struct {
		ctx      sdk.Context
		callback func(purchase types.ShieldStaking) (stop bool)
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.keeper
		})
	}
}

func TestKeeper_ProcessStakeForShieldExpiration(t *testing.T) {
	type args struct {
		ctx        sdk.Context
		poolID     uint64
		purchaseID uint64
		bondDenom  string
		purchaser  sdk.AccAddress
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if err := k.ProcessStakeForShieldExpiration(tt.args.ctx, tt.args.poolID, tt.args.purchaseID, tt.args.bondDenom, tt.args.purchaser); (err != nil) != tt.wantErr {
				t.Errorf("ProcessStakeForShieldExpiration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeeper_SetGlobalShieldStakingPool(t *testing.T) {
	type args struct {
		ctx   sdk.Context
		value sdk.Int
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.keeper
		})
	}
}

func TestKeeper_SetOriginalStaking(t *testing.T) {
	type args struct {
		ctx        sdk.Context
		purchaseID uint64
		amount     sdk.Int
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.keeper
		})
	}
}

func TestKeeper_SetStakeForShield(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		poolID    uint64
		purchaser sdk.AccAddress
		purchase  types.ShieldStaking
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.keeper
		})
	}
}

func TestKeeper_UnstakeFromShield(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		poolID    uint64
		purchaser sdk.AccAddress
		amount    sdk.Int
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if err := k.UnstakeFromShield(tt.args.ctx, tt.args.poolID, tt.args.purchaser, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("UnstakeFromShield() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
