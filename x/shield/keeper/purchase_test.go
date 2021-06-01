package keeper_test

import (
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/types"
)

func TestKeeper_AddPurchase(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		poolID    uint64
		purchaser sdk.AccAddress
		purchase  types.Purchase
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
		want   types.PurchaseList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if got := k.AddPurchase(tt.args.ctx, tt.args.poolID, tt.args.purchaser, tt.args.purchase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddPurchase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_DeletePurchaseList(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		poolID    uint64
		purchaser sdk.AccAddress
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
			if err := k.DeletePurchaseList(tt.args.ctx, tt.args.poolID, tt.args.purchaser); (err != nil) != tt.wantErr {
				t.Errorf("DeletePurchaseList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeeper_DequeuePurchase(t *testing.T) {
	type args struct {
		ctx          sdk.Context
		purchaseList types.PurchaseList
		timestamp    time.Time
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

func TestKeeper_ExpiringPurchaseQueueIterator(t *testing.T) {
	type args struct {
		ctx     sdk.Context
		endTime time.Time
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
		want   sdk.Iterator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if got := k.ExpiringPurchaseQueueIterator(tt.args.ctx, tt.args.endTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExpiringPurchaseQueueIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetAllPurchaseLists(t *testing.T) {
	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name          string
		keeper        keeper.Keeper
		args          args
		wantPurchases []types.PurchaseList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if gotPurchases := k.GetAllPurchaseLists(tt.args.ctx); !reflect.DeepEqual(gotPurchases, tt.wantPurchases) {
				t.Errorf("GetAllPurchaseLists() = %v, want %v", gotPurchases, tt.wantPurchases)
			}
		})
	}
}

func TestKeeper_GetAllPurchases(t *testing.T) {
	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name          string
		keeper        keeper.Keeper
		args          args
		wantPurchases []types.Purchase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if gotPurchases := k.GetAllPurchases(tt.args.ctx); !reflect.DeepEqual(gotPurchases, tt.wantPurchases) {
				t.Errorf("GetAllPurchases() = %v, want %v", gotPurchases, tt.wantPurchases)
			}
		})
	}
}

func TestKeeper_GetExpiringPurchaseQueueTimeSlice(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		timestamp time.Time
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
		want   []types.PoolPurchaser
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if got := k.GetExpiringPurchaseQueueTimeSlice(tt.args.ctx, tt.args.timestamp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExpiringPurchaseQueueTimeSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetNextPurchaseID(t *testing.T) {
	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
		want   uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if got := k.GetNextPurchaseID(tt.args.ctx); got != tt.want {
				t.Errorf("GetNextPurchaseID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetPoolPurchaseLists(t *testing.T) {
	type args struct {
		ctx    sdk.Context
		poolID uint64
	}
	tests := []struct {
		name          string
		keeper        keeper.Keeper
		args          args
		wantPurchases []types.PurchaseList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if gotPurchases := k.GetPoolPurchaseLists(tt.args.ctx, tt.args.poolID); !reflect.DeepEqual(gotPurchases, tt.wantPurchases) {
				t.Errorf("GetPoolPurchaseLists() = %v, want %v", gotPurchases, tt.wantPurchases)
			}
		})
	}
}

func TestKeeper_GetPurchase(t *testing.T) {
	type args struct {
		purchaseList types.PurchaseList
		purchaseID   uint64
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
		want   types.Purchase
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ke := tt.keeper
			got, got1 := ke.GetPurchase(tt.args.purchaseList, tt.args.purchaseID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPurchase() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetPurchase() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKeeper_GetPurchaseList(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		poolID    uint64
		purchaser sdk.AccAddress
	}
	tests := []struct {
		name   string
		keeper keeper.Keeper
		args   args
		want   types.PurchaseList
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			got, got1 := k.GetPurchaseList(tt.args.ctx, tt.args.poolID, tt.args.purchaser)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPurchaseList() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetPurchaseList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKeeper_GetPurchaserPurchases(t *testing.T) {
	type args struct {
		pp      []poolpurchase
		address sdk.AccAddress
	}
	var tests = []struct {
		name    string
		args    args
		wantRes []types.PurchaseList
	}{
		{
			name: "Empty purchases",
			args: args{
				pp:      []poolpurchase{},
				address: acc1,
			},
			wantRes: []types.PurchaseList{},
		},
		{
			name: "one purchase",
			args: args{
				pp: []poolpurchase{
					{poolID: 1,
						purchases: []types.Purchase{
							{
								PurchaseId:        1,
								ProtectionEndTime: time.Time{},
								DeletionTime:      time.Time{},
								Description:       "",
								Shield:            sdk.NewInt(1),
								ServiceFees:       OneMixedDecCoins(common.MicroCTKDenom),
							},
						},
					},
				},
				address: acc1,
			},
			wantRes: []types.PurchaseList{
				{
					PoolId:    1,
					Purchaser: acc1.String(),
					Entries: []types.Purchase{
						{
							PurchaseId:        1,
							ProtectionEndTime: time.Time{},
							DeletionTime:      time.Time{},
							Description:       "",
							Shield:            sdk.NewInt(1),
							ServiceFees:       OneMixedDecCoins(common.MicroCTKDenom),
						},
					},
				},
			},
		},
		{
			name: "one purchase each in two pools",
			args: args{
				pp: []poolpurchase{
					{
						poolID: 1,
						purchases: []types.Purchase{
							{
								PurchaseId:        1,
								ProtectionEndTime: time.Time{},
								DeletionTime:      time.Time{},
								Description:       "",
								Shield:            sdk.NewInt(1),
								ServiceFees:       OneMixedDecCoins(common.MicroCTKDenom),
							},
						},
					},
					{
						poolID: 2,
						purchases: []types.Purchase{
							{
								PurchaseId:        2,
								ProtectionEndTime: time.Time{},
								DeletionTime:      time.Time{},
								Description:       "",
								Shield:            sdk.NewInt(1),
								ServiceFees:       OneMixedDecCoins(common.MicroCTKDenom),
							},
						},
					},
				},
				address: acc1,
			},
			wantRes: []types.PurchaseList{
				{
					PoolId:    1,
					Purchaser: acc1.String(),
					Entries: []types.Purchase{
						{
							PurchaseId:        1,
							ProtectionEndTime: time.Time{},
							DeletionTime:      time.Time{},
							Description:       "",
							Shield:            sdk.NewInt(1),
							ServiceFees:       OneMixedDecCoins(common.MicroCTKDenom),
						},
					},
				},
				{
					PoolId:    2,
					Purchaser: acc1.String(),
					Entries: []types.Purchase{
						{
							PurchaseId:        2,
							ProtectionEndTime: time.Time{},
							DeletionTime:      time.Time{},
							Description:       "",
							Shield:            sdk.NewInt(1),
							ServiceFees:       OneMixedDecCoins(common.MicroCTKDenom),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			k.SetPool(suite.ctx, DummyPool(1))
			k.SetPool(suite.ctx, DummyPool(2))
			for _, p := range tt.args.pp {
				for _, pur := range p.purchases {
					k.AddPurchase(suite.ctx, p.poolID, acc1, pur)
				}
			}
			if gotRes := k.GetPurchaserPurchases(suite.ctx, tt.args.address); !reflect.DeepEqual(gotRes, tt.wantRes) {
				if len(gotRes) != 0 || len(tt.wantRes) != 0 {
					t.Errorf("GetPurchaserPurchases() = %v, want %v", gotRes, tt.wantRes)
				}
			}
		})
	}
}

func TestKeeper_InsertExpiringPurchaseQueue(t *testing.T) {
	type pe struct {
		purchaseList types.PurchaseList
		endTime      time.Time
	}
	tests := []struct {
		name     string
		toInsert []pe
	}{
		{
			name: "One Insertion, no purchase entry",
			toInsert: []pe{
				{
					purchaseList: types.PurchaseList{
						PoolId:    1,
						Purchaser: acc1.String(),
						Entries:   nil,
					},
					endTime: time.Now().Add(12345),
				},
			},
		},
		{
			name: "One Insertion, one purchase entry",
			toInsert: []pe{
				{
					purchaseList: types.PurchaseList{
						PoolId:    1,
						Purchaser: acc1.String(),
						Entries: []types.Purchase{
							{
								PurchaseId:        1,
								ProtectionEndTime: time.Time{},
								DeletionTime:      time.Time{},
								Description:       "",
								Shield:            sdk.NewInt(1),
								ServiceFees:       OneMixedDecCoins(common.MicroCTKDenom),
							},
						},
					},
					endTime: time.Now().Add(12345),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			k.SetPool(suite.ctx, DummyPool(1))
			for _, pair := range tt.toInsert {
				k.InsertExpiringPurchaseQueue(suite.ctx, pair.purchaseList, pair.endTime)
			}
		})
	}
}

func TestKeeper_PurchaseShield(t *testing.T) {
	type args struct {
		ctx         sdk.Context
		poolID      uint64
		shield      sdk.Coins
		description string
		purchaser   sdk.AccAddress
		staking     bool
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    types.Purchase
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			got, err := k.PurchaseShield(tt.args.ctx, tt.args.poolID, tt.args.shield, tt.args.description, tt.args.purchaser, tt.args.staking)
			if (err != nil) != tt.wantErr {
				t.Errorf("PurchaseShield() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PurchaseShield() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_RemoveExpiredPurchasesAndDistributeFees(t *testing.T) {
	type args struct {
		ctx sdk.Context
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

func TestKeeper_SetExpiringPurchaseQueueTimeSlice(t *testing.T) {
	type args struct {
		ctx       sdk.Context
		timestamp time.Time
		ppPairs   []types.PoolPurchaser
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

func TestKeeper_SetNextPurchaseID(t *testing.T) {
	type args struct {
		ctx sdk.Context
		id  uint64
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

func TestKeeper_SetPurchaseList(t *testing.T) {
	type args struct {
		ctx          sdk.Context
		purchaseList types.PurchaseList
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
