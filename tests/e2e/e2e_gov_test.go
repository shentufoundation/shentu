package e2e

import (
	"fmt"
	"path/filepath"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkgovtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

func (s *IntegrationTestSuite) testCommonProposal() {
	s.Run("test_common_proposal", func() {
		var (
			valIdx       = 0
			c            = s.chainA
			grpcEndpoint = s.valResources[s.chainA.id][0].GetHostPort("9090/tcp")
		)

		valA, _ := c.validators[0].keyInfo.GetAddress()
		valB, _ := c.validators[1].keyInfo.GetAddress()

		alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()

		// Create a proposal to send 10ctk from alice to bob
		amount := sdk.NewCoin(uctkDenom, math.NewInt(10000000))
		s.writePoolSpendProposal(c, alice.String(), "pool_proposal.json", amount)

		s.execSubmitProposal(c, valIdx, "pool_proposal.json", bob.String(), feesAmountCoin)
		s.Require().Eventually(
			func() bool {
				proposal, err := queryProposal(grpcEndpoint, proposalCounter)
				return err == nil && proposal.Id == proposalCounter && proposal.Status == sdkgovtypes.StatusVotingPeriod
			},
			20*time.Second,
			5*time.Second,
		)

		s.execVoteProposal(c, 0, proposalCounter, valA.String(), "yes", feesAmountCoin)
		s.execVoteProposal(c, 1, proposalCounter, valB.String(), "yes", feesAmountCoin)
		s.Require().Eventually(
			func() bool {
				proposal, err := queryProposal(grpcEndpoint, proposalCounter)
				return err == nil && proposal.Id == proposalCounter && proposal.Status == sdkgovtypes.StatusPassed
			},
			20*time.Second,
			5*time.Second,
		)
	})
}

func (s *IntegrationTestSuite) writePoolSpendProposal(c *chain, recipient, fileName string, amount sdk.Coin) {
	template := `{
		"messages": [{
			"@type": "/cosmos.distribution.v1beta1.MsgCommunityPoolSpend",
			"authority": "%s",
			"recipient": "%s",
			"amount": [{
				"denom": "%s",
				"amount": "%s"
			}]
		}],
		"metadata": "community pool spend",
		"deposit": "512000000uctk",
		"title": "community pool proposal",
		"summary": "community pool summary"
	}`
	body := fmt.Sprintf(template, govModuleAcct.String(), recipient, amount.Denom, amount.Amount.String())
	err := writeFile(filepath.Join(c.validators[0].configDir(), "config", fileName), []byte(body))
	s.Require().NoError(err)
}
