package keeper_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdksimapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/v2/simapp"
	"github.com/certikfoundation/shentu/v2/x/cert/keeper"
	"github.com/certikfoundation/shentu/v2/x/cert/types"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	acc4 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

type cert struct {
	certTypeStr  string
	contStr      string
	compiler     string
	bytecodeHash string
	description  string
	certifier    sdk.AccAddress
	delete       bool
	assumption   bool
	create       bool
	inputCertId  uint64
}

// shared setup
type KeeperTestSuite struct {
	suite.Suite

	// cdc    *codec.LegacyAmino
	app         *simapp.SimApp
	ctx         sdk.Context
	keeper      keeper.Keeper
	address     []sdk.AccAddress
	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.CertKeeper
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, keeper.Querier{Keeper: suite.keeper})
	suite.queryClient = types.NewQueryClient(queryHelper)
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
	suite.keeper.SetCertifier(suite.ctx, types.NewCertifier(suite.address[0], "", suite.address[0], ""))

}

func (suite *KeeperTestSuite) TestCertificate_GetSet() {
	type args struct {
		cert []cert
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
		{"Certificate(1) Create -> Get: Simple",
			args{
				cert: []cert{
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
					},
				},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Certificate(2) Create -> Get: Two Different Ones",
			args{
				cert: []cert{
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
					},
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash1",
						compiler:     "compiler2",
						bytecodeHash: "bytecodehash2",
						description:  "",
						certifier:    suite.address[1],
					},
				},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		suite.Run(tc.name, func() {
			// construct a new cert
			for _, cert := range tc.args.cert {
				want, err := types.NewCertificate(cert.certTypeStr, cert.contStr, cert.compiler, cert.bytecodeHash, cert.description, cert.certifier)
				suite.Require().NoError(err, tc.name)
				// Get the next available ID and assign it
				id := suite.keeper.GetNextCertificateID(suite.ctx)
				want.CertificateId = id
				// set the cert and its ID in the store
				suite.keeper.SetNextCertificateID(suite.ctx, id+1)
				suite.keeper.SetCertificate(suite.ctx, want)
				// now retrieve its ID from the store
				got, err := suite.keeper.GetCertificateByID(suite.ctx, id)
				if tc.errArgs.shouldPass {
					suite.Require().NoError(err, tc.name)
					suite.Equal(got, want)
				} else {
					suite.Require().Error(err, tc.name)
					suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCertificate_Delete() {
	type args struct {
		cert []cert
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
		{"Certificate(1) Delete: Simple",
			args{
				cert: []cert{
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						delete:       true,
						inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx),
					},
				},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Certificate(2) Delete: Add Three Delete the Second One",
			args{
				cert: []cert{
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						delete:       false,
						inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx),
					},
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						delete:       true,
						inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx) + 1,
					},
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						delete:       false,
						inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx) + 2,
					},
				},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Certificate(3) Delete: Invalid certificate id",
			args{
				cert: []cert{
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						delete:       true,
						inputCertId:  suite.keeper.GetNextCertificateID(suite.ctx) + 10,
					},
				},
			},
			errArgs{
				shouldPass: false,
				contains:   "certificate id does not exist",
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		suite.Run(tc.name, func() {
			// initializing certificate id for each testcase
			suite.keeper.SetNextCertificateID(suite.ctx, 1)
			// construct a new cert
			for _, cert := range tc.args.cert {
				want, err := types.NewCertificate(cert.certTypeStr, cert.contStr, cert.compiler, cert.bytecodeHash, cert.description, cert.certifier)
				suite.Require().NoError(err, tc.name)
				// Get the next available ID and assign it
				id := suite.keeper.GetNextCertificateID(suite.ctx)
				want.CertificateId = id
				// set the cert and its ID in the store
				suite.keeper.SetNextCertificateID(suite.ctx, id+1)
				suite.keeper.SetCertificate(suite.ctx, want)
				// Check if the certificate was set successfully and delete if marked
				if suite.keeper.HasCertificateByID(suite.ctx, want.CertificateId) && cert.delete {
					err := suite.keeper.DeleteCertificate(suite.ctx, want)
					suite.Require().NoError(err, tc.name)
				}
				// now retrieve its ID from the store
				got, err := suite.keeper.GetCertificateByID(suite.ctx, cert.inputCertId)
				if tc.errArgs.shouldPass {
					if cert.delete {
						suite.Require().Error(err, tc.name)
						suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
					} else {
						suite.Require().NoError(err, tc.name)
						suite.Equal(got, want)
					}
				} else {
					suite.Require().Error(err, tc.name)
					suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestCertificate_IsCertified_Issue() {
	type args struct {
		cert []cert
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
		{"Certificate(1) IsCertified: Simple",
			args{
				cert: []cert{
					{
						certTypeStr:  "auditing",
						contStr:      "certik1k4gj07sgy6x3k6ms31aztgu9aajjkaw3ktsydag",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						assumption:   false,
						create:       false,
					},
				},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Certificate(2) IsCertified: Invalid Cert Type",
			args{
				cert: []cert{
					{
						certTypeStr:  "random",
						contStr:      "certik1k4gj07sgy6x3k6ms31aztgu9aajjkaw3ktsydag",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						assumption:   true,
						create:       false,
					},
				},
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
		{"Certificate(1) IsCertified: Issue a New Cert for A Non-Certified Cert",
			args{
				cert: []cert{
					{
						certTypeStr:  "shieldpoolcreator",
						contStr:      "certik1k4gj07sgy6x3k6ms31aztgu9aajjkaw3ktsydag",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						// assumption (i.e. intention) must be false for
						// a non-certified cert to get certified
						assumption: false,
						create:     true,
					},
				},
			},
			errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{"Certificate(1) IsCertified: Issue a New Cert for A Non-Certified Cert with Bad Intention",
			args{
				cert: []cert{
					{
						certTypeStr:  "oracleoperator",
						contStr:      "certik1k4gj07sgy6x3k6ms31aztgu9aajjkaw3ktsydag",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						// assumption (i.e. intention) must be false for
						// a non-certified cert to get certified
						assumption: true,
						create:     true,
					},
				},
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
		{"Certificate(1) IsCertified: Issue a New Cert for A Non-Certified Cert from An Unqualified Certifier",
			args{
				cert: []cert{
					{
						certTypeStr:  "shieldpoolcreator",
						contStr:      "certik1k4gj07sgy6x3k6ms31aztgu9aajjkaw3ktsydag",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[1], // only certifier is address[0]
						// assumption (i.e. intention) must be false for
						// a non-certified cert to get certified
						assumption: false,
						create:     true,
					},
				},
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		suite.Run(tc.name, func() {
			for _, cert := range tc.args.cert {
				got := suite.keeper.IsCertified(suite.ctx, cert.contStr, cert.certTypeStr)
				if !got && !cert.assumption && cert.create {
					// construct a new cert
					newCert, err := types.NewCertificate(cert.certTypeStr, cert.contStr, cert.compiler, cert.bytecodeHash, cert.description, cert.certifier)
					suite.Require().NoError(err, tc.name)
					// issue a new cert
					_, err = suite.keeper.IssueCertificate(suite.ctx, newCert)
					if !reflect.DeepEqual(cert.certifier, suite.address[0]) {
						suite.Require().Error(err, tc.name)
					} else {
						suite.Require().NoError(err, tc.name)
					}

					// check again if Certified
					got = suite.keeper.IsCertified(suite.ctx, cert.contStr, cert.certTypeStr)
					want := !cert.assumption
					if tc.errArgs.shouldPass {
						suite.Require().Equal(got, want)
					} else {
						suite.Require().NotEqual(got, want)
						// suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
					}
				} else {
					want := cert.assumption
					if tc.errArgs.shouldPass {
						suite.Require().Equal(got, want)
					} else {
						suite.Require().NotEqual(got, want)
						// suite.Require().True(strings.Contains(err.Error(), tc.errArgs.contains))
					}
				}
			}
		})
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
