package e2e

import (
	"bytes"
	"encoding/hex"
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
	s.Run("delegate_staking", func() {
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
	})

	s.Run("unbond_staking", func() {
		chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
		validatorA := s.chainA.validators[0]
		validatorAAddr := validatorA.keyInfo.GetAddress()
		valOperA := sdk.ValAddress(validatorAAddr)

		alice := s.chainA.accounts[0].keyInfo.GetAddress()

		delegationAmount, _ := sdk.NewIntFromString("5000000")
		unbondAmount, _ := sdk.NewIntFromString("500000")
		unbond := sdk.NewCoin(uctkDenom, unbondAmount)

		// Alice unbond uatom to Validator A
		s.executeUnbond(s.chainA, 0, unbond.String(), valOperA.String(), alice.String(), feesAmountCoin.String())

		// Validate unbond successful
		s.Require().Eventually(
			func() bool {
				res, err := queryDelegation(chainAAPIEndpoint, valOperA.String(), alice.String())
				amt := res.GetDelegationResponse().GetDelegation().GetShares()
				s.Require().NoError(err)
				return amt.Equal(sdk.NewDecFromInt(delegationAmount.Sub(unbondAmount)))
			},
			20*time.Second,
			5*time.Second,
		)
	})
}

func (s *IntegrationTestSuite) TestGovernance() {

	chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	validatorA := s.chainA.validators[0]
	validatorAAddr := validatorA.keyInfo.GetAddress()

	s.Run("submit_upgrade_proposal", func() {
		height := s.getLatestBlockHeight(s.chainA, 0)
		proposalHeight := height + proposalBlockBuffer

		proposalCounter++
		s.T().Logf("Submiting upgrade proposal %d on chain %s", proposalCounter, s.chainA.id)
		s.executeSubmitUpgradeProposal(s.chainA, 0, proposalHeight, validatorAAddr.String(), "test-upgrade", feesAmountCoin.String())

		s.Require().Eventually(
			func() bool {
				res, err := queryProposal(chainAAPIEndpoint, proposalCounter)
				s.Require().NoError(err)
				return res.Proposal.ProposalId == uint64(proposalCounter)
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("voting_proposal", func() {
		s.T().Logf("Voting upgrade proposal %d", proposalCounter)
		// First round, certifier vote
		s.executeVoteProposal(s.chainA, 0, validatorAAddr.String(), proposalCounter, "yes", feesAmountCoin.String())
		// Second round, validator vote
		s.executeVoteProposal(s.chainA, 0, validatorAAddr.String(), proposalCounter, "yes", feesAmountCoin.String())

		// Validate proposal status
		s.Require().Eventually(
			func() bool {
				res, err := queryProposal(chainAAPIEndpoint, proposalCounter)
				s.Require().NoError(err)
				return res.Proposal.Status == 3
			},
			20*time.Second,
			5*time.Second,
		)
	})
}

func (s *IntegrationTestSuite) TestCoreShield() {
	chainAAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[s.chainA.id][0].GetHostPort("1317/tcp"))
	validatorA := s.chainA.validators[0]
	validatorAAddr := validatorA.keyInfo.GetAddress()

	accountA := s.chainA.accounts[0]
	accountAAddr := accountA.keyInfo.GetAddress()

	// Deposit collaterals
	s.Run("deposit_collateral", func() {
		s.executeDepositCollateral(s.chainA, 0, validatorAAddr.String(), collateralAmountCoin.String(), feesAmountCoin.String())

		// Check shield status
		s.Require().Eventually(
			func() bool {
				status, err := queryShieldStatus(chainAAPIEndpoint)
				s.Require().NoError(err)
				return status.TotalCollateral.GTE(depositAmount)
			},
			20*time.Second,
			5*time.Second,
		)
	})

	// Create pool
	s.Run("create_pool", func() {
		shieldPoolCounter++
		shieldPurchaseCounter++
		s.T().Logf("Creating shield pool %d on chain %s", shieldPoolCounter, s.chainA.id)
		s.executeCreatePool(s.chainA, 0, shieldAmountCoin.String(), shieldPoolName, validatorAAddr.String(), shieldPoolLimit, validatorAAddr.String(), depositAmountCoin.String(), feesAmountCoin.String())

		// Validate pool status
		s.Require().Eventually(
			func() bool {
				status, err := queryShieldPool(chainAAPIEndpoint, shieldPoolCounter)
				s.Require().NoError(err)
				return status.Pool.Active
			},
			20*time.Second,
			5*time.Second,
		)
	})

	// Buy shield
	s.Run("purchase_shield", func() {
		shieldPurchaseCounter++
		s.T().Logf("Purchasing shield pool %d on chain %s", shieldPoolCounter, s.chainA.id)
		s.executePurchaseShield(s.chainA, 0, shieldPoolCounter, shieldAmountCoin.String(), "'shield desc'", accountAAddr.String(), feesAmountCoin.String())

		// Validate purchase status
		s.Require().Eventually(
			func() bool {
				res, err := queryShieldPurchase(chainAAPIEndpoint, accountAAddr.String(), shieldPoolCounter)
				s.Require().NoError(err)
				return len(res.PurchaseList.Entries) >= 1
			},
			20*time.Second,
			5*time.Second,
		)
	})

	// Issue identity certificate
	s.Run("issue_identity", func() {
		certificateCounter++
		s.T().Logf("Issue identity certificate to %s on chain %s", validatorAAddr.String(), s.chainA.id)
		s.executeIssueCertificate(s.chainA, 0, "identity", validatorAAddr.String(), validatorAAddr.String(), feesAmountCoin.String())

		// Validate certificate status
		s.Require().Eventually(
			func() bool {
				res, err := queryCertificate(chainAAPIEndpoint, certificateCounter)
				s.Require().NoError(err)
				return bytes.Contains(res.Certificate.Content.GetValue(), []byte(validatorAAddr.String()))
			},
			20*time.Second,
			5*time.Second,
		)
	})

	// Submit claim
	s.Run("submit_claim", func() {
		proposalFile := "test_claim.json"
		s.T().Logf("Submiting claim from %s on chain %s", accountAAddr.String(), s.chainA.id)
		s.writeClaimProposal(s.chainA, 0, shieldPoolCounter, shieldPurchaseCounter, proposalFile)
		proposalCounter++
		s.executeSubmitClaimProposal(s.chainA, 0, configFile(proposalFile), accountAAddr.String(), feesAmountCoin.String())

		s.Require().Eventually(
			func() bool {
				res, err := queryProposal(chainAAPIEndpoint, proposalCounter)
				s.Require().NoError(err)
				return res.Proposal.ProposalId == uint64(proposalCounter)
			},
			20*time.Second,
			5*time.Second,
		)
	})

	// Vote claim
	s.Run("vote_claim", func() {
		s.T().Logf("Voting claim proposal %d", proposalCounter)
		// First round, certifier vote
		s.executeVoteProposal(s.chainA, 0, validatorAAddr.String(), proposalCounter, "yes", feesAmountCoin.String())
		// Second round, validator vote
		s.executeVoteProposal(s.chainA, 0, validatorAAddr.String(), proposalCounter, "yes", feesAmountCoin.String())

		// Validate proposal status
		s.Require().Eventually(
			func() bool {
				res, err := queryProposal(chainAAPIEndpoint, proposalCounter)
				s.Require().NoError(err)
				return res.Proposal.Status == 4
			},
			20*time.Second,
			5*time.Second,
		)
	})

	// Check reimbursement
	s.Run("check_reimbursement", func() {
		s.T().Logf("Check reimbursement to %s", accountAAddr.String())

		s.Require().Eventually(
			func() bool {
				res, err := queryShieldReimbursement(chainAAPIEndpoint, proposalCounter)
				s.Require().NoError(err)
				return strings.Contains(res.Reimbursement.Beneficiary, accountAAddr.String())
			},
			20*time.Second,
			5*time.Second,
		)
	})
}

func (s *IntegrationTestSuite) TestOracle() {
	chainAAPIEndpoint := s.valResources[s.chainA.id][0].GetHostPort("9090/tcp")

	alice := s.chainA.accounts[0].keyInfo.GetAddress()
	bob := s.chainA.accounts[1].keyInfo.GetAddress()
	charle := s.chainA.accounts[2].keyInfo.GetAddress()

	var txHash, taskHash string
	var txHash2, taskHash2 string
	var txHash3, taskHash3 string
	var err error

	valTime := time.Now().Add(120 * time.Second)
	valTimeStr := valTime.Format(time.RFC3339)

	longValTime := time.Now().Add(1200 * time.Second)
	longValTimeStr := longValTime.Format(time.RFC3339)

	taskContract := "demo-contract"
	taskFunction := "demo-function"

	s.Run("create_operator", func() {
		lessCollateral := sdk.NewCoin(uctkDenom, sdk.NewInt(50000000))
		collateral := sdk.NewCoin(uctkDenom, sdk.NewInt(100000000))
		mostCollateral := sdk.NewCoin(uctkDenom, sdk.NewInt(250000000))

		s.executeOracleCreateOperator(s.chainA, 0, alice.String(), collateral.String(), feesAmountCoin.String())
		s.executeOracleCreateOperator(s.chainA, 0, bob.String(), lessCollateral.String(), feesAmountCoin.String())
		s.executeOracleCreateOperator(s.chainA, 0, charle.String(), mostCollateral.String(), feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleOperator(chainAAPIEndpoint, charle.String())
				s.Require().NoError(e)
				return res.Operator.Address == charle.String()
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("create_tx_task", func() {
		txBytes := base64.StdEncoding.EncodeToString([]byte(valTimeStr + "1"))
		txBytes2 := base64.StdEncoding.EncodeToString([]byte(valTimeStr + "2"))
		txBytes3 := base64.StdEncoding.EncodeToString([]byte(valTimeStr + "3"))
		chainID := "test"
		bountyAmount, _ := sdk.NewIntFromString("500000")
		bounty := sdk.NewCoin(uctkDenom, bountyAmount)
		// normal tx task
		txHash, err = s.executeOracleCreateTxTask(s.chainA, 0, txBytes, chainID, bounty.String(), valTimeStr, alice.String(), feesAmountCoin.String())
		s.Require().NoError(err)
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTaskHash(chainAAPIEndpoint, txHash)
				if e == nil {
					taskHash = res
					return true
				}
				return false
			},
			20*time.Second,
			5*time.Second,
		)
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTxTask(chainAAPIEndpoint, taskHash)
				s.Require().NoError(e)
				return res.Task.Status == 1
			},
			20*time.Second,
			5*time.Second,
		)
		// 0 score task
		txHash2, err = s.executeOracleCreateTxTask(s.chainA, 0, txBytes2, chainID, bounty.String(), valTimeStr, alice.String(), feesAmountCoin.String())
		s.Require().NoError(err)
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTaskHash(chainAAPIEndpoint, txHash2)
				if e == nil {
					taskHash2 = res
					return true
				}
				return false
			},
			20*time.Second,
			5*time.Second,
		)
		// long term task
		txHash3, err = s.executeOracleCreateTxTask(s.chainA, 0, txBytes3, chainID, bounty.String(), longValTimeStr, alice.String(), feesAmountCoin.String())
		s.Require().NoError(err)
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTaskHash(chainAAPIEndpoint, txHash3)
				if e == nil {
					taskHash3 = res
					return true
				}
				return false
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("respond_tx_task", func() {
		s.executeOracleRespondTxTask(s.chainA, 0, 90, taskHash, alice.String(), feesAmountCoin.String())
		s.executeOracleRespondTxTask(s.chainA, 0, 90, taskHash, bob.String(), feesAmountCoin.String())
		s.executeOracleRespondTxTask(s.chainA, 0, 60, taskHash, charle.String(), feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTxTask(chainAAPIEndpoint, taskHash)
				s.Require().NoError(e)
				return len(res.Task.Responses) == 3
			},
			20*time.Second,
			5*time.Second,
		)

		s.executeOracleRespondTxTask(s.chainA, 0, 0, taskHash2, alice.String(), feesAmountCoin.String())
		s.executeOracleRespondTxTask(s.chainA, 0, 0, taskHash2, bob.String(), feesAmountCoin.String())

		s.executeOracleRespondTxTask(s.chainA, 0, 70, taskHash3, charle.String(), feesAmountCoin.String())
	})

	s.Run("close_tx_task", func() {
		if time.Now().Before(valTime) {
			time.Sleep(time.Until(valTime))
		}
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTxTask(chainAAPIEndpoint, taskHash)
				s.Require().NoError(e)
				return res.Task.Status == 2 && res.Task.Score == 71
			},
			20*time.Second,
			5*time.Second,
		)
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTxTask(chainAAPIEndpoint, taskHash2)
				s.Require().NoError(e)
				return res.Task.Status == 2 && res.Task.Score == 0
			},
			20*time.Second,
			5*time.Second,
		)
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTxTask(chainAAPIEndpoint, taskHash3)
				s.Require().NoError(e)
				return res.Task.Status == 2 && res.Task.Score == 70
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("create_task", func() {
		bountyAmount, _ := sdk.NewIntFromString("500000")
		bounty := sdk.NewCoin(uctkDenom, bountyAmount)
		s.executeOracleCreateTask(s.chainA, 0, taskContract, taskFunction, bounty.String(), alice.String(), feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTask(chainAAPIEndpoint, taskContract, taskFunction)
				s.Require().NoError(e)
				return res.Task.Status == 1
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("respond_task", func() {
		s.executeOracleRespondTask(s.chainA, 0, 90, taskContract, taskFunction, alice.String(), feesAmountCoin.String())
		s.executeOracleRespondTask(s.chainA, 0, 90, taskContract, taskFunction, bob.String(), feesAmountCoin.String())
		s.executeOracleRespondTask(s.chainA, 0, 50, taskContract, taskFunction, charle.String(), feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleTask(chainAAPIEndpoint, taskContract, taskFunction)
				s.Require().NoError(e)
				return len(res.Task.Responses) == 3
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("claim_reward", func() {
		s.executeOracleClaimReward(s.chainA, 0, alice.String(), feesAmountCoin.String())
	})

	s.Run("remove_operator", func() {
		s.executeOracleRemoveOperator(s.chainA, 0, charle.String(), feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				res, e := queryOracleOperators(chainAAPIEndpoint)
				s.Require().NoError(e)
				return len(res.Operators) == 2
			},
			20*time.Second,
			5*time.Second,
		)
	})
}
