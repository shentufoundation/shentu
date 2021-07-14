package keeper_test

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/nft/types"
)

type certNFT struct {
	denomID     string
	tokenID     string
	tokenNm     string
	tokenURI    string
	certificate types.Certificate
}

func (suite *KeeperTestSuite) TestCertificate_Issue() {
	type args struct {
		cert certNFT
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
			name: "IssueCertificate: Correct Issue",
			args: args{
				cert: certNFT{
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
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "IssueCertificate: Unqualified Certifier",
			args: args{
				cert: certNFT{
					denomID:  "certificateauditing",
					tokenID:  tokenID,
					tokenNm:  tokenNm,
					tokenURI: tokenURI,
					certificate: types.Certificate{
						Content:     content,
						Description: "",
						Certifier:   acc1.String(),
					},
				},
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "certifier",
			},
		},
		{
			name: "IssueCertificate: Invalid DenomID",
			args: args{
				cert: certNFT{
					denomID:  "CertificateInvalid",
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
			errArgs: errArgs{
				shouldPass: false,
				contains:   "denom",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.IssueCertificate(suite.ctx, tc.args.cert.denomID, tc.args.cert.tokenID,
				tc.args.cert.tokenNm, tc.args.cert.tokenURI, tc.args.cert.certificate)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				certNFT, err := suite.keeper.GetNFT(suite.ctx, tc.args.cert.denomID, tc.args.cert.tokenID)
				suite.Require().NoError(err, tc.name)
				gotCert := suite.keeper.UnmarshalCertificate(suite.ctx, certNFT.GetData())
				suite.Require().Equal(tc.args.cert.certificate, gotCert)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCertificate_Edit() {
	type args struct {
		cert certNFT
		edit certNFT
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
			name: "EditCertificate: Correct Edit",
			args: args{
				cert: certNFT{
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
				edit: certNFT{
					denomID:  "certificateauditing",
					tokenID:  tokenID,
					tokenNm:  tokenNm2,
					tokenURI: tokenURI2,
					certificate: types.Certificate{
						Content:     content,
						Description: "",
						Certifier:   certifier.String(),
					},
				},
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "EditCertificate: Token Does Not Exist",
			args: args{
				cert: certNFT{
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
				edit: certNFT{
					denomID:  "certificateauditing",
					tokenID:  tokenID2,
					tokenNm:  tokenNm2,
					tokenURI: tokenURI2,
					certificate: types.Certificate{
						Content:     content,
						Description: "",
						Certifier:   certifier.String(),
					},
				},
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "not found",
			},
		},
		{
			name: "EditCertificate: Edit Non-Cert NFT",
			args: args{
				cert: certNFT{
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
				edit: certNFT{
					denomID:  "NotCertNFT",
					tokenID:  tokenID,
					tokenNm:  tokenNm,
					tokenURI: tokenURI,
				},
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "denom",
			},
		},
		{
			name: "EditCertificate: Attempts to Edit Certifier",
			args: args{
				cert: certNFT{
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
				edit: certNFT{
					denomID:  "certificateauditing",
					tokenID:  tokenID,
					tokenNm:  tokenNm,
					tokenURI: tokenURI,
					certificate: types.Certificate{
						Content:     content,
						Description: "",
						Certifier:   acc1.String(),
					},
				},
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "certifier",
			},
		},
		{
			name: "EditCertificate: Attempts to Edit Denom",
			args: args{
				cert: certNFT{
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
				edit: certNFT{
					denomID:  "certificateidentity",
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
			errArgs: errArgs{
				shouldPass: false,
				contains:   "denom",
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			err := suite.keeper.IssueCertificate(suite.ctx, tc.args.cert.denomID, tc.args.cert.tokenID,
				tc.args.cert.tokenNm, tc.args.cert.tokenURI, tc.args.cert.certificate)
			suite.Require().NoError(err, tc.name)
			err = suite.keeper.EditCertificate(suite.ctx, tc.args.edit.denomID, tc.args.edit.tokenID,
				tc.args.edit.tokenNm, tc.args.edit.tokenURI, tc.args.edit.certificate)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				certNFT, err := suite.keeper.GetNFT(suite.ctx, tc.args.edit.denomID, tc.args.edit.tokenID)
				suite.Require().NoError(err, tc.name)
				gotCert := suite.keeper.UnmarshalCertificate(suite.ctx, certNFT.GetData())
				suite.Require().Equal(tc.args.edit.certificate, gotCert)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCertificate_Revoke() {
	type args struct {
		cert    certNFT
		revoke  certNFT
		revoker sdk.AccAddress
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
			name: "RevokeCertificate: Correct Revoke",
			args: args{
				cert: certNFT{
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
				revoke: certNFT{
					denomID: "certificateauditing",
					tokenID: tokenID,
					tokenNm: tokenNm,
				},
				revoker: certifier,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "RevokeCertificate: Unqualified Revoker",
			args: args{
				cert: certNFT{
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
				revoke: certNFT{
					denomID: "certificateauditing",
					tokenID: tokenID,
					tokenNm: tokenNm,
				},
				revoker: acc1,
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "certifier",
			},
		},
		{
			name: "RevokeCertificate: Token Does Not Exist",
			args: args{
				cert: certNFT{
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
				revoke: certNFT{
					denomID: "certificateauditing",
					tokenID: tokenID2,
					tokenNm: tokenNm,
				},
				revoker: certifier,
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "not found",
			},
		},
		{
			name: "RevokeCertificate: Revoke Non-Cert NFT",
			args: args{
				cert: certNFT{
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
				revoke: certNFT{
					denomID: "NotCertNFT",
					tokenID: tokenID,
					tokenNm: tokenNm,
				},
				revoker: certifier,
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
			err := suite.keeper.IssueCertificate(suite.ctx, tc.args.cert.denomID, tc.args.cert.tokenID,
				tc.args.cert.tokenNm, tc.args.cert.tokenURI, tc.args.cert.certificate)
			suite.Require().NoError(err, tc.name)
			err = suite.keeper.RevokeCertificate(suite.ctx, tc.args.revoke.denomID, tc.args.revoke.tokenID,
				tc.args.revoker)
			if tc.errArgs.shouldPass {
				suite.Require().NoError(err, tc.name)
				_, err := suite.keeper.GetNFT(suite.ctx, tc.args.revoke.denomID, tc.args.revoke.tokenID)
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), "not found"))
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
			}
		})
	}
}
