package keeper_test

import (
	"fmt"
)

func (suite *KeeperTestSuite) TestHash() {
	hash := suite.keeper.GetProofHash(1, "shentu16gzt5vd0dd5c98ajl3ld2ltvcahxgyygd58n3m", "test")
	fmt.Println(hash)
}
