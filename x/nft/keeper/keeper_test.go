package keeper_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/simapp"
	certtypes "github.com/certikfoundation/shentu/x/cert/types"
	"github.com/certikfoundation/shentu/x/nft/keeper"
	"github.com/certikfoundation/shentu/x/nft/types"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	certifier  = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	certifier2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tokenID  = "tokenid"
	tokenID2 = "tokenid2"
	tokenID3 = "tokenid3"
	tokenID4 = "tokenid4"
	tokenID5 = "tokenid5"

	tokenNm  = "tokennm"
	tokenNm2 = "tokennm2"
	tokenNm3 = "tokennm3"
	tokenNm4 = "tokennm4"
	tokenNm5 = "tokennm5"

	tokenURI  = "https://google.com/token-1.json"
	tokenURI2 = "https://google.com/token-2.json"
	tokenURI3 = "https://google.com/token-3.json"
	tokenURI4 = "https://google.com/token-4.json"
	tokenURI5 = "https://google.com/token-5.json"

	content = "content"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.app = app
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.NFTKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.NFTKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	for _, acc := range []sdk.AccAddress{acc1, acc2} {
		err := app.BankKeeper.AddCoins(
			suite.ctx,
			acc,
			sdk.NewCoins(
				sdk.NewCoin("uctk", sdk.NewInt(10000000000)),
			),
		)
		if err != nil {
			panic(err)
		}
	}

	suite.app.CertKeeper.SetCertifier(suite.ctx, certtypes.NewCertifier(certifier, "", certifier, ""))
	suite.app.CertKeeper.SetCertifier(suite.ctx, certtypes.NewCertifier(certifier2, "", certifier2, ""))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestAdmin_GetSet() {
	type args struct {
		addrs []sdk.AccAddress
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
		{
			name: "NFT(2) Get: One & All",
			args: args{
				addrs: []sdk.AccAddress{acc1, acc2},
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.keeper.SetAdmin(suite.ctx, tc.args.addrs[0])
			suite.keeper.SetAdmin(suite.ctx, tc.args.addrs[1])
			admin1, err := suite.keeper.GetAdmin(suite.ctx, tc.args.addrs[0])
			allAdmins := suite.keeper.GetAdmins(suite.ctx)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Equal(tc.args.addrs[0].String(), admin1.Address)
				suite.Len(allAdmins, 2)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAdmin_Delete() {
	type args struct {
		addrs       []sdk.AccAddress
		deletedAddr sdk.AccAddress
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
		{
			name: "NFT(1) Delete: Simple",
			args: args{
				addrs:       []sdk.AccAddress{acc1},
				deletedAddr: acc1,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "NFT(1) Delete: Add Two, Delete One",
			args: args{
				addrs:       []sdk.AccAddress{acc1, acc2},
				deletedAddr: acc2,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			for _, addr := range tc.args.addrs {
				suite.keeper.SetAdmin(suite.ctx, addr)
			}
			_, err := suite.keeper.GetAdmin(suite.ctx, tc.args.deletedAddr)
			suite.Require().NoError(err, tc.name)
			suite.keeper.DeleteAdmin(suite.ctx, tc.args.deletedAddr)
			allAdmins := suite.keeper.GetAdmins(suite.ctx)
			if tc.errArgs.shouldPass {
				suite.Len(allAdmins, len(tc.args.addrs)-1)
				for _, addr := range tc.args.addrs {
					admin, err := suite.keeper.GetAdmin(suite.ctx, addr)
					if addr.String() == tc.args.deletedAddr.String() {
						suite.Require().Error(err, tc.name)
						suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
					} else {
						suite.Require().NoError(err, tc.name)
						suite.Require().Contains(allAdmins, admin)
					}
				}
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestAdmin_Check() {
	type args struct {
		addrs []sdk.AccAddress
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
		{
			name: "NFT(1) Check: Simple",
			args: args{
				addrs: []sdk.AccAddress{acc1},
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.keeper.SetAdmin(suite.ctx, tc.args.addrs[0])
			_, err := suite.keeper.GetAdmin(suite.ctx, tc.args.addrs[0])
			suite.Require().NoError(err, tc.name)
			isAdmin := suite.keeper.CheckAdmin(suite.ctx, tc.args.addrs[0].String())
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.True(isAdmin)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}
