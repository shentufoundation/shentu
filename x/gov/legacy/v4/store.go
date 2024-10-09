package v4

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/exported"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	v1 "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

func MigrateCustomParams(ctx sdk.Context, storeKey storetypes.StoreKey, legacySubspace exported.ParamSubspace, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	cp := v1.CustomParams{}
	legacySubspace.Get(ctx, v1.ParamStoreKeyCustomParams, &cp)

	securityVoteTally := cp.CertifierUpdateSecurityVoteTally
	stakeVoteTally := cp.CertifierUpdateStakeVoteTally
	// customParams
	certifierUpdateSecurityVoteTally := govtypesv1.NewTallyParams(
		securityVoteTally.Quorum,
		securityVoteTally.Threshold,
		securityVoteTally.VetoThreshold,
	)
	certifierUpdateStakeVoteTally := govtypesv1.NewTallyParams(
		stakeVoteTally.Quorum,
		stakeVoteTally.Threshold,
		stakeVoteTally.VetoThreshold,
	)
	customParams := v1.CustomParams{
		CertifierUpdateSecurityVoteTally: &certifierUpdateSecurityVoteTally,
		CertifierUpdateStakeVoteTally:    &certifierUpdateStakeVoteTally,
	}

	bz, err := cdc.Marshal(&customParams)
	if err != nil {
		return err
	}

	// set migrate params
	store.Set(CustomParamsKey, bz)

	return nil
}
