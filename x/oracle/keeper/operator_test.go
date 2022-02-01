package keeper_test

import (
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestOperator_Create() {
	type args struct {
		collateral   int64
		senderAddr   sdk.AccAddress
		proposerAddr sdk.AccAddress
		operatorName string
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
		{"Operator(1) Create: min collateral",
			args{
				collateral:   50000,
				senderAddr:   suite.address[0],
				proposerAddr: suite.address[1],
				operatorName: "Operator",
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Operator(1) Create: under min collateral",
			args{
				collateral:   10000,
				senderAddr:   suite.address[0],
				proposerAddr: suite.address[1],
				operatorName: "Operator",
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, tc.args.proposerAddr, tc.args.operatorName)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOperator_Get() {
	type args struct {
		collateral    int64
		senderAddr1   sdk.AccAddress
		proposerAddr1 sdk.AccAddress
		operatorName1 string
		senderAddr2   sdk.AccAddress
		proposerAddr2 sdk.AccAddress
		operatorName2 string
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
		{"Operator(2) Get: One & All",
			args{
				collateral:    50000,
				senderAddr1:   suite.address[0],
				proposerAddr1: suite.address[1],
				operatorName1: "Operator1",
				senderAddr2:   suite.address[2],
				proposerAddr2: suite.address[3],
				operatorName2: "Operator2",
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr1, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, tc.args.proposerAddr1, tc.args.operatorName1)
			suite.Require().NoError(err, tc.name)
			err = suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr2, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, tc.args.proposerAddr2, tc.args.operatorName2)
			suite.Require().NoError(err, tc.name)
			operator1, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr1)
			allOperators := suite.keeper.GetAllOperators(suite.ctx)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Equal(tc.args.senderAddr1.String(), operator1.Address)
				suite.Equal(sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, operator1.Collateral)
				suite.Equal(tc.args.proposerAddr1.String(), operator1.Proposer)
				suite.Len(allOperators, 2)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOperator_Remove() {
	type args struct {
		collateral    int64
		senderAddr1   sdk.AccAddress
		proposerAddr1 sdk.AccAddress
		operatorName1 string
		senderAddr2   sdk.AccAddress
		proposerAddr2 sdk.AccAddress
		operatorName2 string
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
		{"Operator(2) Remove: One",
			args{
				collateral:    50000,
				senderAddr1:   suite.address[0],
				proposerAddr1: suite.address[1],
				operatorName1: "Operator1",
				senderAddr2:   suite.address[2],
				proposerAddr2: suite.address[3],
				operatorName2: "Operator2",
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr1, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, tc.args.proposerAddr1, tc.args.operatorName1)
			suite.Require().NoError(err, tc.name)
			err = suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr2, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, tc.args.proposerAddr2, tc.args.operatorName2)
			suite.Require().NoError(err, tc.name)
			operator1, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr1)
			suite.Require().NoError(err, tc.name)
			// convert operator1.Address (string) back to sdk.AccAddress
			operator1Addr, _ := sdk.AccAddressFromBech32(operator1.Address)
			operator2, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr2)
			suite.Require().NoError(err, tc.name)
			// remove operator1
			err = suite.keeper.RemoveOperator(suite.ctx, operator1Addr.String(), operator1Addr.String())
			allOperators := suite.keeper.GetAllOperators(suite.ctx)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Len(allOperators, 1)
				suite.Equal(operator2, allOperators[0])
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOperator_Collateral() {
	type args struct {
		collateral         int64
		senderAddr         sdk.AccAddress
		proposerAddr       sdk.AccAddress
		operatorName       string
		collateralToAdd    int64
		collateralToReduce int64
		reduce             bool
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
		{"Operator(1) Add: 100,000 uctk",
			args{
				collateral:      50000,
				senderAddr:      suite.address[0],
				proposerAddr:    suite.address[1],
				operatorName:    "Operator1",
				collateralToAdd: 100000,
				reduce:          false,
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Operator(1) Add: 10,000 uctk -> Reduce: 5,000 uctk",
			args{
				collateral:         50000,
				senderAddr:         suite.address[2],
				proposerAddr:       suite.address[3],
				operatorName:       "Operator2",
				collateralToAdd:    10000,
				collateralToReduce: 5000,
				reduce:             true,
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Operator(1) Add: 10,000 uctk -> Reduce: 15,000 uctk", // 10,000+10,000-15,000=5000 < min(10,000)
			args{
				collateral:         50000,
				senderAddr:         suite.address[0],
				proposerAddr:       suite.address[1],
				operatorName:       "Operato2",
				collateralToAdd:    10000,
				collateralToReduce: 15000,
				reduce:             true,
			},
			errArgs{
				shouldPass: false,
				contains:   "collateral not enough",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, tc.args.proposerAddr, tc.args.operatorName)
			suite.Require().NoError(err, tc.name)
			collateral, err := suite.keeper.GetCollateralAmount(suite.ctx, tc.args.senderAddr)
			suite.Require().NoError(err, tc.name)
			suite.Equal(sdk.NewInt(tc.args.collateral), collateral)
			// add collateral
			suite.NoError(suite.keeper.AddCollateral(suite.ctx, tc.args.senderAddr, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateralToAdd)}))
			// check collateral amount
			collateral, err = suite.keeper.GetCollateralAmount(suite.ctx, tc.args.senderAddr)
			if tc.errArgs.shouldPass {
				if tc.args.reduce {
					// reduce collateral
					suite.NoError(suite.keeper.ReduceCollateral(suite.ctx, tc.args.senderAddr, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateralToReduce)}))
					collateral, err = suite.keeper.GetCollateralAmount(suite.ctx, tc.args.senderAddr)
					suite.Require().NoError(err, tc.name)
					suite.Equal(sdk.NewInt(tc.args.collateral+tc.args.collateralToAdd-tc.args.collateralToReduce), collateral)
				} else {
					suite.Require().NoError(err, tc.name)
					suite.Equal(sdk.NewInt(tc.args.collateral+tc.args.collateralToAdd), collateral)
				}
			} else {
				if tc.args.reduce {
					// reduce collateral
					err := suite.keeper.ReduceCollateral(suite.ctx, tc.args.senderAddr, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateralToReduce)})
					suite.Require().Error(err, tc.name)
					suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
				} else {
					suite.Require().Error(err, tc.name)
					suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOperator_Reward() {
	type args struct {
		collateral   int64
		senderAddr   sdk.AccAddress
		proposerAddr sdk.AccAddress
		operatorName string
		rewardToAdd  int64
		withdrawAll  bool
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
		{"Operator(1) Add: 100000 uctk",
			args{
				collateral:   50000,
				senderAddr:   suite.address[0],
				proposerAddr: suite.address[1],
				operatorName: "Operator",
				rewardToAdd:  100000,
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Operator(1) WithdrawAll: Reward <= Collateral",
			args{
				collateral:   50000,
				senderAddr:   suite.address[0],
				proposerAddr: suite.address[1],
				operatorName: "Operator",
				rewardToAdd:  50000,
				withdrawAll:  true,
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Operator(1) WithdrawAll: Reward > Collateral",
			args{
				collateral:   50000,
				senderAddr:   suite.address[0],
				proposerAddr: suite.address[1],
				operatorName: "Operator",
				rewardToAdd:  100000,
				withdrawAll:  true,
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, tc.args.proposerAddr, tc.args.operatorName)
			suite.Require().NoError(err, tc.name)
			err = suite.keeper.CreateTask(suite.ctx, "contract", "function", sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.rewardToAdd)}, "description", time.Now().Add(time.Hour).UTC(), tc.args.proposerAddr, int64(50))
			suite.Require().NoError(err, tc.name)
			err = suite.keeper.AddReward(suite.ctx, tc.args.senderAddr, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.rewardToAdd)})
			suite.Require().NoError(err, tc.name)
			if tc.errArgs.shouldPass {
				if tc.args.withdrawAll {
					withdrawAmt, err := suite.keeper.WithdrawAllReward(suite.ctx, tc.args.senderAddr)
					suite.Require().NoError(err, tc.name)

					suite.Equal(sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.rewardToAdd)}, withdrawAmt)

					operator, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr)
					suite.Require().NoError(err, tc.name)
					suite.Nil(operator.AccumulatedRewards)
				} else {
					operator, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr)
					suite.Require().NoError(err, tc.name)
					suite.Equal(sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.rewardToAdd)}, operator.AccumulatedRewards)
				}
			} else {
				if tc.args.withdrawAll {
					_, err := suite.keeper.WithdrawAllReward(suite.ctx, tc.args.senderAddr)
					suite.Require().Error(err, tc.name)
					suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
				} else {
					suite.Require().Error(err, tc.name)
					suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}
		})
	}
}
