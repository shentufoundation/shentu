package types_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/certikfoundation/shentu/x/nft/types"
)

var (
	acc1      = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String()
	acc2      = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String()
	certifier = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())

	tokenID  = "tokenid"
	tokenNm  = "tokennm"
	tokenURI = "https://google.com/token.json"
	content  = "content"
)

func TestMsgCreateAdmin(t *testing.T) {
	type args struct {
		creator string
		address string
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
			name: "MsgCreateAdmin: Valid",
			args: args{
				creator: acc1,
				address: acc2,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "MsgCreateAdmin: Missing Creator",
			args: args{
				address: acc2,
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "empty address",
			},
		},
		{
			name: "MsgCreateAdmin: Missing Address",
			args: args{
				creator: acc1,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		func() {
			defer func() {
				if r := recover(); r != nil {
					err := r.(error)
					require.False(t, tc.errArgs.shouldPass)
					require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}()
			msg := types.NewMsgCreateAdmin(tc.args.creator, tc.args.address)
			err := msg.ValidateBasic()
			if tc.errArgs.shouldPass {
				require.NoError(t, err, tc.name)
			} else {
				require.Error(t, err, tc.name)
				require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
			}
		}()
	}
}

func TestMsgRevokeAdmin(t *testing.T) {
	type args struct {
		issuer  string
		address string
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
			name: "MsgRevokeAdmin: Valid",
			args: args{
				issuer:  acc1,
				address: acc2,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "MsgRevokeAdmin: Missing Issuer",
			args: args{
				address: acc2,
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "empty address",
			},
		},
		{
			name: "MsgRevokeAdmin: Missing Address",
			args: args{
				issuer: acc1,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
	}
	for _, tc := range tests {
		func() {
			defer func() {
				if r := recover(); r != nil {
					err := r.(error)
					require.False(t, tc.errArgs.shouldPass)
					require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}()
			msg := types.NewMsgRevokeAdmin(tc.args.issuer, tc.args.address)
			err := msg.ValidateBasic()
			if tc.errArgs.shouldPass {
				require.NoError(t, err, tc.name)
			} else {
				require.Error(t, err, tc.name)
				require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
			}
		}()
	}
}

func TestMsgIssueCertificate(t *testing.T) {
	type args struct {
		denomID     string
		tokenID     string
		name        string
		uri         string
		content     string
		description string
		certifier   sdk.AccAddress
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
			name: "MsgIssueCertificate: Valid",
			args: args{
				denomID:     "certificateauditing",
				tokenID:     tokenID,
				name:        tokenNm,
				uri:         tokenURI,
				content:     content,
				description: "",
				certifier:   certifier,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "MsgIssueCertificate: Missing Certifier",
			args: args{
				denomID:     "certificateauditing",
				tokenID:     tokenID,
				name:        tokenNm,
				uri:         tokenURI,
				content:     content,
				description: "",
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "empty address",
			},
		},
	}

	for _, tc := range tests {
		func() {
			defer func() {
				if r := recover(); r != nil {
					err := r.(error)
					require.False(t, tc.errArgs.shouldPass)
					require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}()
			msg := types.NewMsgIssueCertificate(tc.args.denomID, tc.args.tokenID, tc.args.name,
				tc.args.uri, tc.args.content, tc.args.description, tc.args.certifier)
			err := msg.ValidateBasic()
			if tc.errArgs.shouldPass {
				require.NoError(t, err, tc.name)
			} else {
				require.Error(t, err, tc.name)
				require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
			}
		}()
	}
}

func TestMsgEditCertificate(t *testing.T) {
	type args struct {
		denomID     string
		tokenID     string
		name        string
		uri         string
		content     string
		description string
		owner       sdk.AccAddress
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
			name: "MsgEditCertificate: Valid",
			args: args{
				denomID:     "certificateauditing",
				tokenID:     tokenID,
				name:        tokenNm,
				uri:         tokenURI,
				content:     content,
				description: "",
				owner:       certifier,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "MsgEditCertificate: Missing Owner",
			args: args{
				denomID:     "certificateauditing",
				tokenID:     tokenID,
				name:        tokenNm,
				uri:         tokenURI,
				content:     content,
				description: "",
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "empty address",
			},
		},
	}

	for _, tc := range tests {
		func() {
			defer func() {
				if r := recover(); r != nil {
					err := r.(error)
					require.False(t, tc.errArgs.shouldPass)
					require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}()
			msg := types.NewMsgEditCertificate(tc.args.denomID, tc.args.tokenID, tc.args.name,
				tc.args.uri, tc.args.content, tc.args.description, tc.args.owner)
			err := msg.ValidateBasic()
			if tc.errArgs.shouldPass {
				require.NoError(t, err, tc.name)
			} else {
				require.Error(t, err, tc.name)
				require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
			}
		}()
	}
}

func TestMsgRevokeCertificate(t *testing.T) {
	type args struct {
		denomID     string
		tokenID     string
		description string
		revoker     sdk.AccAddress
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
			name: "MsgRevokeCertificate: Valid",
			args: args{
				denomID:     "certificateauditing",
				tokenID:     tokenID,
				description: "",
				revoker:     certifier,
			},
			errArgs: errArgs{
				shouldPass: true,
				contains:   "",
			},
		},
		{
			name: "MsgRevokeCertificate: Missing Revoker",
			args: args{
				denomID:     "certificateauditing",
				tokenID:     tokenID,
				description: "",
			},
			errArgs: errArgs{
				shouldPass: false,
				contains:   "empty address",
			},
		},
	}

	for _, tc := range tests {
		func() {
			defer func() {
				if r := recover(); r != nil {
					err := r.(error)
					require.False(t, tc.errArgs.shouldPass)
					require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
				}
			}()
			msg := types.NewMsgRevokeCertificate(tc.args.denomID, tc.args.tokenID,
				tc.args.description, tc.args.revoker)
			err := msg.ValidateBasic()
			if tc.errArgs.shouldPass {
				require.NoError(t, err, tc.name)
			} else {
				require.Error(t, err, tc.name)
				require.True(t, strings.Contains(err.Error(), tc.errArgs.contains))
			}
		}()
	}
}
