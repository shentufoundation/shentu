package keeper_test

import (
	"fmt"
	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"testing"
)

func TestKeeper_AddStaking(t *testing.T) {
	type args struct {
		pools      []types.Pool
		poolID     uint64
		purchaser  sdk.AccAddress
		purchaseID uint64
		stakingAmt sdk.Int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Pool doesn't exist",
			args: args{
				poolID:     1,
				purchaser:  acc1,
				purchaseID: 1,
				stakingAmt: sdk.NewInt(1),
			},
			wantErr: true,
		},
		{
			name: "Pool exists",
			args: args{
				pools:      []types.Pool{DummyPool(1)},
				poolID:     1,
				purchaser:  acc1,
				purchaseID: 1,
				stakingAmt: sdk.NewInt(1),
			},
			wantErr: false,
		},
		{
			name: "Zero amount",
			args: args{
				pools:      []types.Pool{DummyPool(1)},
				poolID:     1,
				purchaser:  acc1,
				purchaseID: 1,
				stakingAmt: sdk.NewInt(0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			for _, p := range tt.args.pools {
				k.SetPool(suite.ctx, p)
			}
			if err := k.AddStaking(suite.ctx, tt.args.poolID, tt.args.purchaser, tt.args.purchaseID, tt.args.stakingAmt); (err != nil) != tt.wantErr {
				t.Errorf("AddStaking() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeeper_FundShieldBlockRewards(t *testing.T) {
	type args struct {
		amount sdk.Coins
		sender sdk.AccAddress
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Normal path",
			args: args{
				amount: sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.OneInt())),
				sender: acc1,
			},
			wantErr: false,
		},
		{
			name: "Zero Amount",
			args: args{
				amount: sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.ZeroInt())),
				sender: acc1,
			},
			wantErr: false,
		},
		{
			name: "Nil Account",
			args: args{
				amount: sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.OneInt())),
				sender: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			if err := k.FundShieldBlockRewards(suite.ctx, tt.args.amount, tt.args.sender); (err != nil) != tt.wantErr {
				t.Errorf("FundShieldBlockRewards() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeeper_GetAllOriginalStakings(t *testing.T) {
	type args struct {
		stakings []types.ShieldStaking
	}
	tests := []struct {
		name                 string
		args                 args
		wantOriginalStakings []types.OriginalStaking
	}{
		{
			name:                 "Empty stakings",
			args:                 args{},
			wantOriginalStakings: nil,
		},
		{
			name: "One staking",
			args: args{
				stakings: []types.ShieldStaking{
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.OneInt(),
						WithdrawRequested: sdk.ZeroInt(),
					},
				},
			},
			wantOriginalStakings: []types.OriginalStaking{
				{
					PurchaseId: 0,
					Amount:     sdk.OneInt(),
				},
			},
		},
		{
			name: "Two staking",
			args: args{
				stakings: []types.ShieldStaking{
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.OneInt(),
						WithdrawRequested: sdk.ZeroInt(),
					},
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.NewInt(2),
						WithdrawRequested: sdk.ZeroInt(),
					},
				},
			},
			wantOriginalStakings: []types.OriginalStaking{
				{
					PurchaseId: 0,
					Amount:     sdk.OneInt(),
				},
				{
					PurchaseId: 1,
					Amount:     sdk.NewInt(2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			k.SetPool(suite.ctx, DummyPool(1))
			for i, s := range tt.args.stakings {
				purchaser, err := sdk.AccAddressFromBech32(s.Purchaser)
				if err != nil {
					t.Errorf("got err = %v, want no error while setting up the test", err)
				}
				err = k.AddStaking(suite.ctx, s.PoolId, purchaser, uint64(i), s.Amount)
				if err != nil {
					t.Errorf("got err = %v, want no error while setting up the test", err)
				}
			}
			if gotOriginalStakings := k.GetAllOriginalStakings(suite.ctx); !reflect.DeepEqual(gotOriginalStakings, tt.wantOriginalStakings) {
				t.Errorf("GetAllOriginalStakings() = %v, want %v", gotOriginalStakings, tt.wantOriginalStakings)
			}
		})
	}
}

func TestKeeper_GetAllStakeForShields(t *testing.T) {
	type args struct {
		sfs []types.ShieldStaking
	}
	tests := []struct {
		name          string
		args          args
		wantPurchases []types.ShieldStaking
	}{
		{
			name:          "Empty list",
			args:          args{},
			wantPurchases: []types.ShieldStaking{},
		},
		{
			name: "One stake for shield",
			args: args{
				sfs: []types.ShieldStaking{
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.NewInt(2),
						WithdrawRequested: sdk.ZeroInt(),
					},
				},
			},
			wantPurchases: []types.ShieldStaking{
				{
					PoolId:            1,
					Purchaser:         acc1.String(),
					Amount:            sdk.NewInt(2),
					WithdrawRequested: sdk.ZeroInt(),
				},
			},
		},
		{
			name: "Two stake for shield from one purchaser to one pool",
			args: args{
				sfs: []types.ShieldStaking{
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.OneInt(),
						WithdrawRequested: sdk.ZeroInt(),
					},
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.NewInt(2),
						WithdrawRequested: sdk.ZeroInt(),
					},
				},
			},
			wantPurchases: []types.ShieldStaking{
				{
					PoolId:            1,
					Purchaser:         acc1.String(),
					Amount:            sdk.NewInt(3),
					WithdrawRequested: sdk.ZeroInt(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			k.SetPool(suite.ctx, DummyPool(1))
			for i, sfs := range tt.args.sfs {
				purchaser, err := sdk.AccAddressFromBech32(sfs.Purchaser)
				if err != nil {
					panic(err)
				}
				err = k.AddStaking(suite.ctx, 1, purchaser, uint64(i), sfs.Amount)
				if err != nil {
					panic(err)
				}
			}
			if gotPurchases := k.GetAllStakeForShields(suite.ctx); !reflect.DeepEqual(gotPurchases, tt.wantPurchases) {
				for i, sfs := range gotPurchases {
					if !reflect.DeepEqual(sfs, tt.wantPurchases[i]) {
						t.Errorf("GetAllStakeForShields() = %v, want %v", gotPurchases, tt.wantPurchases)
					}
				}
			}
		})
	}
}

func TestKeeper_GetGlobalShieldStakingPool(t *testing.T) {
	type args struct {
		sfs []types.ShieldStaking
	}
	tests := []struct {
		name     string
		args     args
		wantPool sdk.Int
	}{
		{
			name:     "Empty pool",
			args:     args{},
			wantPool: sdk.ZeroInt(),
		},
		{
			name: "One pool",
			args: args{
				sfs: []types.ShieldStaking{
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.OneInt(),
						WithdrawRequested: sdk.ZeroInt(),
					},
				},
			},
			wantPool: sdk.OneInt(),
		},
		{
			name: "Three pool",
			args: args{
				sfs: []types.ShieldStaking{
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.OneInt(),
						WithdrawRequested: sdk.ZeroInt(),
					},
					{
						PoolId:            1,
						Purchaser:         acc1.String(),
						Amount:            sdk.NewInt(2),
						WithdrawRequested: sdk.ZeroInt(),
					},
				},
			},
			wantPool: sdk.NewInt(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			k.SetPool(suite.ctx, DummyPool(1))
			for i, sfs := range tt.args.sfs {
				purchaser, err := sdk.AccAddressFromBech32(sfs.Purchaser)
				if err != nil {
					panic(err)
				}
				err = k.AddStaking(suite.ctx, 1, purchaser, uint64(i), sfs.Amount)
				if err != nil {
					panic(err)
				}
			}
			if gotPool := k.GetGlobalShieldStakingPool(suite.ctx); !reflect.DeepEqual(gotPool, tt.wantPool) {
				t.Errorf("GetGlobalShieldStakingPool() = %v, want %v", gotPool, tt.wantPool)
			}
		})
	}
}

func TestKeeper_GetOriginalStaking(t *testing.T) {
	modpurchase := basePurchase
	modpurchase.Shield = sdk.NewInt(50000000)
	type args struct {
		purchases  []types.Purchase
		purchaseID uint64
	}
	tests := []struct {
		name string
		args args
		want sdk.Int
	}{
		{
			name: "No original staking",
			args: args{},
			want: sdk.ZeroInt(),
		},
		{
			name: "Valid original staking",
			args: args{
				purchases: []types.Purchase{
					modpurchase,
				},
				purchaseID: 1,
			},
			want: sdk.NewInt(50000000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			for _, p := range tt.args.purchases {
				simapp.AddTestAddrsFromPubKeys(suite.app, suite.ctx, PKS, sdk.NewInt(2e8))
				suite.tstaking.CreateValidatorWithValPower(sdk.ValAddress(PKS[0].Address()), PKS[0], 10000, true)
				suite.tshield.DepositCollateral(acc1, 500000000, true)
				k.SetPool(suite.ctx, DummyPool(1))
				_, err := k.PurchaseShield(suite.ctx, 1, sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, p.Shield)), "", acc1, false)
				if err != nil {
					panic(err)
				}
			}
			fmt.Println(k.GetOriginalStaking(suite.ctx, 1))
			if got := k.GetOriginalStaking(suite.ctx, tt.args.purchaseID); !reflect.DeepEqual(got, tt.want) {
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
