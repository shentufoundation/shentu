package keeper_test

import (
	"reflect"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/shield/keeper"
	"github.com/certikfoundation/shentu/x/shield/types"
)

func TestKeeper_GetSetProvider(t *testing.T) {
	type args struct {
		delegator sdk.AccAddress
	}
	tests := []struct {
		name         string
		args         args
		wantProvider types.Provider
		wantFound    bool
	}{
		{
			name: "Get Valid Address",
			args: args{
				delegator: acc2,
			},
			wantProvider: types.Provider{
				Address:          acc2.String(),
				DelegationBonded: sdk.NewInt(10000000000),
				Collateral:       sdk.ZeroInt(),
				TotalLocked:      sdk.ZeroInt(),
				Withdrawing:      sdk.ZeroInt(),
			},
			wantFound: true,
		},
		{
			name: "Get Invalid Address",
			args: args{
				delegator: acc2,
			},
			wantProvider: types.Provider{},
			wantFound:    false,
		},
		{
			name:         "Get Nil Address",
			args:         args{},
			wantProvider: types.Provider{},
			wantFound:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			if tt.wantFound {
				gotProvider := k.AddProvider(suite.ctx, tt.args.delegator)
				if !reflect.DeepEqual(gotProvider, tt.wantProvider) {
					t.Errorf("AddProvider() gotProvider = %v, want %v", gotProvider, tt.wantProvider)
				}
			}
			gotProvider, gotFound := k.GetProvider(suite.ctx, tt.args.delegator)
			if !reflect.DeepEqual(gotProvider, tt.wantProvider) {
				t.Errorf("GetProvider() gotProvider = %v, want %v", gotProvider, tt.wantProvider)
			}
			if gotFound != tt.wantFound {
				t.Errorf("GetProvider() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestKeeper_GetAllProviders(t *testing.T) {
	type args struct {
		providersToAdd []sdk.AccAddress
	}
	tests := []struct {
		name          string
		args          args
		wantProviders []types.Provider
	}{
		{
			name: "Empty Providers",
			args: args{
				providersToAdd: []sdk.AccAddress{},
			},
			wantProviders: []types.Provider{
				{
					Address:          acc1.String(),
					DelegationBonded: sdk.NewInt(10000000000),
					Collateral:       sdk.NewInt(2500000000),
					TotalLocked:      sdk.ZeroInt(),
					Withdrawing:      sdk.ZeroInt(),
				},
			},
		},
		{
			name: "Same Address",
			args: args{
				providersToAdd: []sdk.AccAddress{acc2, acc2},
			},
			wantProviders: []types.Provider{
				{
					Address:          acc1.String(),
					DelegationBonded: sdk.NewInt(10000000000),
					Collateral:       sdk.NewInt(2500000000),
					TotalLocked:      sdk.ZeroInt(),
					Withdrawing:      sdk.ZeroInt(),
				},
				{
					Address:          acc2.String(),
					DelegationBonded: sdk.NewInt(10000000000),
					Collateral:       sdk.ZeroInt(),
					TotalLocked:      sdk.ZeroInt(),
					Withdrawing:      sdk.ZeroInt(),
				},
			},
		},
		{
			name: "Different Addresses",
			args: args{
				providersToAdd: []sdk.AccAddress{acc2, acc3, acc4, acc5},
			},
			wantProviders: []types.Provider{
				{
					Address:          acc1.String(),
					DelegationBonded: sdk.NewInt(10000000000),
					Collateral:       sdk.NewInt(2500000000),
					TotalLocked:      sdk.ZeroInt(),
					Withdrawing:      sdk.ZeroInt(),
				},
				{
					Address:          acc2.String(),
					DelegationBonded: sdk.NewInt(10000000000),
					Collateral:       sdk.ZeroInt(),
					TotalLocked:      sdk.ZeroInt(),
					Withdrawing:      sdk.ZeroInt(),
				},
				{
					Address:          acc3.String(),
					DelegationBonded: sdk.NewInt(10000000000),
					Collateral:       sdk.ZeroInt(),
					TotalLocked:      sdk.ZeroInt(),
					Withdrawing:      sdk.ZeroInt(),
				},
				{
					Address:          acc4.String(),
					DelegationBonded: sdk.NewInt(10000000000),
					Collateral:       sdk.ZeroInt(),
					TotalLocked:      sdk.ZeroInt(),
					Withdrawing:      sdk.ZeroInt(),
				},
				{
					Address:          acc5.String(),
					DelegationBonded: sdk.NewInt(10000000000),
					Collateral:       sdk.ZeroInt(),
					TotalLocked:      sdk.ZeroInt(),
					Withdrawing:      sdk.ZeroInt(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup(t)
			k := suite.keeper
			for _, addr := range tt.args.providersToAdd {
				k.AddProvider(suite.ctx, addr)
			}
			sort.Slice(tt.wantProviders, func(i, j int) bool {
				return tt.wantProviders[i].Address < tt.wantProviders[j].Address
			})
			gotProviders := k.GetAllProviders(suite.ctx)
			sort.Slice(gotProviders, func(i, j int) bool {
				return gotProviders[i].Address < gotProviders[j].Address
			})
			if !reflect.DeepEqual(gotProviders, tt.wantProviders) {
				t.Errorf("GetAllProviders() = %v, want %v", gotProviders, tt.wantProviders)
			}
		})
	}
}

func TestKeeper_GetProvidersPaginated(t *testing.T) {
	suite := setup(t)
	k := suite.keeper
	for _, addr := range []sdk.AccAddress{acc1, acc2, acc3, acc4, acc5} {
		k.AddProvider(suite.ctx, addr)
	}
	providers := k.GetAllProviders(suite.ctx)
	if len(providers) != 5 {
		t.Errorf("GetProvidersPaginated() setup error")
	}
	type args struct {
		ctx   sdk.Context
		page  uint
		limit uint
	}
	tests := []struct {
		name          string
		keeper        keeper.Keeper
		args          args
		wantProviders []types.Provider
	}{
		{
			name:   "Valid Pagination: Page 1, Limit 3",
			keeper: k,
			args: args{
				ctx:   suite.ctx,
				page:  1,
				limit: 3,
			},
			wantProviders: providers[:3],
		},
		{
			name:   "Valid Pagination: Page 2, Limit 3",
			keeper: k,
			args: args{
				ctx:   suite.ctx,
				page:  2,
				limit: 3,
			},
			wantProviders: providers[3:],
		},
		{
			name:   "Invalid Page Number: Page 3, Limit 3",
			keeper: k,
			args: args{
				ctx:   suite.ctx,
				page:  3,
				limit: 3,
			},
			wantProviders: nil,
		},
		{
			name:   "Nil Page Number",
			keeper: k,
			args: args{
				ctx:   suite.ctx,
				limit: 0,
			},
			wantProviders: nil,
		},
		{
			name:   "Large Limit: Page 1, Limit 10",
			keeper: k,
			args: args{
				ctx:   suite.ctx,
				page:  1,
				limit: 10,
			},
			wantProviders: providers,
		},
		{
			name:   "Nil Limit",
			keeper: k,
			args: args{
				ctx:  suite.ctx,
				page: 1,
			},
			wantProviders: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.keeper
			if gotProviders := k.GetProvidersPaginated(tt.args.ctx, tt.args.page, tt.args.limit); !reflect.DeepEqual(gotProviders, tt.wantProviders) {
				t.Errorf("GetProvidersPaginated() = %v, want %v", gotProviders, tt.wantProviders)
			}
		})
	}
}

func TestKeeper_RemoveDelegation(t *testing.T) {
	type args struct {
		delAddr sdk.AccAddress
		valAddr sdk.ValAddress
	}
	tests := []struct {
		name         string
		args         args
		wantProvider types.Provider
		wantFound    bool
	}{
		{
			name: "Valid Remove Delegation",
			args: args{
				delAddr: acc2,
				valAddr: val2,
			},
			wantProvider: types.Provider{
				Address:          acc2.String(),
				DelegationBonded: sdk.ZeroInt(),
				Collateral:       sdk.ZeroInt(),
				TotalLocked:      sdk.ZeroInt(),
				Withdrawing:      sdk.ZeroInt(),
			},
			wantFound: true,
		},
		{
			name: "Invalid Remove: Address Mismatch",
			args: args{
				delAddr: acc2,
				valAddr: val1,
			},
			wantProvider: types.Provider{},
			wantFound:    false,
		},
		{
			name: "Invalid Remove: Nil Delegator Address",
			args: args{
				valAddr: val1,
			},
			wantProvider: types.Provider{},
			wantFound:    false,
		},
		{
			name: "Invalid Remove: Nil Validator Address",
			args: args{
				delAddr: acc2,
			},
			wantProvider: types.Provider{},
			wantFound:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && tt.wantFound {
					t.Errorf("RemoveDelegation() gotFound = %v, want %v", false, tt.wantProvider)
				}
			}()
			suite := setup(t)
			k := suite.keeper
			k.AddProvider(suite.ctx, tt.args.delAddr)
			k.RemoveDelegation(suite.ctx, tt.args.delAddr, tt.args.valAddr)
			gotProvider, gotFound := k.GetProvider(suite.ctx, tt.args.delAddr)
			if !reflect.DeepEqual(gotProvider, tt.wantProvider) {
				t.Errorf("RemoveDelegation() gotProvider = %v, want %v", gotProvider, tt.wantProvider)
			}
			if gotFound != tt.wantFound {
				t.Errorf("RemoveDelegation() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestKeeper_UpdateDelegationAmount(t *testing.T) {
	type args struct {
		delAddr sdk.AccAddress
	}
	tests := []struct {
		name         string
		args         args
		wantProvider types.Provider
		wantFound    bool
	}{
		{
			name: "Valid Update Delegation, Updated Provider",
			args: args{
				delAddr: acc5,
			},
			wantProvider: types.Provider{
				Address:          acc5.String(),
				DelegationBonded: sdk.NewInt(9000000000),
				Collateral:       sdk.ZeroInt(),
				TotalLocked:      sdk.ZeroInt(),
				Withdrawing:      sdk.ZeroInt(),
			},
			wantFound: true,
		},
		{
			name: "Valid Update Delegation, No Updates",
			args: args{
				delAddr: acc2,
			},
			wantProvider: types.Provider{
				Address:          acc2.String(),
				DelegationBonded: sdk.NewInt(10000000000),
				Collateral:       sdk.ZeroInt(),
				TotalLocked:      sdk.ZeroInt(),
				Withdrawing:      sdk.ZeroInt(),
			},
			wantFound: true,
		},
		{
			name:         "Invalid Update Delegation, Nil Address",
			args:         args{},
			wantProvider: types.Provider{},
			wantFound:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && tt.wantFound {
					t.Errorf("UpdateDelegationAmount() gotFound = %v, want %v", false, tt.wantProvider)
				}
			}()
			suite := setup(t)
			k := suite.keeper
			k.AddProvider(suite.ctx, tt.args.delAddr)
			suite.setupUndelegate()
			k.UpdateDelegationAmount(suite.ctx, tt.args.delAddr)
			gotProvider, gotFound := k.GetProvider(suite.ctx, tt.args.delAddr)
			if !reflect.DeepEqual(gotProvider, tt.wantProvider) {
				t.Errorf("UpdateDelegationAmount() gotProvider = %v, want %v", gotProvider, tt.wantProvider)
			}
			if gotFound != tt.wantFound {
				t.Errorf("UpdateDelegationAmount() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
