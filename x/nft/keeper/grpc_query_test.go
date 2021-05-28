package keeper_test

import (
	gocontext "context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/nft/types"
)

func (suite *KeeperTestSuite) TestQueryAdmin() {
	type args struct {
		adminAddr   sdk.AccAddress
		requestAddr string
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
		{"Admin(1) Query: Empty Address",
			args{
				adminAddr: suite.address[0],
			},
			errArgs{
				shouldPass: false,
				contains:   "empty address",
			},
		},
		{"Admin(1) Query: Admin Address",
			args{
				adminAddr:   suite.address[0],
				requestAddr: suite.address[0].String(),
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Admin(1) Query: Non-admin Address",
			args{
				adminAddr:   suite.address[0],
				requestAddr: suite.address[1].String(),
			},
			errArgs{
				shouldPass: false,
				contains:   "not found",
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			suite.keeper.SetAdmin(suite.ctx, tc.args.adminAddr)
			res, err := suite.queryClient.Admin(gocontext.Background(), &types.QueryAdminRequest{Address: tc.args.requestAddr})
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotNil(res.Admin)
				suite.Require().Equal(res.Admin.Address, tc.args.adminAddr.String())
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryAdmins() {
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
		{"Admins(1) Query: No Admins",
			args{
				addrs: []sdk.AccAddress{},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Admins(2) Query: Two Admins",
			args{
				addrs: []sdk.AccAddress{suite.address[0], suite.address[1]},
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
			for _, addr := range tc.args.addrs {
				suite.keeper.SetAdmin(suite.ctx, addr)
			}
			res, err := suite.queryClient.Admins(gocontext.Background(), &types.QueryAdminsRequest{})
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotNil(res)
				suite.Require().Len(res.Admins, len(tc.args.addrs))
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}
