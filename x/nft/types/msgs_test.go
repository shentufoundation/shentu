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
	acc1 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String()
	acc2 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes()).String()
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
