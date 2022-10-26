package e2e

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *IntegrationTestSuite) TestIBCTokenTransfer() {
	var ibcStakeDenom string

	s.Run("send_photon_to_chainB", func() {
		recipient := s.chainB.validators[0].keyInfo.GetAddress().String()
		token := sdk.NewInt64Coin(photonDenom, 3300000000) // 3,300photon
		s.sendIBC(s.chainA.id, s.chainB.id, recipient, token)

		chainBAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainB.id][0].GetHostPort("1317/tcp"))

		// require the recipient account receives the IBC tokens (IBC packets ACKd)
		var (
			balances sdk.Coins
			err      error
		)
		s.Require().Eventually(
			func() bool {
				balances, err = queryShentuAllBalances(chainBAPIEndpoint, recipient)
				s.Require().NoError(err)

				return balances.Len() == 3
			},
			time.Minute,
			5*time.Second,
		)

		for _, c := range balances {
			if strings.Contains(c.Denom, "ibc/") {
				ibcStakeDenom = c.Denom
				s.Require().Equal(token.Amount.Int64(), c.Amount.Int64())
				break
			}
		}

		s.Require().NotEmpty(ibcStakeDenom)
	})
}

func (s *IntegrationTestSuite) TestStaking() {
	chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	validatorA := s.chainA.validators[0]
	validatorAAddr := validatorA.keyInfo.GetAddress()
	valOperA := sdk.ValAddress(validatorAAddr)

	alice := s.chainA.accounts[0].keyInfo.GetAddress()

	delegationAmount, _ := sdk.NewIntFromString("5000000")
	delegation := sdk.NewCoin(uctkDenom, delegationAmount)

	// Alice delegate uatom to Validator A
	s.executeDelegate(s.chainA, 0, delegation.String(), valOperA.String(), alice.String(), feesAmountCoin.String())

	// Validate delegation successful
	s.Require().Eventually(
		func() bool {
			res, err := queryDelegation(chainAAPIEndpoint, valOperA.String(), alice.String())
			amt := res.GetDelegationResponse().GetDelegation().GetShares()
			s.Require().NoError(err)

			return amt.Equal(sdk.NewDecFromInt(delegationAmount))
		},
		20*time.Second,
		5*time.Second,
	)
}

func (s *IntegrationTestSuite) TestGoverment() {
	chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	validatorA := s.chainA.validators[0]
	validatorAAddr := validatorA.keyInfo.GetAddress()

	height := s.getLatestBlockHeight(s.chainA, 0)
	proposalHeight := height + proposalBlockBuffer

	proposalCounter++
	s.T().Logf("Submiting upgrade proposal %d on chain %s", proposalCounter, s.chainA.id)
	s.executeSubmitUpgradeProposal(s.chainA, 0, proposalHeight, validatorAAddr.String(), "test-upgrade", feesAmountCoin.String())

	s.T().Logf("Voting upgrade proposal %d", proposalCounter)
	// First round, certifier vote
	s.executeVoteProposal(s.chainA, 0, validatorAAddr.String(), proposalCounter, "yes", feesAmountCoin.String())
	// Second round, validator vote
	s.executeVoteProposal(s.chainA, 0, validatorAAddr.String(), proposalCounter, "yes", feesAmountCoin.String())

	// Validate proposal status
	s.Require().Eventually(
		func() bool {
			status, err := queryProposal(chainAAPIEndpoint, proposalCounter)
			s.Require().NoError(err)
			return status == "PROPOSAL_STATUS_VALIDATOR_VOTING_PERIOD"
		},
		20*time.Second,
		5*time.Second,
	)
}

func (s *IntegrationTestSuite) TestShieldCreatePool() {
	chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	validatorA := s.chainA.validators[0]
	validatorAAddr := validatorA.keyInfo.GetAddress()

	// First, deposit collaterals
	s.executeDepositCollateral(s.chainA, 0, validatorAAddr.String(), depositAmountCoin.String(), feesAmountCoin.String())
	// Second, create pool
	shieldPoolCounter++
	s.T().Logf("Creating shield pool %d on chain %s", shieldPoolCounter, s.chainA.id)
	s.executeCreatePool(s.chainA, 0, depositAmountCoin.String(), shieldPoolName, validatorAAddr.String(), shieldPoolLimit, validatorAAddr.String(), depositAmountCoin.String(), feesAmountCoin.String())

	// Validate pool status
	s.Require().Eventually(
		func() bool {
			status, err := queryShieldPool(chainAAPIEndpoint, proposalCounter)
			s.Require().NoError(err)
			return status.Pool.Active
		},
		20*time.Second,
		5*time.Second,
	)
}
