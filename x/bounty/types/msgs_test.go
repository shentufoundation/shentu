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
		{"", "name", "desc", addrs[0], false},
		{"1", "", "desc", addrs[0], false},
		{"1", "name", "", addrs[0], false},
		{"1", "name", "desc", sdk.AccAddress{}, false},
	}

	for i, test := range tests {
		msg := NewMsgCreateProgram(test.pid, test.name, test.detail, test.creatorAddress)
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

func TestMsgEditProgram(t *testing.T) {
	tests := []struct {
		pid            string
		name           string
		detail         string
		creatorAddress sdk.AccAddress
		expectPass     bool
	}{
		{"1", "name", "desc", addrs[0], true},
		{"1", "", "desc", addrs[0], true},
		{"1", "name", "", addrs[0], true},
		{"", "name", "desc", addrs[0], false},
		{"1", "name", "desc", sdk.AccAddress{}, false},
	}

	for i, test := range tests {
		msg := NewMsgEditProgram(test.pid, test.name, test.detail, test.creatorAddress)
		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgEditProgram)

		if test.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{test.creatorAddress})
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
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
		{"", addrs[0], false},
		{"1", sdk.AccAddress{}, false},
	}

	for i, tc := range testCases {
		msg := NewMsgActivateProgram(tc.pid, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgActivateProgram)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
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
		{"", addrs[0], false},
		{"1", sdk.AccAddress{}, false},
		{"", sdk.AccAddress{}, false},
	}

	for i, tc := range testCases {
		msg := NewMsgCloseProgram(tc.pid, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgCloseProgram)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgSubmitFinding(t *testing.T) {
	testCases := []struct {
		pid, fid, hash string
		addr           sdk.AccAddress
		severityLevel  SeverityLevel
		expectPass     bool
	}{
		{"1", "1", "hash", addrs[0], 3, true},
		{"", "1", "hash", addrs[0], 3, false},
		{"1", "", "hash", addrs[0], 3, false},
		{"1", "1", "", addrs[0], 3, false},
		{"1", "1", "hash", sdk.AccAddress{}, 3, false},
		{"1", "1", "hash", addrs[0], 10, false},
	}

	for i, tc := range testCases {
		msg := NewMsgSubmitFinding(tc.pid, tc.fid, tc.hash, tc.addr, tc.severityLevel)
		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgSubmitFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgActivateFinding(t *testing.T) {
	testCases := []struct {
		fid        string
		addr       sdk.AccAddress
		expectPass bool
	}{
		{"1", addrs[0], true},
		{"1", sdk.AccAddress{}, false},
		{"", addrs[0], false},
	}

	for i, tc := range testCases {
		msg := NewMsgActivateFinding(tc.fid, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgActivateFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgConfirmFinding(t *testing.T) {
	testCases := []struct {
		fid, fingerprint string
		addr             sdk.AccAddress
		expectPass       bool
	}{
		{"1", "fingerprint", addrs[0], true},
		{"", "fingerprint", addrs[0], false},
		{"1", "", addrs[0], false},
		{"1", "fingerprint", sdk.AccAddress{}, false},
	}

	for i, tc := range testCases {
		msg := NewMsgConfirmFinding(tc.fid, tc.fingerprint, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgConfirmFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgConfirmFindingPaid(t *testing.T) {
	testCases := []struct {
		fid        string
		addr       sdk.AccAddress
		expectPass bool
	}{
		{"1", addrs[0], true},
		{"1", sdk.AccAddress{}, false},
		{"", addrs[0], false},
	}

	for _, tc := range testCases {
		msg := NewMsgConfirmFindingPaid(tc.fid, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgConfirmFindingPaid)

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
		{"1", sdk.AccAddress{}, false},
		{"", addrs[0], false},
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

func TestMsgPublishFinding(t *testing.T) {
	testCases := []struct {
		fid, desc, poc string
		addr           sdk.AccAddress
		expectPass     bool
	}{
		{"1", "desc", "poc", addrs[0], true},
		{"", "desc", "poc", addrs[0], false},
		{"1", "", "poc", addrs[0], false},
		{"1", "desc", "", addrs[0], false},
		{"1", "desc", "poc", sdk.AccAddress{}, false},
	}

	for _, tc := range testCases {
		msg := NewMsgPublishFinding(tc.fid, tc.desc, tc.poc, tc.addr)

		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgPublishFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.addr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}
