package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

// nolint
const (
	subkeyProtectionPeriod  = "protection_period"
	subkeyWithdrawPeriod    = "withdraw_period"
	subkeyShieldFeesRate    = "shield_fees_rate"
	subkeyPoolShieldLimit   = "pool_shield_limit"
	subkeyMinShieldPurchase = "min_shield_purchase"

	subkeyClaimPeriod  = "claim_period"
	subkeyPayoutPeriod = "payout_period"
	subkeyMinDeposit   = "min_deposit"
	subkeyDepositRate  = "deposit_rate"
	subkeyFeesRate     = "fees_rate"
)

// ParamChanges defines the parameters that can be modified by param change proposals on the simulation.
func ParamChanges(_ *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyPoolParams),
			func(r *rand.Rand) string {
				pp := GenPoolParams(r)
				changes := []struct {
					key   string
					value string
				}{
					{subkeyProtectionPeriod, fmt.Sprintf("%d", pp.ProtectionPeriod)},
					{subkeyWithdrawPeriod, fmt.Sprintf("%d", pp.WithdrawPeriod)},
					{subkeyShieldFeesRate, pp.ShieldFeesRate.String()},
					{subkeyPoolShieldLimit, pp.PoolShieldLimit.String()},
				}

				pc := make(map[string]string)
				numChanges := len(changes)
				for i := 0; i < numChanges; i++ {
					c := changes[i]
					pc[c.key] = c.value
				}
				bz, _ := json.Marshal(pc)
				return string(bz)
			},
		),

		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyClaimProposalParams),
			func(r *rand.Rand) string {
				cpp := GenClaimProposalParams(r)
				changes := []struct {
					key   string
					value string
				}{
					{subkeyClaimPeriod, fmt.Sprintf("%d", cpp.ClaimPeriod)},
					{subkeyPayoutPeriod, fmt.Sprintf("%d", cpp.PayoutPeriod)},
					{subkeyDepositRate, cpp.DepositRate.String()},
					{subkeyFeesRate, cpp.FeesRate.String()},
				}

				pc := make(map[string]string)
				numChanges := len(changes)
				for i := 0; i < numChanges; i++ {
					c := changes[i]
					pc[c.key] = c.value
				}
				bz, _ := json.Marshal(pc)
				return string(bz)
			},
		),

		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyStakingShieldRate),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenShieldStakingRateParam(r))
			},
		),
	}
}
