package e2e

import (
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) testBankTokenTransfer() {
	s.Run("send_uctk_between_accounts", func() {
		var (
			err      error
			valIdx   = 0
			c        = s.chainA
			endpoint = fmt.Sprintf("http://%s", s.valResources[c.id][valIdx].GetHostPort("1317/tcp"))
		)

		alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()
		charlie, _ := c.genesisAccounts[3].keyInfo.GetAddress()

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

		// Alice sends 10ctk to Bob
		amountCoin := sdk.NewCoin(uctkDenom, math.NewInt(10000000))

		s.execBankSend(c, valIdx, alice.String(), bob.String(), amountCoin, feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				afterAliceUctk, err = queryShentuBalance(endpoint, alice.String(), uctkDenom)
				s.Require().NoError(err)
				afterBobUctk, err = queryShentuBalance(endpoint, bob.String(), uctkDenom)
				s.Require().NoError(err)
				outgoing := beforeAliceUctk.Sub(amountCoin).Sub(feesAmountCoin).IsEqual(afterAliceUctk)
				incoming := beforeBobUctk.Add(amountCoin).IsEqual(afterBobUctk)
				return outgoing && incoming
			},
			20*time.Second,
			5*time.Second,
		)

		beforeAliceUctk, beforeBobUctk = afterAliceUctk, afterBobUctk

		// alice sends 10ctk to bob and charlie
		s.execBankMultiSend(c, valIdx, alice.String(), []string{bob.String(), charlie.String()}, amountCoin, feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				afterAliceUctk, err = queryShentuBalance(endpoint, alice.String(), uctkDenom)
				s.Require().NoError(err)
				afterBobUctk, err = queryShentuBalance(endpoint, bob.String(), uctkDenom)
				s.Require().NoError(err)
				afterCharlieUctk, err = queryShentuBalance(endpoint, charlie.String(), uctkDenom)
				s.Require().NoError(err)
				outgoing := beforeAliceUctk.Sub(amountCoin).Sub(amountCoin).Sub(feesAmountCoin).IsEqual(afterAliceUctk)
				incoming := beforeBobUctk.Add(amountCoin).IsEqual(afterBobUctk) && beforeCharlieUctk.Add(amountCoin).IsEqual(afterCharlieUctk)
				return outgoing && incoming
			},
			20*time.Second,
			5*time.Second,
		)
	})
}
