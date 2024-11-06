package e2e

import (
	"time"

	bountytypes "github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (s *IntegrationTestSuite) testBounty() {
	s.Run("test_bounty", func() {
		var (
			// err      error
			valIdx       = 0
			c            = s.chainA
			grpcEndpoint = s.valResources[s.chainA.id][0].GetHostPort("9090/tcp")
		)

		alice, _ := c.genesisAccounts[1].keyInfo.GetAddress()
		bob, _ := c.genesisAccounts[2].keyInfo.GetAddress()
		charlie, _ := c.genesisAccounts[3].keyInfo.GetAddress()

		programID := "fc28a970-f977-4bcb-bbfb-560baaaf7dd2"
		programName := "e2e-program-name"
		programDetail := `{"desc":"Refer to https://bounty.desc/cosmos for more details.","targets":["https://github.com/bounty/repo"],"total_bounty":500000,"bounty_denom":"USDT","bounty_levels":[{"severity":"critical","bounty":{"min_amount":"1","max_amount":"25000"}},{"severity":"high","bounty":{"min_amount":"1","max_amount":"3000"}},{"severity":"medium","bounty":{"min_amount":"1","max_amount":"1000"}},{"severity":"low","bounty":{"min_amount":"1","max_amount":"500"}},{"severity":"informational","bounty":{"min_amount":"1","max_amount":"1"}}]}`

		findingID := "4b34ff64-ad6a-4dda-98f5-6da02db7106c"
		findingDesc := "e2e-finding-desc"
		findingPoc := "e2e-finding-poc"
		// Create a program
		s.execCreateProgram(c, valIdx, programID, programName, programDetail, alice.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				program, err := queryProgram(grpcEndpoint, programID)
				return err == nil && program.ProgramId == programID && program.Status == bountytypes.ProgramStatusInactive
			},
			20*time.Second,
			5*time.Second,
		)

		// Create a duplicate program
		s.execCreateProgram(c, valIdx, programID, "dupe-name", programDetail, alice.String(), feesAmountCoin, true)

		// Active a program by non-admin
		s.execActivateProgram(c, valIdx, programID, alice.String(), feesAmountCoin, true)

		// Issue admin certificate
		certifierAcct, _ := c.certifier.keyInfo.GetAddress()
		s.execIssueCertificate(c, valIdx, charlie.String(), "bountyadmin", "set bounty admin", certifierAcct.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				ok, _ := queryCertificate(grpcEndpoint, charlie.String(), "bountyadmin")
				return ok
			},
			20*time.Second,
			5*time.Second,
		)

		// Submit finding to inactive program
		s.execSubmitFinding(c, valIdx, programID, findingID, "MEDIUM", findingDesc, findingPoc, bob.String(), feesAmountCoin, true)

		// Active a program by admin
		s.execActivateProgram(c, valIdx, programID, charlie.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				program, err := queryProgram(grpcEndpoint, programID)
				return err == nil && program.ProgramId == programID && program.Status == bountytypes.ProgramStatusActive
			},
			20*time.Second,
			5*time.Second,
		)

		// Submit a finding
		s.execSubmitFinding(c, valIdx, programID, findingID, "MEDIUM", findingDesc, findingPoc, bob.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				finding, err := queryFinding(grpcEndpoint, findingID)
				return err == nil && finding.FindingId == findingID && finding.Status == bountytypes.FindingStatusSubmitted
			},
			20*time.Second,
			5*time.Second,
		)

		// Edit a finding
		s.execEditFinding(c, valIdx, findingID, "LOW", findingDesc, findingPoc, bob.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				finding, err := queryFinding(grpcEndpoint, findingID)
				return err == nil && finding.FindingId == findingID && finding.SeverityLevel == bountytypes.Low
			},
			20*time.Second,
			5*time.Second,
		)

		// Edit a finding by non-creator
		s.execEditFinding(c, valIdx, findingID, "CRITICAL", findingDesc, findingPoc, alice.String(), feesAmountCoin, true)

		// Active a finding by non-admin
		s.execActivateFinding(c, valIdx, findingID, bob.String(), feesAmountCoin, true)

		// Active a finding by admin
		s.execActivateFinding(c, valIdx, findingID, charlie.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				finding, err := queryFinding(grpcEndpoint, findingID)
				return err == nil && finding.FindingId == findingID && finding.Status == bountytypes.FindingStatusActive
			},
			20*time.Second,
			5*time.Second,
		)

		// Close program by admin
		s.execCloseProgram(c, valIdx, programID, charlie.String(), feesAmountCoin, true)

		findingFingerprint, err := queryFindingFingerprint(grpcEndpoint, findingID)
		s.Require().NoError(err)
		// Confirm a finding by non-client
		s.execConfirmFinding(c, valIdx, findingID, findingFingerprint, bob.String(), feesAmountCoin, true)

		// Confirm a finding by client
		s.execConfirmFinding(c, valIdx, findingID, findingFingerprint, alice.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				finding, err := queryFinding(grpcEndpoint, findingID)
				return err == nil && finding.FindingId == findingID && finding.Status == bountytypes.FindingStatusConfirmed
			},
			20*time.Second,
			5*time.Second,
		)

		// Edit payment by creator
		s.execEditPayment(c, valIdx, findingID, "payment-hash", bob.String(), feesAmountCoin, true)

		// Edit payment by client
		s.execEditPayment(c, valIdx, findingID, "payment-hash", alice.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				finding, err := queryFinding(grpcEndpoint, findingID)
				return err == nil && finding.FindingId == findingID && finding.PaymentHash == "payment-hash"
			},
			20*time.Second,
			5*time.Second,
		)

		// Confirm paid by non-creator
		s.execConfirmPayment(c, valIdx, findingID, alice.String(), feesAmountCoin, true)

		// Confirm paid by creator
		s.execConfirmPayment(c, valIdx, findingID, bob.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				finding, err := queryFinding(grpcEndpoint, findingID)
				return err == nil && finding.FindingId == findingID && finding.Status == bountytypes.FindingStatusPaid
			},
			20*time.Second,
			5*time.Second,
		)

		// Publish a finding by creator
		s.execPublishFinding(c, valIdx, findingID, findingDesc, findingPoc, bob.String(), feesAmountCoin, true)

		// Publish a finding by client
		s.execPublishFinding(c, valIdx, findingID, findingDesc, findingPoc, alice.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				finding, err := queryFinding(grpcEndpoint, findingID)
				return err == nil && finding.FindingId == findingID && finding.ProofOfConcept == findingPoc
			},
			20*time.Second,
			5*time.Second,
		)

		// Close a program by non-client
		s.execCloseProgram(c, valIdx, programID, bob.String(), feesAmountCoin, true)

		// Close a program by client
		s.execCloseProgram(c, valIdx, programID, alice.String(), feesAmountCoin, false)
		s.Require().Eventually(
			func() bool {
				program, err := queryProgram(grpcEndpoint, programID)
				return err == nil && program.ProgramId == programID && program.Status == bountytypes.ProgramStatusClosed
			},
			20*time.Second,
			5*time.Second,
		)
	})
}
