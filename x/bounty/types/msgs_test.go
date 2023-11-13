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
		pid            string
		name           string
		detail         string
		creatorAddress sdk.AccAddress
		expectPass     bool
	}{
		{"1", "name", "desc", addrs[0], true},
		{"1", "name", "desc", sdk.AccAddress{}, false},
	}

	for i, test := range tests {
		msg := NewMsgCreateProgram(test.pid, test.name, test.detail, test.creatorAddress, nil)
		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgCreateProgram)

		if test.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{test.creatorAddress})
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgSubmitFinding(t *testing.T) {
	testCases := []struct {
		pid, fid, title, desc, hash string
		addr                        sdk.AccAddress
		severityLevel               int8
		expectPass                  bool
	}{
		{"1", "1", "title", "desc", "hash", addrs[0], 3, true},
		{"", "1", "title", "desc", "hash", addrs[0], 3, false},
		{"2", "2", "title", "desc", "hash", sdk.AccAddress{}, 3, false},
	}

	for _, tc := range testCases {
		msg := NewMsgSubmitFinding(tc.pid, tc.fid, tc.title, tc.desc, tc.hash, tc.addr, SeverityLevel(tc.severityLevel))
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

func TestMsgActivateProgram(t *testing.T) {
	testCases := []struct {
		pid        string
		addr       sdk.AccAddress
		expectPass bool
	}{
		{"1", addrs[0], true},
		{"2", sdk.AccAddress{}, false},
	}

	for _, tc := range testCases {
		msg := NewMsgActivateProgram(tc.pid, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgActivateProgram)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}

func TestMsgCloseProgram(t *testing.T) {
	testCases := []struct {
		pid        string
		addr       sdk.AccAddress
		expectPass bool
	}{
		{"1", addrs[0], true},
		{"2", sdk.AccAddress{}, false},
	}

	for _, tc := range testCases {
		msg := NewMsgCloseProgram(tc.pid, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgCloseProgram)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}

func TestMsgConfirmFinding(t *testing.T) {
	testCases := []struct {
		fid        string
		addr       sdk.AccAddress
		expectPass bool
	}{
		{"1", addrs[0], true},
		{"2", sdk.AccAddress{}, false},
	}

	for _, tc := range testCases {
		msg := NewMsgConfirmFinding(tc.fid, "", tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgConfirmFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}

func TestMsgCloseFinding(t *testing.T) {
	testCases := []struct {
		fid        string
		addr       sdk.AccAddress
		expectPass bool
	}{
		{"1", addrs[0], true},
		{"2", sdk.AccAddress{}, false},
	}

	for _, tc := range testCases {
		msg := NewMsgCloseFinding(tc.fid, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgCloseFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}

func TestMsgReleaseFinding(t *testing.T) {
	testCases := []struct {
		fid, desc, poc string
		addr           sdk.AccAddress
		expectPass     bool
	}{
		{"1", "desc", "poc", addrs[0], true},
		{"2", "desc", "poc", sdk.AccAddress{}, false},
	}

	for _, tc := range testCases {
		msg := NewMsgReleaseFinding(tc.fid, tc.desc, tc.poc, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgReleaseFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}
