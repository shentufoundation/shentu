package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

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
	proposer := validShentuAddr(t, "proposer_")

	tests := []struct {
		name    string
		msg     *types.MsgUpdateCertifier
		wantErr bool
	}{
		{
			"valid add",
			types.NewMsgUpdateCertifier(authority, certifier, "desc", types.Add, proposer),
			false,
		},
		{
			"valid remove",
			types.NewMsgUpdateCertifier(authority, certifier, "desc", types.Remove, nil),
			false,
		},
		{
			"invalid authority",
			&types.MsgUpdateCertifier{
				Authority: "invalid",
				Certifier: certifier.String(),
				Operation: "add",
			},
			true,
		},
		{
			"empty certifier",
			&types.MsgUpdateCertifier{
				Authority: authority.String(),
				Certifier: "",
				Operation: "add",
			},
			true,
		},
		{
			"invalid operation",
			types.NewMsgUpdateCertifier(authority, certifier, "", types.AddOrRemove(false), nil),
			false, // "add" is valid
		},
		{
			"bad operation string",
			&types.MsgUpdateCertifier{
				Authority: authority.String(),
				Certifier: certifier.String(),
				Operation: "invalid",
			},
			true,
		},
		{
			"invalid proposer",
			&types.MsgUpdateCertifier{
				Authority: authority.String(),
				Certifier: certifier.String(),
				Operation: "add",
				Proposer:  "invalid",
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
	msg := types.NewMsgUpdateCertifier(authority, certifier, "", types.Add, nil)

	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, authority, signers[0])
}

func TestMsgIssueCertificate_ValidateBasic(t *testing.T) {
	// ValidateBasic currently returns nil unconditionally.
	addr := validShentuAddr(t, "certifier")
	content := types.AssembleContent("general", "some-content")
	msg := types.NewMsgIssueCertificate(content, "", "", "desc", addr)
	require.NoError(t, msg.ValidateBasic())
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
	msg := types.NewMsgRevokeCertificate(addr, 1, "reason")
	require.NoError(t, msg.ValidateBasic())
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
	msg := types.NewMsgUpdateCertifier(authority, certifier, "", types.Add, nil)
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
