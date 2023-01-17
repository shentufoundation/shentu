package types

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/shentufoundation/shentu/v2/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var addrs = []sdk.AccAddress{sdk.AccAddress("test1"), sdk.AccAddress("test2")}

func TestMsgCreateProgram(t *testing.T) {
	decKey, err := ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	require.NoError(t, err)
	encKey := crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)
	deposit := sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e5)))
	var sET, jET, cET time.Time

	tests := []struct {
		creatorAddress string
		description    string
		encKey         []byte
		commissionRate sdk.Dec
		deposit        sdk.Coins
		expectPass     bool
	}{
		{"Test Program", "test pass", encKey,
			sdk.ZeroDec(), deposit, true,
		},
		{"Test Program", "test fail, encKey is nil", nil,
			sdk.ZeroDec(), deposit, false,
		},
	}

	for i, test := range tests {
		msg, err := NewMsgCreateProgram(test.creatorAddress, test.description, test.encKey, test.commissionRate,
			test.deposit, sET, jET, cET)

		if test.expectPass {
			require.NoError(t, err)
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			//
			//require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgHostAcceptFinding(t *testing.T) {
	testCases := []struct {
		findingId  uint64
		hostAddr   sdk.AccAddress
		comment    string
		expectPass bool
	}{
		{0, addrs[0], "comment", false},
		{1, sdk.AccAddress{}, "comment", false},
		{1, addrs[0], "comment", true},
		{1, addrs[0], "", true},
	}

	for _, tc := range testCases {
		msg := NewMsgHostAcceptFinding(tc.findingId, tc.comment, tc.hostAddr)
		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgAcceptFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.hostAddr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}

func TestMsgHostRejectFinding(t *testing.T) {
	testCases := []struct {
		findingId  uint64
		hostAddr   sdk.AccAddress
		comment    string
		expectPass bool
	}{
		{0, addrs[0], "comment", false},
		{1, sdk.AccAddress{}, "comment", false},
		{1, addrs[0], "comment", true},
		{1, addrs[0], "", true},
	}

	for _, tc := range testCases {
		msg := NewMsgHostRejectFinding(tc.findingId, tc.comment, tc.hostAddr)
		require.Equal(t, msg.Route(), RouterKey)
		require.Equal(t, msg.Type(), TypeMsgRejectFinding)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic())
			require.Equal(t, msg.GetSigners(), []sdk.AccAddress{tc.hostAddr})
		} else {
			require.Error(t, msg.ValidateBasic())
		}
	}
}

func TestHostAcceptGetSignBytes(t *testing.T) {
	msg := NewMsgHostAcceptFinding(1, "comment", addrs[0])
	msg.GetSignBytes()

	//expected := ""
	//require.Equal(t, expected, string(res))

	msg1 := NewMsgHostRejectFinding(1, "comment", addrs[0])
	msg1.GetSignBytes()
}
