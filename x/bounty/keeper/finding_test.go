package keeper_test

func (suite *KeeperTestSuite) TestFindingList_GetSet() {
	findIDs := []uint64{10, 20, 30, 40}
	var pid uint64
	pid = 2
	err := suite.keeper.SetPidFindingIDList(suite.ctx, pid, findIDs)
	suite.Require().NoError(err)

	findIDs2, err := suite.keeper.GetPidFindingIDList(suite.ctx, pid)
	suite.Require().NoError(err)
	suite.Require().Equal(findIDs, findIDs2)
}
