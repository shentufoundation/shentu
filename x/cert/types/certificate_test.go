package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/test-go/testify/suite"
)

var (
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
)

type cert struct {
	certTypeStr  string
	contStr      string
	compiler     string
	bytecodeHash string
	description  string
	certifier    sdk.AccAddress
}

type TypesTestSuite struct {
	suite.Suite
}

func (suite *TypesTestSuite) TestNewCertificate() {
	tests := []struct {
		cert       cert
		shouldPass bool
	}{
		{
			cert: cert{
				certTypeStr:  "compilation",
				contStr:      "sourcodehash0",
				compiler:     "compiler1",
				bytecodeHash: "bytecodehash1",
				description:  "",
				certifier:    acc1,
			},
			shouldPass: true,
		},
		{
			cert: cert{
				certTypeStr:  "invalid",
				contStr:      "sourcodehash0",
				compiler:     "compiler1",
				bytecodeHash: "bytecodehash1",
				description:  "",
				certifier:    acc1,
			},
			shouldPass: false,
		},
	}

	for _, tc := range tests {
		certificate, err := NewCertificate(tc.cert.certTypeStr, tc.cert.contStr, tc.cert.compiler, tc.cert.bytecodeHash, tc.cert.description, tc.cert.certifier)
		if tc.shouldPass {
			suite.Require().NoError(err)
			suite.Require().Equal(tc.cert.certifier.String(), certificate.Certifier)
		} else {
			suite.Require().Error(err)
		}
	}
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}
