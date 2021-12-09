package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/test-go/testify/suite"
)

type KeysTestSuite struct {
	suite.Suite
}

func (suite *KeysTestSuite) TestCertifierStoreKey() {
	tests := []struct {
		keyBytes   []byte
		address    sdk.AccAddress
		shouldPass bool
	}{
		{
			keyBytes:   concat(certifierStoreKeyPrefix, []byte{10}),
			address:    sdk.AccAddress([]byte{10}),
			shouldPass: true,
		},
		{
			// different key value
			keyBytes:   concat(validatorStoreKeyPrefix, []byte{10}),
			address:    sdk.AccAddress([]byte{10}),
			shouldPass: false,
		},
	}

	for _, tc := range tests {
		certifierStoreKeyBytes := CertifierStoreKey(tc.address)
		if tc.shouldPass {
			suite.Require().Equal(tc.keyBytes, certifierStoreKeyBytes)
		} else {
			suite.Require().NotEqual(tc.keyBytes, certifierStoreKeyBytes)
		}
	}

}

func (suite *KeysTestSuite) TestCertifierAliasStoreKey() {
	tests := []struct {
		keyBytes   []byte
		alias      string
		shouldPass bool
	}{
		{
			keyBytes:   concat(certifierAliasStoreKeyPrefix, []byte("address1")),
			alias:      "address1",
			shouldPass: true,
		},
		{
			// different alias
			keyBytes:   concat(certifierAliasStoreKeyPrefix, []byte("address1")),
			alias:      "address2",
			shouldPass: false,
		},
		{
			// different key value
			keyBytes:   concat(certifierStoreKeyPrefix, []byte("address1")),
			alias:      "address1",
			shouldPass: false,
		},
	}

	for _, tc := range tests {
		certifierAliasStoreKeyBytes := CertifierAliasStoreKey(tc.alias)
		if tc.shouldPass {
			suite.Require().Equal(tc.keyBytes, certifierAliasStoreKeyBytes)
		} else {
			suite.Require().NotEqual(tc.keyBytes, certifierAliasStoreKeyBytes)
		}
	}

}

func (suite *KeysTestSuite) TestLibraryStoreKey() {
	tests := []struct {
		keyBytes   []byte
		address    sdk.AccAddress
		shouldPass bool
	}{
		{
			keyBytes:   concat(libraryStoreKeyPrefix, []byte{10}),
			address:    sdk.AccAddress([]byte{10}),
			shouldPass: true,
		},
		{
			// different key value
			keyBytes:   concat(validatorStoreKeyPrefix, []byte{10}),
			address:    sdk.AccAddress([]byte{10}),
			shouldPass: false,
		},
	}

	for _, tc := range tests {
		libraryStoreKeyBytes := LibraryStoreKey(tc.address)
		if tc.shouldPass {
			suite.Require().Equal(tc.keyBytes, libraryStoreKeyBytes)
		} else {
			suite.Require().NotEqual(tc.keyBytes, libraryStoreKeyBytes)
		}
	}

}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(KeysTestSuite))
}
