package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

func validShentuAddr(t *testing.T, seed string) sdk.AccAddress {
	t.Helper()
	// Pad seed to 20 bytes for a valid AccAddress.
	padded := make([]byte, 20)
	copy(padded, seed)
	return sdk.AccAddress(padded)
}

func TestMsgUpdateCertifier_ValidateBasic(t *testing.T) {
	authority := validShentuAddr(t, "authority")
	certifier := validShentuAddr(t, "certifier")

	tests := []struct {
		name    string
		msg     *types.MsgUpdateCertifier
		wantErr bool
	}{
		{
			"valid add",
			types.NewMsgUpdateCertifier(authority, certifier, "desc", types.Add),
			false,
		},
		{
			"valid remove",
			types.NewMsgUpdateCertifier(authority, certifier, "desc", types.Remove),
			false,
		},
		{
			"invalid authority",
			&types.MsgUpdateCertifier{
				Authority: "invalid",
				Certifier: certifier.String(),
				Operation: types.CertifierUpdateOperationAdd,
			},
			true,
		},
		{
			"empty certifier",
			&types.MsgUpdateCertifier{
				Authority: authority.String(),
				Certifier: "",
				Operation: types.CertifierUpdateOperationAdd,
			},
			true,
		},
		{
			"invalid operation",
			&types.MsgUpdateCertifier{
				Authority: authority.String(),
				Certifier: certifier.String(),
				Operation: types.CertifierUpdateOperationUnspecified,
			},
			true,
		},
		{
			"bad operation enum",
			&types.MsgUpdateCertifier{
				Authority: authority.String(),
				Certifier: certifier.String(),
				Operation: types.CertifierUpdateOperation(99),
			},
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateCertifier_GetSigners(t *testing.T) {
	authority := validShentuAddr(t, "authority")
	certifier := validShentuAddr(t, "certifier")
	msg := types.NewMsgUpdateCertifier(authority, certifier, "", types.Add)

	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, authority, signers[0])
}

func TestMsgIssueCertificate_ValidateBasic(t *testing.T) {
	addr := validShentuAddr(t, "certifier")
	tests := []struct {
		name    string
		msg     *types.MsgIssueCertificate
		wantErr bool
	}{
		{
			name:    "valid general",
			msg:     types.NewMsgIssueCertificate(types.AssembleContent("general", "some-content"), "", "", "desc", addr),
			wantErr: false,
		},
		{
			name:    "valid compilation",
			msg:     types.NewMsgIssueCertificate(types.AssembleContent("compilation", "source-hash"), "solc", "0xdeadbeef", "desc", addr),
			wantErr: false,
		},
		{
			name: "invalid certifier",
			msg: &types.MsgIssueCertificate{
				Content:   mustAny(t, types.AssembleContent("general", "some-content")),
				Certifier: "invalid",
			},
			wantErr: true,
		},
		{
			name: "missing content",
			msg: &types.MsgIssueCertificate{
				Certifier: addr.String(),
			},
			wantErr: true,
		},
		{
			name: "invalid content type",
			msg: &types.MsgIssueCertificate{
				Content:   &codectypes.Any{TypeUrl: "/shentu.cert.v1alpha1.NotAContent"},
				Certifier: addr.String(),
			},
			wantErr: true,
		},
		{
			name: "compilation missing compiler",
			msg: &types.MsgIssueCertificate{
				Content:      mustAny(t, types.AssembleContent("compilation", "source-hash")),
				BytecodeHash: "0xdeadbeef",
				Certifier:    addr.String(),
			},
			wantErr: true,
		},
		{
			name: "compilation missing bytecode hash",
			msg: &types.MsgIssueCertificate{
				Content:   mustAny(t, types.AssembleContent("compilation", "source-hash")),
				Compiler:  "solc",
				Certifier: addr.String(),
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgIssueCertificate_GetSigners(t *testing.T) {
	addr := validShentuAddr(t, "certifier")
	content := types.AssembleContent("general", "content")
	msg := types.NewMsgIssueCertificate(content, "", "", "", addr)

	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0])
}

func TestMsgRevokeCertificate_ValidateBasic(t *testing.T) {
	addr := validShentuAddr(t, "revoker_addr____")
	tests := []struct {
		name    string
		msg     *types.MsgRevokeCertificate
		wantErr bool
	}{
		{
			name:    "valid",
			msg:     types.NewMsgRevokeCertificate(addr, 1, "reason"),
			wantErr: false,
		},
		{
			name: "invalid revoker",
			msg: &types.MsgRevokeCertificate{
				Revoker: "invalid",
				Id:      1,
			},
			wantErr: true,
		},
		{
			name:    "zero id",
			msg:     types.NewMsgRevokeCertificate(addr, 0, "reason"),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgRevokeCertificate_GetSigners(t *testing.T) {
	addr := validShentuAddr(t, "revoker_addr____")
	msg := types.NewMsgRevokeCertificate(addr, 1, "reason")

	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, addr, signers[0])
}

func TestMsgUpdateCertifier_RouteAndType(t *testing.T) {
	authority := validShentuAddr(t, "authority")
	certifier := validShentuAddr(t, "certifier")
	msg := types.NewMsgUpdateCertifier(authority, certifier, "", types.Add)
	require.Equal(t, types.ModuleName, msg.Route())
	require.Equal(t, "update_certifier", msg.Type())
}

func TestMsgIssueCertificate_RouteAndType(t *testing.T) {
	addr := validShentuAddr(t, "certifier")
	content := types.AssembleContent("general", "content")
	msg := types.NewMsgIssueCertificate(content, "", "", "", addr)
	require.Equal(t, types.ModuleName, msg.Route())
	require.Equal(t, "issue_certificate", msg.Type())
}

func TestMsgRevokeCertificate_RouteAndType(t *testing.T) {
	addr := validShentuAddr(t, "revoker_addr____")
	msg := types.NewMsgRevokeCertificate(addr, 1, "")
	require.Equal(t, types.ModuleName, msg.Route())
	require.Equal(t, "revoke_certificate", msg.Type())
}

func mustAny(t *testing.T, content types.Content) *codectypes.Any {
	t.Helper()
	require.NotNil(t, content)

	any, err := codectypes.NewAnyWithValue(content)
	require.NoError(t, err)
	return any
}
