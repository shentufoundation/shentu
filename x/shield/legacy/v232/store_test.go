package v232_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/certikfoundation/shentu/v2/app"
	"github.com/certikfoundation/shentu/v2/x/shield/keeper"
	"github.com/certikfoundation/shentu/v2/x/shield/types"
)

//shared test
type MigrationTestSuite struct {
	suite.Suite

	app    *shentuapp.ShentuApp
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (suite *MigrationTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})
	suite.keeper = suite.app.ShieldKeeper
}

func (suite *MigrationTestSuite) TestMigrateParams() {
	tests := []struct {
		name          string
		description   string
		newParams     types.DistributionParams
		expectedError bool
	}{
		{
			name:        "Zero for all",
			description: "Set all parameter constants to zero",
			newParams:   types.NewDistributionParams(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		},
		{
			name:        "One for all",
			description: "Set all parameter constants to one",
			newParams:   types.NewDistributionParams(sdk.OneDec(), sdk.OneDec(), sdk.OneDec()),
		},
	}

	for _, tc := range tests {
		suite.keeper.SetDistributionParams(suite.ctx, tc.newParams)
		actual := suite.keeper.GetDistributionParams(suite.ctx)
		suite.Require().Equal(actual, tc.newParams)
	}
}

func TestPoolTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationTestSuite))
}
