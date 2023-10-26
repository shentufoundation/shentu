package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	addrs = []sdk.AccAddress{sdk.AccAddress("test1"), sdk.AccAddress("test2")}
)

func TestMsgCreateProgram(t *testing.T) {
	tests := []struct {
		name           string
		description    string
		pid            string
		creatorAddress sdk.AccAddress
		expectPass     bool
	}{
		{"name", "desc", "1", addrs[0], true},
		{"name", "desc", "1", sdk.AccAddress{}, false},
	}

	for i, test := range tests {
		msg, err := NewMsgCreateProgram(test.name, test.description, test.pid, test.creatorAddress, nil, nil)
		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgCreateProgram)

		if test.expectPass {
			require.NoError(t, err)
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{test.creatorAddress})
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgSubmitFinding(t *testing.T) {
	testCases := []struct {
		pid, fid, title, desc string
		addr                  sdk.AccAddress
		severityLevel         int8
		expectPass            bool
	}{
		{"1", "1", "title", "desc", addrs[0], 3, true},
		{"2", "2", "title", "desc", sdk.AccAddress{}, 3, false},
	}

	for _, tc := range testCases {
		msg := NewMsgSubmitFinding(tc.pid, tc.fid, tc.title, tc.desc, tc.addr, SeverityLevel(tc.severityLevel))
		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgSubmitFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}
