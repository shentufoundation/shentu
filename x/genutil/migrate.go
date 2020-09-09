package genutil

import (
	"github.com/cosmos/cosmos-sdk/codec"
	v038auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v0_38"
	v039auth "github.com/cosmos/cosmos-sdk/x/auth/legacy/v0_39"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/oracle"
)

// Migrate makes a state exported from v0.11 compatible with v0.13
// NOTE: It simply removes certificates and certificate ID fields,
// which are incompatible with v0.13.
func Migrate(appState types.AppMap) types.AppMap {
	oldCodec := codec.New()
	codec.RegisterCrypto(oldCodec)
	v038auth.RegisterCodec(oldCodec)

	newCodec := codec.New()
	codec.RegisterCrypto(newCodec)
	v039auth.RegisterCodec(newCodec)
	cert.RegisterCodec(newCodec)

	// migrate auth state to cosmos sdk v0.39
	if appState[v038auth.ModuleName] != nil {
		var authGenState v038auth.GenesisState
		oldCodec.MustUnmarshalJSON(appState[v038auth.ModuleName], &authGenState)
		appState[v038auth.ModuleName] = newCodec.MustMarshalJSON(v039auth.Migrate(authGenState))
	}

	// remove certificates and certificate ID field
	if appState[cert.ModuleName] != nil {
		type CertMigrationGenesisState struct {
			Certifiers []cert.Certifier `json:"certifiers"`
			Validators []cert.Validator `json:"validators"`
			Platforms  []cert.Platform  `json:"platforms"`
			Libraries  []cert.Library   `json:"libraries"`
		}
		var certGenState CertMigrationGenesisState
		newCodec.MustUnmarshalJSON(appState[cert.ModuleName], &certGenState)
		appState[cert.ModuleName] = newCodec.MustMarshalJSON(certGenState)
	}

	var oracleGenState = oracle.DefaultGenesisState()
	appState[oracle.ModuleName] = newCodec.MustMarshalJSON(oracleGenState)

	return appState
}
