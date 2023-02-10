package e2e

import (
	"bytes"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bountytypes "github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (s *IntegrationTestSuite) TestIBCTokenTransfer() {
	var ibcStakeDenom string

	s.Run("send_photon_to_chainB", func() {
		recipient := s.chainB.validators[0].keyInfo.GetAddress().String()
		token := sdk.NewInt64Coin(photonDenom, 3300000000) // 3,300photon
		s.sendIBC(s.chainA.id, s.chainB.id, recipient, token)

		chainBAPIEndpoint := s.valResources[s.chainB.id][0].GetHostPort("9090/tcp")

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
		chainAAPIEndpoint := s.valResources[s.chainA.id][0].GetHostPort("9090/tcp")
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
		chainAAPIEndpoint := s.valResources[s.chainA.id][0].GetHostPort("9090/tcp")
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

func (s *IntegrationTestSuite) TestGovernment() {

	chainAAPIEndpoint := s.valResources[s.chainA.id][0].GetHostPort("9090/tcp")
	validatorA := s.chainA.validators[0]
	validatorAAddr := validatorA.keyInfo.GetAddress()

	s.Run("submit_upgrade_proposal", func() {
		height, err := s.getLatestBlockHeight(chainAAPIEndpoint)
		s.Require().NoError(err)
		proposalHeight := int(height) + proposalBlockBuffer

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
	chainAAPIEndpoint := s.valResources[s.chainA.id][0].GetHostPort("9090/tcp")
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
				s.T().Logf("===> status: %#v\n", status)
				s.T().Logf("===> coll: %#v\n", status.TotalCollateral)
				s.T().Logf("===> amount: %#v\n", depositAmount)
				return true
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

func (s *IntegrationTestSuite) TestBounty() {
	chainAAPIEndpoint := s.valResources[s.chainA.id][0].GetHostPort("9090/tcp")
	validatorA := s.chainA.validators[0]
	accountA := s.chainA.accounts[0]
	accountAAddr := accountA.keyInfo.GetAddress()
	accountB := s.chainA.accounts[1]
	accountBAddr := accountB.keyInfo.GetAddress()

	bountyKeyFile := "e2e_bounty_key.json"
	generateBountyKeyFile(validatorA.configDir() + "/" + bountyKeyFile)
	bountyKeyPath := "/root/.shentud/" + bountyKeyFile

	s.Run("create_program", func() {
		bountyProgramCounter++
		s.T().Logf("Creating program %d on chain %s", bountyProgramCounter, s.chainA.id)
		var (
			programDesc    = "program-desc"
			commissionRate = "2"
			endTime        = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		)
		s.executeCreateProgram(s.chainA, 0, accountAAddr.String(), programDesc, bountyKeyPath, commissionRate, depositAmountCoin.String(), endTime, feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				rsp, err := queryBountyProgram(chainAAPIEndpoint, bountyProgramCounter)
				s.Require().NoError(err)
				return rsp.GetProgram().ProgramId == uint64(bountyProgramCounter)
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("submit_finding", func() {
		bountyFindingCounter++
		s.T().Logf("Submit finding %d on program %d chain %s", bountyFindingCounter, bountyProgramCounter, s.chainA.id)
		var (
			findingDesc  = "finding-desc"
			findingTitle = "finding-title"
			findingPoc   = "finding-poc"
		)
		s.executeSubmitFinding(s.chainA, 0, bountyProgramCounter, accountBAddr.String(), findingDesc, findingTitle, findingPoc, feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				rsp, err := queryBountyFinding(chainAAPIEndpoint, bountyFindingCounter)
				s.Require().NoError(err)
				return rsp.GetFinding().FindingStatus == 0
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("reject_finding", func() {
		s.T().Logf("Accept finding %d on program %d chain %s", bountyFindingCounter, bountyProgramCounter, s.chainA.id)
		var (
			findingComment = "reject-comment"
		)
		s.executeRejectFinding(s.chainA, 0, bountyFindingCounter, accountAAddr.String(), findingComment, feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				rsp, err := queryBountyFinding(chainAAPIEndpoint, bountyFindingCounter)
				s.Require().NoError(err)
				return rsp.GetFinding().FindingStatus == 2
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("accept_finding", func() {
		s.T().Logf("Accept finding %d on program %d chain %s", bountyFindingCounter, bountyProgramCounter, s.chainA.id)
		var (
			findingComment = "accept-comment"
		)
		s.executeAcceptFinding(s.chainA, 0, bountyFindingCounter, accountAAddr.String(), findingComment, feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				rsp, err := queryBountyFinding(chainAAPIEndpoint, bountyFindingCounter)
				s.Require().NoError(err)
				return rsp.GetFinding().FindingStatus == 1
			},
			20*time.Second,
			5*time.Second,
		)
	})

	s.Run("release_finding", func() {
		s.T().Logf("Release finding %d on program %d chain %s", bountyFindingCounter, bountyProgramCounter, s.chainA.id)
		s.executeReleaseFinding(s.chainA, 0, bountyFindingCounter, accountAAddr.String(), bountyKeyPath, feesAmountCoin.String())
		s.Require().Eventually(
			func() bool {
				rsp, err := queryBountyFinding(chainAAPIEndpoint, bountyFindingCounter)
				s.Require().NoError(err)
				var poc bountytypes.PlainTextPoc
				proto.Unmarshal(rsp.Finding.FindingPoc.Value, &poc)
				return string(poc.FindingPoc) == "finding-poc"
			},
			20*time.Second,
			5*time.Second,
		)
	})
}
