package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types1 "github.com/cosmos/cosmos-sdk/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/bounty/keeper"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

// shared setup
type KeeperTestSuite struct {
	suite.Suite
	app     *shentuapp.ShentuApp
	ctx     sdk.Context
	keeper  keeper.Keeper
	address []sdk.AccAddress
	// queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = shentuapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.BountyKeeper

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := sdksimapp.FundAccount(
			suite.app.BankKeeper,
			suite.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin("uctk", sdk.NewInt(10000000000)), // 1,000 CTK
			),
		)
		if err != nil {
			panic(err)
		}
	}

	suite.address = []sdk.AccAddress{acc1, acc2, acc3, acc4}
	//suite.keeper.SetCertifier(suite.ctx, types.NewCertifier(suite.address[0], "", suite.address[0], ""))
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

	deposit1 := types1.NewInt(10000)

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{"Program(1)  -> Set: Simple",
			args{
				program: []types.Program{
					{
						ProgramId:         1,
						CreatorAddress:    suite.address[0].String(),
						SubmissionEndTime: time.Now(),
						Description:       "for test1",
						Deposit: []types1.Coin{
							{
								Denom:  "uctk",
								Amount: deposit1,
							},
						},
						CommissionRate: types1.NewDec(1),
					},
				},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	suite.keeper.SetNextProgramID(suite.ctx, 1)

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			for _, program := range tc.args.program {
				fmt.Println("----", program)
				nextID := suite.keeper.GetNextProgramID(suite.ctx)
				fmt.Println(nextID)

				suite.keeper.SetProgram(suite.ctx, program)
			}
		})
	}
}
