package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/oracle/internal/types"
)

type Keeper struct {
	cdc           *codec.Codec
	storeKey      sdk.StoreKey
	authKeeper    types.AuthKeeper
	distrKeeper   types.DistrKeeper
	stakingKeeper types.StakingKeeper
	supplyKeeper  types.SupplyKeeper
	CertKeeper    types.CertKeeper
	paramSpace    types.ParamSubspace
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, authKeeper types.AuthKeeper, certKeeper types.CertKeeper, distriKeeper types.DistrKeeper,
	stakingKeeper types.StakingKeeper, supplyKeeper types.SupplyKeeper, paramSpace types.ParamSubspace) Keeper {
	return Keeper{
		cdc:           cdc,
		paramSpace:    paramSpace.WithKeyTable(types.ParamKeyTable()),
		storeKey:      storeKey,
		authKeeper:    authKeeper,
		distrKeeper:   distriKeeper,
		stakingKeeper: stakingKeeper,
		supplyKeeper:  supplyKeeper,
		CertKeeper:    certKeeper,
	}
}

// GetAuthKeeper returns the auth keeper wrapped in module keeper.
func (k Keeper) GetAuthKeeper() types.AuthKeeper {
	return k.authKeeper
}
