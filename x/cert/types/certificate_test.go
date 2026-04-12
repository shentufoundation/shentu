package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func TestCertificateTypeFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected types.CertificateType
	}{
		{"COMPILATION", types.CertificateTypeCompilation},
		{"compilation", types.CertificateTypeCompilation},
		{"CERT_TYPE_COMPILATION", types.CertificateTypeCompilation},
		{"AUDITING", types.CertificateTypeAuditing},
		{"auditing", types.CertificateTypeAuditing},
		{"CERT_TYPE_AUDITING", types.CertificateTypeAuditing},
		{"PROOF", types.CertificateTypeProof},
		{"proof", types.CertificateTypeProof},
		{"CERT_TYPE_PROOF", types.CertificateTypeProof},
		{"ORACLEOPERATOR", types.CertificateTypeOracleOperator},
		{"CERT_TYPE_ORACLE_OPERATOR", types.CertificateTypeOracleOperator},
		{"SHIELDPOOLCREATOR", types.CertificateTypeShieldPoolCreator},
		{"CERT_TYPE_SHIELD_POOL_CREATOR", types.CertificateTypeShieldPoolCreator},
		{"IDENTITY", types.CertificateTypeIdentity},
		{"identity", types.CertificateTypeIdentity},
		{"CERT_TYPE_IDENTITY", types.CertificateTypeIdentity},
		{"GENERAL", types.CertificateTypeGeneral},
		{"general", types.CertificateTypeGeneral},
		{"CERT_TYPE_GENERAL", types.CertificateTypeGeneral},
		{"BOUNTYADMIN", types.CertificateTypeBountyAdmin},
		{"CERT_TYPE_BOUNTY_ADMIN", types.CertificateTypeBountyAdmin},
		{"CERT_TYPE_BOUNTYADMIN", types.CertificateTypeBountyAdmin},
		{"OPENMATH", types.CertificateTypeOpenMath},
		{"openmath", types.CertificateTypeOpenMath},
		{"CERT_TYPE_OPENMATH", types.CertificateTypeOpenMath},
		// Invalid types.
		{"", types.CertificateTypeNil},
		{"unknown", types.CertificateTypeNil},
		{"INVALID", types.CertificateTypeNil},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			require.Equal(t, tc.expected, types.CertificateTypeFromString(tc.input))
		})
	}
}

func TestAssembleContent(t *testing.T) {
	validTypes := []struct {
		name    string
		content string
	}{
		{"compilation", "compiler-hash"},
		{"auditing", "audit-result"},
		{"proof", "proof-data"},
		{"identity", "identity-info"},
		{"general", "general-content"},
		{"openmath", "prover-addr"},
	}
	for _, tc := range validTypes {
		t.Run(tc.name, func(t *testing.T) {
			c := types.AssembleContent(tc.name, tc.content)
			require.NotNil(t, c)
			require.Equal(t, tc.content, c.GetContent())
		})
	}

	// Invalid type returns nil.
	c := types.AssembleContent("nonexistent", "data")
	require.Nil(t, c)
}

func TestNewCertificate_InvalidType(t *testing.T) {
	addr := []byte("test_addr_for_cert__")
	_, err := types.NewCertificate("nonexistent", "content", "", "", "", addr)
	require.Error(t, err)
}

func TestNewCertificate_ValidTypes(t *testing.T) {
	addr := []byte("test_addr_for_cert__")
	for _, certType := range []string{"general", "auditing", "proof", "identity", "openmath"} {
		t.Run(certType, func(t *testing.T) {
			cert, err := types.NewCertificate(certType, "test-content", "", "", "description", addr)
			require.NoError(t, err)
			require.Equal(t, "test-content", cert.GetContentString())
			require.Equal(t, "description", cert.Description)
		})
	}
}

func TestAddOrRemoveFromString(t *testing.T) {
	tests := []struct {
		input   string
		want    types.AddOrRemove
		wantErr bool
	}{
		{"add", types.Add, false},
		{"Add", types.Add, false},
		{"ADD", types.Add, false},
		{"remove", types.Remove, false},
		{"Remove", types.Remove, false},
		{"REMOVE", types.Remove, false},
		{"", types.Add, true},
		{"invalid", types.Add, true},
		{"delete", types.Add, true},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := types.AddOrRemoveFromString(tc.input)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestAddOrRemove_String(t *testing.T) {
	require.Equal(t, "add", types.Add.String())
	require.Equal(t, "remove", types.Remove.String())
}

func TestAddOrRemove_MarshalJSON(t *testing.T) {
	bz, err := types.Add.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"add"`, string(bz))

	bz, err = types.Remove.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, `"remove"`, string(bz))
}

func TestAddOrRemove_UnmarshalJSON(t *testing.T) {
	var aor types.AddOrRemove

	require.NoError(t, aor.UnmarshalJSON([]byte(`"add"`)))
	require.Equal(t, types.Add, aor)

	require.NoError(t, aor.UnmarshalJSON([]byte(`"remove"`)))
	require.Equal(t, types.Remove, aor)

	require.Error(t, aor.UnmarshalJSON([]byte(`"invalid"`)))
	require.Error(t, aor.UnmarshalJSON([]byte(`not-json`)))
}
