package types

import (
	"bytes"
	"testing"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestMsgWithdrawRewardsValidateBasic(t *testing.T) {
	validAddress := sdk.AccAddress(bytes.Repeat([]byte{1}, 20)).String()

	tests := []struct {
		name          string
		from          string
		wantErr       bool
		wantEmptyErr  bool
		wantErrSubstr string
	}{
		{
			name: "valid sender",
			from: validAddress,
		},
		{
			name:         "empty sender",
			wantErr:      true,
			wantEmptyErr: true,
		},
		{
			name:          "invalid sender address",
			from:          "invalid-address",
			wantErr:       true,
			wantErrSubstr: "invalid sender address invalid-address",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := MsgWithdrawRewards{From: tc.from}

			var err error
			require.NotPanics(t, func() {
				err = msg.ValidateBasic()
			})

			if !tc.wantErr {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			if tc.wantEmptyErr {
				require.True(t, errorsmod.IsOf(err, ErrEmptySender))
			}
			if tc.wantErrSubstr != "" {
				require.ErrorContains(t, err, tc.wantErrSubstr)
			}
		})
	}
}
