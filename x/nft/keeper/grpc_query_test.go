package keeper_test

import (
	gocontext "context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

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
		{
			name: "Admin(1) Query: Empty Address",
			args: args{
				adminAddr: acc1,
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "empty address",
			},
		},
		{
			name: "Admin(1) Query: Admin Address",
			args: args{
				adminAddr:   acc1,
				requestAddr: acc1.String(),
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "Admin(1) Query: Non-admin Address",
			args: args{
				adminAddr:   acc1,
				requestAddr: acc2.String(),
			},
			errArgs: errArgs{
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
		{
			name: "Admins(1) Query: No Admins",
			args: args{
				addrs: []sdk.AccAddress{},
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "Admins(2) Query: Two Admins",
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

func (suite *KeeperTestSuite) TestQueryCertificate() {
	type args struct {
		denomID     string
		tokenID     string
		certificate *types.Certificate
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
			name: "Certificate(1) Query: Certificate Exists",
			args: args{
				denomID: "certificateauditing",
				tokenID: tokenID,
				certificate: &types.Certificate{
					Content:     content,
					Description: "",
					Certifier:   certifier.String(),
				},
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "Certificate(1) Query: Certificate Does Not Exists",
			args: args{
				denomID: "certificateauditing",
				tokenID: tokenID,
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "not found",
			},
		},
		{
			name: "Certificate(1) Query: Call on Non-Cert NFT",
			args: args{
				denomID: "NotCertNFT",
				tokenID: tokenID,
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "denom",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.args.certificate != nil {
				err := suite.keeper.IssueCertificate(suite.ctx, tc.args.denomID, tc.args.tokenID,
					tokenNm, tokenURI, *tc.args.certificate)
				suite.Require().NoError(err, tc.name)
			}
			certificateReq := types.QueryCertificateRequest{
				DenomId: tc.args.denomID,
				TokenId: tc.args.tokenID,
			}
			res, err := suite.queryClient.Certificate(gocontext.Background(), &certificateReq)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotNil(res.Certificate)
				suite.Require().Equal(*tc.args.certificate, res.Certificate)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryCertificates() {
	type args struct {
		certNFTs   []certNFT
		certifier  string
		denomID    string
		pagination *query.PageRequest
		wantLen    int
	}

	type errArgs struct {
		shouldPass bool
		contains   string
	}

	manyNFTs := []certNFT{
		{
			denomID:  "certificateauditing",
			tokenID:  tokenID,
			tokenNm:  tokenNm,
			tokenURI: tokenURI,
			certificate: types.Certificate{
				Content:     content,
				Description: "",
				Certifier:   certifier.String(),
			},
		},
		{
			denomID:  "certificateidentity",
			tokenID:  tokenID2,
			tokenNm:  tokenNm2,
			tokenURI: tokenURI2,
			certificate: types.Certificate{
				Content:     content,
				Description: "",
				Certifier:   certifier.String(),
			},
		},
		{
			denomID:  "certificateidentity",
			tokenID:  tokenID3,
			tokenNm:  tokenNm3,
			tokenURI: tokenURI3,
			certificate: types.Certificate{
				Content:     content,
				Description: "",
				Certifier:   certifier.String(),
			},
		},
		{
			denomID:  "certificateidentity",
			tokenID:  tokenID4,
			tokenNm:  tokenNm4,
			tokenURI: tokenURI4,
			certificate: types.Certificate{
				Content:     content,
				Description: "",
				Certifier:   certifier.String(),
			},
		},
		{
			denomID:  "certificateidentity",
			tokenID:  tokenID5,
			tokenNm:  tokenNm5,
			tokenURI: tokenURI5,
			certificate: types.Certificate{
				Content:     content,
				Description: "",
				Certifier:   certifier2.String(),
			},
		},
	}

	tests := []struct {
		name    string
		args    args
		errArgs errArgs
	}{
		{
			name: "Certificates(1) Query: One Certificate, One Result",
			args: args{
				certNFTs: []certNFT{
					{
						denomID:  "certificateauditing",
						tokenID:  tokenID,
						tokenNm:  tokenNm,
						tokenURI: tokenURI,
						certificate: types.Certificate{
							Content:     content,
							Description: "",
							Certifier:   certifier.String(),
						},
					},
				},
				certifier: certifier.String(),
				denomID:   "certificateauditing",
				wantLen:   1,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "Certificates(1) Query: One Certificate, No Results",
			args: args{
				certNFTs: []certNFT{
					{
						denomID:  "certificateauditing",
						tokenID:  tokenID,
						tokenNm:  tokenNm,
						tokenURI: tokenURI,
						certificate: types.Certificate{
							Content:     content,
							Description: "",
							Certifier:   certifier.String(),
						},
					},
				},
				certifier: certifier.String(),
				denomID:   "certificateidentity",
				wantLen:   0,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "Certificates(1) Query: Call on Non-Cert NFT",
			args: args{
				certNFTs: []certNFT{
					{
						denomID:  "certificateauditing",
						tokenID:  tokenID,
						tokenNm:  tokenNm,
						tokenURI: tokenURI,
						certificate: types.Certificate{
							Content:     content,
							Description: "",
							Certifier:   certifier.String(),
						},
					},
				},
				certifier: certifier.String(),
				denomID:   "NonCertNFT",
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "denom",
			},
		},
		{
			name: "Certificates(5) Query: Five Certificates, Three Results",
			args: args{
				certNFTs:  manyNFTs,
				certifier: certifier.String(),
				denomID:   "certificateidentity",
				wantLen:   3,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "Certificates(5) Query: Five Certificates, Pagination with Limit Two",
			args: args{
				certNFTs:  manyNFTs,
				certifier: certifier.String(),
				denomID:   "certificateidentity",
				pagination: &query.PageRequest{
					Offset: 0,
					Limit:  2,
				},
				wantLen: 2,
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
			for _, c := range tc.args.certNFTs {
				err := suite.keeper.IssueCertificate(suite.ctx, c.denomID, c.tokenID,
					c.tokenNm, c.tokenURI, c.certificate)
				suite.Require().NoError(err, tc.name)
			}
			certificatesReq := types.QueryCertificatesRequest{
				Certifier:  tc.args.certifier,
				DenomId:    tc.args.denomID,
				Pagination: tc.args.pagination,
			}
			res, err := suite.queryClient.Certificates(gocontext.Background(), &certificatesReq)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().NotNil(res)
				suite.Require().Len(res.Certificates, tc.args.wantLen)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}
