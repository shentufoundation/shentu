package e2e

import (
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testDistribution() {
	s.Run("test_distribution", func() {
		var (
			err      error
			valIdx   = 0
			c        = s.chainA
			endpoint = fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
		)

		delegator, _ := c.genesisAccounts[2].keyInfo.GetAddress()
		withdrawer, _ := c.genesisAccounts[3].keyInfo.GetAddress()
		s.execSetWithdrawAddress(c, valIdx, delegator.String(), withdrawer.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				res, err := queryDelegatorWithdrawalAddress(endpoint, delegator.String())
				s.Require().NoError(err)
				return res.WithdrawAddress == withdrawer.String()
			},
			20*time.Second,
			5*time.Second,
		)

		// var beforeBalance, afterBalance sdk.Coin

		// s.Require().Eventually(
		// 	func() bool {
		// 		beforeBalance, err = queryShentuBalance(endpoint, delegator.String(), uctkDenom)
		// 		s.Require().NoError(err)
		// 		return beforeBalance.IsValid()
		// 	},
		// 	20*time.Second,
		// 	5*time.Second,
		// )

		// valaddr, _ := c.validators[0].keyInfo.GetAddress()
		// validator := sdk.ValAddress(valaddr)
		// s.execWithdrawReward(c, valIdx, delegator.String(), validator.String(), feesAmountCoin, false)
		// s.Require().Eventually(
		// 	func() bool {
		// 		afterBalance, err = queryShentuBalance(endpoint, delegator.String(), uctkDenom)
		// 		s.Require().NoError(err)
		// 		return afterBalance.IsGTE(beforeBalance.Sub(feesAmountCoin))
		// 	},
		// 	20*time.Second,
		// 	5*time.Second,
		// )

		var beforeDistribBalance, afterDistribBalance sdk.Coin

		s.Require().Eventually(
			func() bool {
				beforeDistribBalance, err = queryShentuBalance(endpoint, distribModuleAcct.String(), uctkDenom)
				s.Require().NoError(err)
				return beforeDistribBalance.IsValid()
			},
			20*time.Second,
			5*time.Second,
		)
		amount := sdk.NewCoin(uctkDenom, math.NewInt(10000000))
		s.execFundCommunityPool(c, valIdx, delegator.String(), amount, feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				afterDistribBalance, err = queryShentuBalance(endpoint, distribModuleAcct.String(), uctkDenom)
				s.Require().NoError(err)
				return afterDistribBalance.IsGTE(beforeDistribBalance.Add(amount))
			},
			20*time.Second,
			5*time.Second,
		)
	})
}
