package e2e

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (s *IntegrationTestSuite) testFeeGrant() {
	s.Run("test_fee_grant", func() {
		var (
			err      error
			valIdx   = 0
			c        = s.chainA
			endpoint = fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
		)

		alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()
		charlie, _ := c.genesisAccounts[3].keyInfo.GetAddress()

		amount := sdk.NewCoin(uctkDenom, math.NewInt(10000000))
		s.execFeeGrant(c, valIdx, alice.String(), bob.String(), amount, feesAmountCoin, false, withExtraFlag("allowed-messages", sdk.MsgTypeURL(&banktypes.MsgSend{})))

		var beforeAliceUctk, beforeBobUctk, beforeCharlieUctk sdk.Coin
		var afterAliceUctk, afterBobUctk, afterCharlieUctk sdk.Coin

		s.Require().Eventually(
			func() bool {
				beforeAliceUctk, err = queryShentuBalance(endpoint, alice.String(), uctkDenom)
				s.Require().NoError(err)
				beforeBobUctk, err = queryShentuBalance(endpoint, bob.String(), uctkDenom)
				s.Require().NoError(err)
				beforeCharlieUctk, err = queryShentuBalance(endpoint, charlie.String(), uctkDenom)
				s.Require().NoError(err)
				return beforeAliceUctk.IsValid() && beforeBobUctk.IsValid() && beforeCharlieUctk.IsValid()
			},
			20*time.Second,
			5*time.Second,
		)

		// Bob sends 10ctk to Charlie, Alice pays the fees
		s.execBankSend(c, valIdx, bob.String(), charlie.String(), amount, feesAmountCoin, false, withExtraFlag("fee-granter", alice.String()))
		s.Require().Eventually(
			func() bool {
				afterAliceUctk, err = queryShentuBalance(endpoint, alice.String(), uctkDenom)
				s.Require().NoError(err)
				afterBobUctk, err = queryShentuBalance(endpoint, bob.String(), uctkDenom)
				s.Require().NoError(err)
				afterCharlieUctk, err = queryShentuBalance(endpoint, charlie.String(), uctkDenom)
				s.Require().NoError(err)
				outgoing := beforeBobUctk.Sub(amount).IsEqual(afterBobUctk)
				incoming := beforeCharlieUctk.Add(amount).IsEqual(afterCharlieUctk)
				feepayment := beforeAliceUctk.Sub(feesAmountCoin).IsEqual(afterAliceUctk)
				return outgoing && incoming && feepayment
			},
			20*time.Second,
			5*time.Second,
		)
	})
}
