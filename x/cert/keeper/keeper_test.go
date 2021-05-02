package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/certikfoundation/shentu/simapp"
	"github.com/certikfoundation/shentu/x/cert/keeper"
	"github.com/certikfoundation/shentu/x/cert/types"
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
}

// shared setup
type KeeperTestSuite struct {
	suite.Suite

	// cdc    *codec.LegacyAmino
	app     *simapp.SimApp
	ctx     sdk.Context
	keeper  keeper.Keeper
	address []sdk.AccAddress
	// queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = simapp.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = suite.app.CertKeeper

	// queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	// types.RegisterQueryServer(queryHelper, suite.app.CertKeeper)
	// suite.queryClient = types.NewQueryClient(queryHelper)

	for _, acc := range []sdk.AccAddress{acc1, acc2, acc3, acc4} {
		err := suite.app.BankKeeper.AddCoins(
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
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestCertificateGetSet() {
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
		{"Certificate(1) Creat -> Get: Simple",
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
		{"Certificate(2) Creat -> Get: Two Different Ones",
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
		suite.Run(tc.name, func() {
			suite.SetupTest()
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

func (suite *KeeperTestSuite) TestCertificateDelete() {
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
					},
				},
			},
			errArgs{
				shouldPass: false,
				contains:   "",
			},
		},
		{"Certificate(1) Delete: Add Three Delete the Second One",
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
					},
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						delete:       false,
					},
					{
						certTypeStr:  "compilation",
						contStr:      "sourcodehash0",
						compiler:     "compiler1",
						bytecodeHash: "bytecodehash1",
						description:  "",
						certifier:    suite.address[0],
						delete:       true,
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
		suite.Run(tc.name, func() {
			suite.SetupTest()
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
				// delete if marked
				if cert.delete {
					err := suite.keeper.DeleteCertificate(suite.ctx, want)
					suite.Require().NoError(err, tc.name)
				}
				// now retrieve its ID from the store
				got, err := suite.keeper.GetCertificateByID(suite.ctx, id)
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
