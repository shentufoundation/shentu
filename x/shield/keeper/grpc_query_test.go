package keeper_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/types"
)

func TestKeeper_ClaimParams(t *testing.T) {
	emptyParams := types.ClaimProposalParams{
		ClaimPeriod:  0,
		PayoutPeriod: 0,
		MinDeposit:   sdk.NewCoins(sdk.NewInt64Coin("stake", 1234)),
		DepositRate:  sdk.NewDec(1),
		FeesRate:     sdk.NewDec(1),
	}
	type args struct {
		params types.ClaimProposalParams
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryClaimParamsResponse
		wantErr bool
	}{
		{
			name: "Default params",
			args: args{
				params: types.DefaultClaimProposalParams(),
			},
			want: &types.QueryClaimParamsResponse{
				Params: types.DefaultClaimProposalParams(),
			},
		},
		{
			name: "Empty params",
			args: args{
				params: emptyParams,
			},
			want: &types.QueryClaimParamsResponse{
				Params: emptyParams,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			k.SetClaimProposalParams(suite.ctx, tt.args.params)
			got, err := k.ClaimParams(sdk.WrapSDKContext(suite.ctx), &types.QueryClaimParamsRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("ClaimParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClaimParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Pool(t *testing.T) {
	type args struct {
		poolsToAdd []types.Pool
		req        *types.QueryPoolRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryPoolResponse
		wantErr bool
	}{
		{
			name:    "No Pool",
			args:    args{req: &types.QueryPoolRequest{PoolId: 1}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "One Pool",
			args: args{
				[]types.Pool{{
					Id:          1,
					Description: "w",
					Sponsor:     "w",
					SponsorAddr: "w",
					ShieldLimit: sdk.NewInt(1),
					Active:      false,
					Shield:      sdk.NewInt(1),
				}},
				&types.QueryPoolRequest{
					PoolId: 1,
				},
			},
			want: &types.QueryPoolResponse{
				Pool: types.Pool{
					Id:          1,
					Description: "w",
					Sponsor:     "w",
					SponsorAddr: "w",
					ShieldLimit: sdk.NewInt(1),
					Active:      false,
					Shield:      sdk.NewInt(1),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			for _, p := range tt.args.poolsToAdd {
				k.SetPool(suite.ctx, p)
			}
			got, err := k.Pool(sdk.WrapSDKContext(suite.ctx), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pool() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_PoolParams(t *testing.T) {
	randomParams := types.PoolParams{
		ProtectionPeriod:  0,
		ShieldFeesRate:    sdk.NewDec(123),
		WithdrawPeriod:    0,
		PoolShieldLimit:   sdk.NewDec(1234),
		MinShieldPurchase: sdk.NewCoins(sdk.NewInt64Coin("stake", 12345)),
	}
	type args struct {
		params types.PoolParams
		req    *types.QueryPoolParamsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryPoolParamsResponse
		wantErr bool
	}{
		{
			name: "Default Params",
			args: args{
				params: types.DefaultPoolParams(),
				req:    &types.QueryPoolParamsRequest{},
			},
			want: &types.QueryPoolParamsResponse{
				Params: types.DefaultPoolParams(),
			},
			wantErr: false,
		},
		{
			name: "Random Params",
			args: args{
				params: randomParams,
				req:    &types.QueryPoolParamsRequest{},
			},
			want: &types.QueryPoolParamsResponse{
				Params: randomParams,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			k.SetPoolParams(suite.ctx, tt.args.params)
			got, err := k.PoolParams(sdk.WrapSDKContext(suite.ctx), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PoolParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PoolParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_PoolPurchaseLists(t *testing.T) {
	pl := types.PurchaseList{
		PoolId:    1,
		Purchaser: acc3.String(),
		Entries: []types.Purchase{{
			PurchaseId:        0,
			ProtectionEndTime: time.Time{},
			DeletionTime:      time.Time{},
			Description:       "",
			Shield:            sdk.NewInt(1),
			ServiceFees:       OneMixedDecCoins("uctk"),
		}},
	}
	type args struct {
		pls []types.PurchaseList
		req *types.QueryPoolPurchaseListsRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryPurchaseListsResponse
		wantErr bool
	}{
		{
			name: "Empty lists",
			args: args{
				req: &types.QueryPoolPurchaseListsRequest{
					PoolId: 1,
				},
			},
			want:    &types.QueryPurchaseListsResponse{},
			wantErr: false,
		},
		{
			name: "One lists",
			args: args{
				pls: []types.PurchaseList{pl},
				req: &types.QueryPoolPurchaseListsRequest{
					PoolId: 1,
				},
			},
			want: &types.QueryPurchaseListsResponse{
				PurchaseLists: []types.PurchaseList{pl},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			p := DummyPool(1)
			k := suite.keeper
			k.SetPool(suite.ctx, p)
			for _, pl := range tt.args.pls {
				k.SetPurchaseList(suite.ctx, pl)
			}
			got, err := k.PoolPurchaseLists(sdk.WrapSDKContext(suite.ctx), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PoolPurchaseLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PoolPurchaseLists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Pools(t *testing.T) {
	type args struct {
		pools []types.Pool
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryPoolsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			got, err := k.Pools(sdk.WrapSDKContext(suite.ctx), &types.QueryPoolsRequest{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Pools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pools() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Provider(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryProviderRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryProviderResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.Provider(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Providers(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryProvidersRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryProvidersResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.Providers(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Providers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Providers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_PurchaseList(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryPurchaseListRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryPurchaseListResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.PurchaseList(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PurchaseList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PurchaseList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_PurchaseLists(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryPurchaseListsRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryPurchaseListsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.PurchaseLists(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PurchaseLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PurchaseLists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Purchases(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryPurchasesRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryPurchasesResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.Purchases(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Purchases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Purchases() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Reimbursement(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryReimbursementRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryReimbursementResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.Reimbursement(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reimbursement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reimbursement() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Reimbursements(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryReimbursementsRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryReimbursementsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.Reimbursements(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reimbursements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reimbursements() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_ShieldStaking(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryShieldStakingRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryShieldStakingResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.ShieldStaking(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShieldStaking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShieldStaking() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_ShieldStakingRate(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryShieldStakingRateRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryShieldStakingRateResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.ShieldStakingRate(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShieldStakingRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShieldStakingRate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_ShieldStatus(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QueryShieldStatusRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QueryShieldStatusResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.ShieldStatus(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShieldStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShieldStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_Sponsor(t *testing.T) {
	type args struct {
		c   context.Context
		req *types.QuerySponsorRequest
	}
	tests := []struct {
		name    string
		keeper  keeper.Keeper
		args    args
		want    *types.QuerySponsorResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.keeper
			got, err := q.Sponsor(tt.args.c, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sponsor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sponsor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
