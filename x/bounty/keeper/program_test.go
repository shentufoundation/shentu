package keeper_test

import "github.com/shentufoundation/shentu/v2/x/bounty/types"

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

				storedPrograms := suite.keeper.GetAllPrograms(suite.ctx)
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
		Description:  "desc",
		AdminAddress: suite.address[0].String(),
		Status:       types.ProgramStatusInactive,
	}
	suite.keeper.SetProgram(suite.ctx, program)
	storedProgram, isExist := suite.keeper.GetProgram(suite.ctx, program.ProgramId)
	suite.Require().Equal(true, isExist)
	suite.Require().Equal(program.ProgramId, storedProgram.ProgramId)

	isCert := suite.app.CertKeeper.IsCertifier(suite.ctx, suite.address[3])
	suite.Require().True(isCert)

	// normal addr open program
	err := suite.keeper.OpenProgram(suite.ctx, program.ProgramId, suite.address[1])
	suite.Require().Error(err)
	// certifier open program
	err = suite.keeper.OpenProgram(suite.ctx, program.ProgramId, suite.address[3])
	suite.Require().NoError(err)

	// normal addr close program
	err = suite.keeper.CloseProgram(suite.ctx, program.ProgramId, suite.address[1])
	suite.Require().Error(err)
	// admin close program
	err = suite.keeper.CloseProgram(suite.ctx, program.ProgramId, suite.address[3])
	suite.Require().NoError(err)
}
