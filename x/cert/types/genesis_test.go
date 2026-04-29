package types_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// validAddr returns a deterministic, well-formed bech32 address for the
// currently configured account prefix.
func validAddr(seed string) string {
	padded := make([]byte, 20)
	copy(padded, seed)
	return sdk.AccAddress(padded).String()
}

func TestValidateGenesis(t *testing.T) {
	certifierA := validAddr("certifier-a")
	certifierB := validAddr("certifier-b")

	cert1 := types.Certificate{CertificateId: 1}
	cert2 := types.Certificate{CertificateId: 2}
	cert5 := types.Certificate{CertificateId: 5}

	cases := []struct {
		name    string
		gs      types.GenesisState
		wantErr string
	}{
		{
			name: "empty state with default next id",
			gs:   types.GenesisState{NextCertificateId: 1},
		},
		{
			name: "valid state",
			gs: types.GenesisState{
				Certifiers:        []types.Certifier{{Address: certifierA}},
				Certificates:      []types.Certificate{cert1, cert2, cert5},
				NextCertificateId: 6,
			},
		},
		{
			name: "zero certificate id is rejected",
			gs: types.GenesisState{
				Certificates:      []types.Certificate{{CertificateId: 0}},
				NextCertificateId: 1,
			},
			wantErr: "zero id",
		},
		{
			name: "duplicate certificate id is rejected",
			gs: types.GenesisState{
				Certificates:      []types.Certificate{cert1, cert1},
				NextCertificateId: 2,
			},
			wantErr: "duplicate certificate id",
		},
		{
			name: "next id equal to max imported id is rejected",
			gs: types.GenesisState{
				Certificates:      []types.Certificate{cert5},
				NextCertificateId: 5,
			},
			wantErr: "must be strictly greater",
		},
		{
			name: "next id less than max imported id is rejected",
			gs: types.GenesisState{
				Certificates:      []types.Certificate{cert5},
				NextCertificateId: 3,
			},
			wantErr: "must be strictly greater",
		},
		{
			name:    "next id zero with no certificates is rejected",
			gs:      types.GenesisState{NextCertificateId: 0},
			wantErr: "must be strictly greater",
		},
		{
			name: "empty certifier address is rejected",
			gs: types.GenesisState{
				Certifiers:        []types.Certifier{{Address: ""}},
				NextCertificateId: 1,
			},
			wantErr: "empty address",
		},
		{
			name: "malformed certifier address is rejected",
			gs: types.GenesisState{
				Certifiers:        []types.Certifier{{Address: "not-a-bech32-address"}},
				NextCertificateId: 1,
			},
			wantErr: "invalid certifier address",
		},
		{
			name: "duplicate certifier is rejected",
			gs: types.GenesisState{
				Certifiers: []types.Certifier{
					{Address: certifierA},
					{Address: certifierB},
					{Address: certifierA},
				},
				NextCertificateId: 1,
			},
			wantErr: "duplicate certifier",
		},
		{
			name: "case-variant duplicate certifier is rejected",
			gs: types.GenesisState{
				Certifiers: []types.Certifier{
					{Address: certifierA},
					{Address: strings.ToUpper(certifierA)},
				},
				NextCertificateId: 1,
			},
			wantErr: "duplicate certifier",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.ValidateGenesis(tc.gs)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}
