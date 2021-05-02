package keeper_test

// import (
// 	gocontext "context"
// 	"strings"

// 	"github.com/certikfoundation/shentu/x/oracle/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func (suite *KeeperTestSuite) TestQueryOperator() {
// 	queryClient := suite.queryClient
// 	type args struct {
// 		collateral   int64
// 		senderAddr   sdk.AccAddress
// 		proposerAddr sdk.AccAddress
// 		operatorName string
// 		operatorAddr sdk.AccAddress
// 	}

// 	type errArgs struct {
// 		shouldPass bool
// 		contains   string
// 	}

// 	tests := []struct {
// 		name    string
// 		args    args
// 		errArgs errArgs
// 	}{
// 		{"Operator(1) Query: Empty Operator Address",
// 			args{
// 				collateral:   50000,
// 				senderAddr:   suite.address[0],
// 				proposerAddr: suite.address[1],
// 				operatorName: "Operator",
// 			},
// 			errArgs{
// 				shouldPass: false,
// 				contains:   "",
// 			},
// 		},
// 		{"Operator(1) Query: Non-Existent Operator",
// 			args{
// 				collateral:   50000,
// 				senderAddr:   suite.address[0],
// 				proposerAddr: suite.address[1],
// 				operatorName: "Operator",
// 				operatorAddr: suite.address[3],
// 			},
// 			errArgs{
// 				shouldPass: false,
// 				contains:   "",
// 			},
// 		},
// 		{"Operator(1) Query: Valid Request",
// 			args{
// 				collateral:   50000,
// 				senderAddr:   suite.address[0],
// 				proposerAddr: suite.address[1],
// 				operatorName: "Operator",
// 				operatorAddr: suite.address[0],
// 			},
// 			errArgs{
// 				shouldPass: true,
// 				contains:   "",
// 			},
// 		},
// 	}
// 	for _, tc := range tests {
// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()
// 			// create an operator
// 			err := suite.keeper.CreateOperator(suite.ctx, tc.args.senderAddr, sdk.Coins{sdk.NewInt64Coin("uctk", tc.args.collateral)}, tc.args.proposerAddr, tc.args.operatorName)
// 			suite.Require().NoError(err, tc.name)
// 			// what we got from query
// 			got, err := queryClient.Operator(gocontext.Background(), &types.QueryOperatorRequest{Address: tc.args.operatorAddr.String()})
// 			suite.Require().NoError(err, tc.name)
// 			// what we want
// 			want, err := suite.keeper.GetOperator(suite.ctx, tc.args.senderAddr)
// 			suite.Require().NoError(err, tc.name)
// 			if tc.errArgs.shouldPass {
// 				suite.Require().NoError(err, tc.name)
// 				suite.Equal(got, want)
// 			} else {
// 				suite.Require().Error(err, tc.name)
// 				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
// 			}
// 		})
// 	}
// }
