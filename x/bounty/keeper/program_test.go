package keeper_test

import (
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

func (suite *KeeperTestSuite) TestSetGetProgram() {
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
						Detail:       "detail",
						AdminAddress: suite.programAddr.String(),
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
				err := suite.keeper.Programs.Set(suite.ctx, program.ProgramId, program)
				suite.Require().NoError(err)

				storedProgram, err := suite.keeper.Programs.Get(suite.ctx, program.ProgramId)
				suite.Require().NoError(err)

				var storedPrograms []types.Program
				err = suite.keeper.Programs.Walk(suite.ctx, nil, func(_ string, p types.Program) (bool, error) {
					storedPrograms = append(storedPrograms, p)
					return false, nil
				})
				suite.Require().NoError(err)
				suite.Require().Equal(1, len(storedPrograms))

				if tc.errArgs.shouldPass {
					suite.Require().Equal(program.ProgramId, storedProgram.ProgramId)
				} else {
					suite.Require().NotEqual(program.ProgramId, storedProgram.ProgramId)
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOpenCloseProgram() {
	type errArgs struct {
		shouldPass bool
		contains   string
	}

	// create program
	program := types.Program{
		ProgramId:    "1",
		Name:         "name",
		Detail:       "detail",
		AdminAddress: suite.programAddr.String(),
		Status:       types.ProgramStatusInactive,
	}
	err := suite.keeper.Programs.Set(suite.ctx, program.ProgramId, program)
	suite.Require().NoError(err)

	storedProgram, err := suite.keeper.Programs.Get(suite.ctx, program.ProgramId)
	suite.Require().NoError(err)
	suite.Require().Equal(program.ProgramId, storedProgram.ProgramId)

	isCert := suite.app.CertKeeper.IsBountyAdmin(suite.ctx, suite.bountyAdminAddr)
	suite.Require().True(isCert)

	// normal addr open program
	err = suite.keeper.ActivateProgram(suite.ctx, program.ProgramId, suite.normalAddr)
	suite.Require().Error(err)
	// certifier open program
	err = suite.keeper.ActivateProgram(suite.ctx, program.ProgramId, suite.bountyAdminAddr)
	suite.Require().NoError(err)

	// normal addr close program
	err = suite.keeper.CloseProgram(suite.ctx, program.ProgramId, suite.normalAddr)
	suite.Require().Error(err)
	// admin close program
	err = suite.keeper.CloseProgram(suite.ctx, program.ProgramId, suite.programAddr)
	suite.Require().NoError(err)
}
