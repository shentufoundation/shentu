package e2e

func (s *IntegrationTestSuite) TestIBC() {
	s.testIBCTokanTransfer()
}

func (s *IntegrationTestSuite) TestBank() {
	s.testBankTokenTransfer()
}

func (s *IntegrationTestSuite) TestDestribution() {
	s.testDistribution()
}

func (s *IntegrationTestSuite) TestFeegrant() {
	s.testFeeGrant()
}

func (s *IntegrationTestSuite) TestStaking() {
	s.testStaking()
}

func (s *IntegrationTestSuite) TestBounty() {
	s.testBounty()
}

func (s *IntegrationTestSuite) TestGov() {
	s.testCommonProposal()
}
