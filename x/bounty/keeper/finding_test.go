package keeper_test

import (
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestFindingList_GetSet() {
	findIDs := []uint64{10, 20, 30, 40}
	var pid uint64 = 2
	err := suite.keeper.SetPidFindingIDList(suite.ctx, pid, findIDs)
	suite.Require().NoError(err)

	findIDs2, err := suite.keeper.GetPidFindingIDList(suite.ctx, pid)
	suite.Require().NoError(err)
	suite.Require().Equal(findIDs, findIDs2)
}

func (suite *KeeperTestSuite) TestFindingWithdrawal() {
	var findingID uint64 = 1
	testAcct := suite.address[0]
	finding := types.Finding{
		FindingId:        findingID,
		SubmitterAddress: testAcct.String(),
		Active:           true,
	}
	suite.keeper.SetFinding(suite.ctx, finding)

	f2, err := suite.keeper.WithdrawalFinding(suite.ctx, testAcct, findingID)
	suite.Require().NoError(err)
	suite.Require().False(f2.Active)

	_, err = suite.keeper.WithdrawalFinding(suite.ctx, testAcct, findingID)
	suite.Require().Error(err)

	f4, err := suite.keeper.ReactivateFinding(suite.ctx, testAcct, findingID)
	suite.Require().NoError(err)
	suite.Require().True(f4.Active)

	_, err = suite.keeper.ReactivateFinding(suite.ctx, testAcct, findingID)
	suite.Require().Error(err)
}
