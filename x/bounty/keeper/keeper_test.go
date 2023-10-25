package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite

	app         *shentuapp.ShentuApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	address     []sdk.AccAddress
	msgServer   types.MsgServer
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.BountyKeeper
	suite.address = shentuapp.AddTestAddrs(suite.app, suite.ctx, 4, sdk.NewInt(1e10))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.app.BountyKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
	suite.msgServer = keeper.NewMsgServerImpl(suite.keeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestProgram_GetSet() {
	type args struct {
		program []types.Program
	}

	type errArgs struct {
		shouldPass bool
		contains   string
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Program(1)  -> Set: Simple",
			args{
				program: []types.Program{
					{
						ProgramId:    "1",
						Name:         "name",
						Description:  "desc",
						AdminAddress: suite.address[0].String(),
						Status:       1,
					},
				},
			},
			errArgs{
				shouldPass: true,
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			for _, program := range tc.args.program {
				suite.keeper.SetProgram(suite.ctx, program)
				storedProgram, isExist := suite.keeper.GetProgram(suite.ctx, program.ProgramId)
				suite.Require().Equal(true, isExist)

				if tc.errArgs.shouldPass {
					suite.Require().Equal(program.ProgramId, storedProgram.ProgramId)
				} else {
					suite.Require().NotEqual(program.ProgramId, storedProgram.ProgramId)
				}
			}
		})
	}
}
