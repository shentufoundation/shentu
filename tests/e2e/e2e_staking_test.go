package e2e

import (
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testStaking() {
	s.Run("test_staking", func() {
		var (
			err      error
			valIdx   = 0
			c        = s.chainA
			endpoint = fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
		)

		valAaddr, _ := c.validators[0].keyInfo.GetAddress()
		// valBaddr, _ := c.validators[1].keyInfo.GetAddress()

		validatorA := sdk.ValAddress(valAaddr)
		// validatorB := sdk.ValAddress(valBaddr)

		delegator, _ := c.genesisAccounts[2].keyInfo.GetAddress()

		currentDelegation := math.LegacyZeroDec()
		resp, err := queryDelegation(endpoint, validatorA.String(), delegator.String())
		if err == nil {
			currentDelegation = resp.DelegationResponse.Delegation.Shares
		}

		amount := sdk.NewCoin(uctkDenom, math.NewInt(100000000))
		s.execDelegate(c, valIdx, delegator.String(), validatorA.String(), amount, feesAmountCoin, false)

		s.Require().Eventually(
			func() bool {
				resp, err := queryDelegation(endpoint, validatorA.String(), delegator.String())
				s.Require().NoError(err)
				return resp.DelegationResponse.Delegation.Shares.Equal(currentDelegation.Add(math.LegacyNewDecFromInt(amount.Amount)))
			},
			20*time.Second,
			5*time.Second,
		)
	})
}
