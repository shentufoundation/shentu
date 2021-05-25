package keeper_test

import (
	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"testing"
	"time"
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
		ctx     sdk.Context
		address sdk.AccAddress
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		wantRes []types.PurchaseList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if gotRes := k.GetPurchaserPurchases(tt.args.ctx, tt.args.address); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("GetPurchaserPurchases() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestKeeper_InsertExpiringPurchaseQueue(t *testing.T) {
	type args struct {
		ctx          sdk.Context
		purchaseList types.PurchaseList
		endTime      time.Time
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

func TestKeeper_IteratePoolPurchaseLists(t *testing.T) {
	type args struct {
		ctx      sdk.Context
		poolID   uint64
		callback func(purchaseList types.PurchaseList) (stop bool)
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

func TestKeeper_IteratePurchaseListEntries(t *testing.T) {
	type args struct {
		ctx      sdk.Context
		callback func(purchase types.Purchase) (stop bool)
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

func TestKeeper_IteratePurchaseLists(t *testing.T) {
	type args struct {
		ctx      sdk.Context
		callback func(purchase types.PurchaseList) (stop bool)
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
