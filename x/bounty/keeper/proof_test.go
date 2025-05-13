package keeper_test

import (
	"fmt"
)

func (suite *KeeperTestSuite) TestHash() {
	hash := suite.keeper.GetProofHash(1, "shentu14ayuhu60zyc7a5chxy65s5g2cfamvufwm7vd52", "test")
	fmt.Println(hash)
}
