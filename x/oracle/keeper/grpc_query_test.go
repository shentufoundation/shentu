package keeper_test

import (
	"encoding/base64"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/oracle/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryOperator() {
	tests := []struct {
		name    string
		req     *types.QueryOperatorRequest
		preRun  func()
		postRun func(_ *types.QueryOperatorResponse)
		expPass bool
	}{
		{
			"invalid address",
			&types.QueryOperatorRequest{
				Address: "invalid",
			},
			func() {},
			func(*types.QueryOperatorResponse) {},
			false,
		},
		{
			"operator does not exist",
			&types.QueryOperatorRequest{
				Address: suite.address[2].String(),
			},
			func() {},
			func(*types.QueryOperatorResponse) {},
			false,
		},
		{
			"valid request",
			&types.QueryOperatorRequest{
				Address: suite.address[0].String(),
			},
			func() {
				suite.createOperator(suite.address[0], suite.address[1])
			},
			func(res *types.QueryOperatorResponse) {
				suite.Require().Equal(suite.address[0].String(), res.Operator.Address)
				suite.Require().Equal(suite.address[1].String(), res.Operator.Proposer)
			},
			true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.preRun()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.Operator(ctx, tc.req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
			tc.postRun(res)
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryOperators() {
	tests := []struct {
		name    string
		req     *types.QueryOperatorsRequest
		preRun  func()
		postRun func(_ *types.QueryOperatorsResponse)
		expPass bool
	}{
		{
			"valid request",
			&types.QueryOperatorsRequest{},
			func() {
				suite.createOperator(suite.address[0], suite.address[1])
			},
			func(res *types.QueryOperatorsResponse) {
				suite.Require().Len(res.Operators, 1)
				suite.Require().Equal(suite.address[0].String(), res.Operators[0].Address)
				suite.Require().Equal(suite.address[1].String(), res.Operators[0].Proposer)
			},
			true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.preRun()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.Operators(ctx, tc.req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
			tc.postRun(res)
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryWithdraws() {
	tests := []struct {
		name    string
		req     *types.QueryWithdrawsRequest
		preRun  func()
		postRun func(_ *types.QueryWithdrawsResponse)
		expPass bool
	}{
		{
			"valid request",
			&types.QueryWithdrawsRequest{},
			func() {
				suite.createWithdraw(suite.address[0])
			},
			func(res *types.QueryWithdrawsResponse) {
				suite.Require().Len(res.Withdraws, 1)
				suite.Require().Equal(suite.address[0].String(), res.Withdraws[0].Address)
				suite.Require().Equal(suite.keeper.GetLockedPoolParams(suite.ctx).LockedInBlocks, res.Withdraws[0].DueBlock)
			},
			true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.preRun()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.Withdraws(ctx, tc.req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
			tc.postRun(res)
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryTask() {
	tests := []struct {
		name    string
		req     *types.QueryTaskRequest
		preRun  func()
		postRun func(_ *types.QueryTaskResponse)
		expPass bool
	}{
		{
			"empty contract",
			&types.QueryTaskRequest{
				Contract: "",
				Function: "func",
			},
			func() {},
			func(*types.QueryTaskResponse) {},
			false,
		},
		{
			"empty function",
			&types.QueryTaskRequest{
				Contract: "0x1234567890abcdef",
				Function: "",
			},
			func() {},
			func(*types.QueryTaskResponse) {},
			false,
		},
		{
			"valid request",
			&types.QueryTaskRequest{
				Contract: "0x1234567890abcdef",
				Function: "func",
			},
			func() {
				suite.createTask("0x1234567890abcdef", "func", suite.address[0])
			},
			func(res *types.QueryTaskResponse) {
				suite.Require().Equal("0x1234567890abcdef", res.Task.Contract)
				suite.Require().Equal("func", res.Task.Function)
			},
			true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.preRun()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.Task(ctx, tc.req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
			tc.postRun(res)
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryTxTask() {
	tests := []struct {
		name    string
		req     *types.QueryTxTaskRequest
		preRun  func()
		postRun func(_ *types.QueryTxTaskResponse)
		expPass bool
	}{
		{
			"empty tx hash",
			&types.QueryTxTaskRequest{
				TxHash: "",
			},
			func() {},
			func(*types.QueryTxTaskResponse) {},
			false,
		},
		{
			"invalid tx hash",
			&types.QueryTxTaskRequest{
				TxHash: "1234567",
			},
			func() {},
			func(*types.QueryTxTaskResponse) {},
			false,
		},
		{
			"valid request",
			&types.QueryTxTaskRequest{
				TxHash: base64.StdEncoding.EncodeToString([]byte("valid request")),
			},
			func() {
				suite.createTxTask([]byte("valid request"), suite.address[0])
			},
			func(res *types.QueryTxTaskResponse) {
				suite.Require().Equal([]byte("valid request"), res.Task.TxHash)
			},
			true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.preRun()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.TxTask(ctx, tc.req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
			tc.postRun(res)
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryResponse() {
	tests := []struct {
		name    string
		req     *types.QueryResponseRequest
		preRun  func()
		postRun func(_ *types.QueryResponseResponse)
		expPass bool
	}{
		{
			"no task found",
			&types.QueryResponseRequest{
				Contract:        "",
				Function:        "func",
				OperatorAddress: suite.address[0].String(),
			},
			func() {},
			func(*types.QueryResponseResponse) {},
			false,
		},
		{
			"no operator found",
			&types.QueryResponseRequest{
				Contract:        "0x1234567890abcdef",
				Function:        "func",
				OperatorAddress: suite.address[1].String(),
			},
			func() {
				suite.createOperator(suite.address[0], suite.address[0])
				suite.createTask("0x1234567890abcdef", "func", suite.address[0])
				suite.respondToTask("0x1234567890abcdef", "func", suite.address[0])
			},
			func(*types.QueryResponseResponse) {},
			false,
		},
		{
			"valid request",
			&types.QueryResponseRequest{
				Contract:        "0x1234567890abcdef",
				Function:        "func",
				OperatorAddress: suite.address[0].String(),
			},
			func() {
				suite.createOperator(suite.address[0], suite.address[0])
				suite.createTask("0x1234567890abcdef", "func", suite.address[0])
				suite.respondToTask("0x1234567890abcdef", "func", suite.address[0])
			},
			func(res *types.QueryResponseResponse) {
				suite.Require().Equal(suite.address[0].String(), res.Response.Operator)
			},
			true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.preRun()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.Response(ctx, tc.req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
			tc.postRun(res)
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryTxResponse() {
	tests := []struct {
		name    string
		req     *types.QueryTxResponseRequest
		preRun  func()
		postRun func(_ *types.QueryTxResponseResponse)
		expPass bool
	}{
		{
			"no task found",
			&types.QueryTxResponseRequest{
				TxHash:          "",
				OperatorAddress: suite.address[0].String(),
			},
			func() {},
			func(response *types.QueryTxResponseResponse) {},
			false,
		},
		{
			"no operator found",
			&types.QueryTxResponseRequest{
				TxHash:          base64.StdEncoding.EncodeToString([]byte("no operator")),
				OperatorAddress: suite.address[1].String(),
			},
			func() {
				suite.createOperator(suite.address[0], suite.address[0])
				suite.createTxTask([]byte("no operator"), suite.address[0])
				suite.respondToTxTask([]byte("no operator"), suite.address[0])
			},
			func(*types.QueryTxResponseResponse) {},
			false,
		},
		{
			"valid request",
			&types.QueryTxResponseRequest{
				TxHash:          base64.StdEncoding.EncodeToString([]byte("0x1234567890abcdef")),
				OperatorAddress: suite.address[0].String(),
			},
			func() {
				suite.createOperator(suite.address[0], suite.address[0])
				suite.createTxTask([]byte("0x1234567890abcdef"), suite.address[0])
				suite.respondToTxTask([]byte("0x1234567890abcdef"), suite.address[0])
			},
			func(res *types.QueryTxResponseResponse) {
				suite.Require().Equal(suite.address[0].String(), res.Response.Operator)
			},
			true,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			tc.preRun()
			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.TxResponse(ctx, tc.req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
			tc.postRun(res)
		})
	}
}

func (suite *KeeperTestSuite) createOperator(address, proposer sdk.AccAddress) {
	err := suite.keeper.CreateOperator(suite.ctx, address, sdk.Coins{sdk.NewInt64Coin("uctk", 50000)}, proposer, "operator")
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) createWithdraw(address sdk.AccAddress) {
	err := suite.keeper.CreateWithdraw(suite.ctx, address, sdk.Coins{sdk.NewInt64Coin("uctk", 1000)})
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) createTask(contract, function string, creator sdk.AccAddress) {
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	expiration := time.Now().Add(time.Hour).UTC()
	waitingBlocks := int64(5)
	scTask := types.NewTask(
		contract, function, suite.ctx.BlockHeight(),
		bounty, "task", expiration,
		creator, suite.ctx.BlockHeight()+waitingBlocks, waitingBlocks)
	err := suite.keeper.CreateTask(suite.ctx, creator, &scTask)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) createTxTask(txHash []byte, creator sdk.AccAddress) {
	bounty := sdk.Coins{sdk.NewInt64Coin("uctk", 100000)}
	expiration := time.Now().Add(time.Hour).UTC()

	task := types.NewTxTask(
		txHash,
		creator.String(),
		bounty,
		expiration,
		types.TaskStatusPending,
	)
	err := suite.keeper.CreateTask(suite.ctx, creator, task)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) respondToTask(contract, function string, operatorAddress sdk.AccAddress) {
	err := suite.keeper.RespondToTask(suite.ctx, types.NewTaskID(contract, function), 100, operatorAddress)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) respondToTxTask(txHash []byte, operatorAddress sdk.AccAddress) {
	err := suite.keeper.RespondToTask(suite.ctx, txHash, 100, operatorAddress)
	suite.Require().NoError(err)
}
