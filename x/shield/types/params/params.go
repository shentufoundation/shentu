package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/certikfoundation/shentu/v2/x/shield/types/v1alpha1"
	"github.com/certikfoundation/shentu/v2/x/shield/types/v1beta1"
)

// ParamKeyTable is the key declaration for parameters.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(v1alpha1.ParamStoreKeyPoolParams, v1alpha1.PoolParams{}, v1alpha1.ValidatePoolParams),
		paramtypes.NewParamSetPair(v1beta1.ParamStoreKeyPoolParams, v1beta1.PoolParams{}, v1beta1.ValidatePoolParams),
		paramtypes.NewParamSetPair(v1beta1.ParamStoreKeyClaimProposalParams, v1beta1.ClaimProposalParams{}, v1beta1.ValidateClaimProposalParams),
		paramtypes.NewParamSetPair(v1beta1.ParamStoreKeyStakingShieldRate, sdk.Dec{}, v1beta1.ValidateStakingShieldRateParams),
		paramtypes.NewParamSetPair(v1beta1.ParamStoreKeyBlockRewardParams, v1beta1.BlockRewardParams{}, v1beta1.ValidateBlockRewardParams),
	)
}
