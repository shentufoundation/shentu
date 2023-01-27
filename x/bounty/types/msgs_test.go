package types

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/common"
)

var (
	addrs         = []sdk.AccAddress{sdk.AccAddress("test1"), sdk.AccAddress("test2")}
	decKey, _     = ecies.GenerateKey(rand.Reader, ecies.DefaultCurve, nil)
	encKey        = crypto.FromECDSAPub(&decKey.ExportECDSA().PublicKey)
	deposit       = sdk.NewCoins(sdk.NewCoin(common.MicroCTKDenom, sdk.NewInt(1e5)))
	sET, jET, cET time.Time
)

func TestMsgCreateProgram(t *testing.T) {
	tests := []struct {
		creatorAddress sdk.AccAddress
		description    string
		encKey         []byte
		commissionRate sdk.Dec
		deposit        sdk.Coins
		expectPass     bool
	}{
		{addrs[0], "desc", encKey,
			sdk.ZeroDec(), deposit, true,
		},
		{sdk.AccAddress{}, "desc", encKey,
			sdk.ZeroDec(), deposit, false,
		},
	}

	for i, test := range tests {
		msg, err := NewMsgCreateProgram(test.creatorAddress.String(), test.description, test.encKey, test.commissionRate,
			test.deposit, sET, jET, cET)
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
		pid           uint64
		severityLevel int32
		addr          sdk.AccAddress
		title         string
		desc, poc     *codectypes.Any
		expectPass    bool
	}{
		{0, 0, addrs[0], "title", nil, nil, false},
		{1, 0, sdk.AccAddress{}, "title", nil, nil, false},
		{1, 0, addrs[0], "title", nil, nil, true},
	}

	for _, tc := range testCases {
		msg := NewMsgSubmitFinding(tc.addr.String(), tc.title, tc.desc, tc.poc, tc.pid, tc.severityLevel)
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

func TestMsgHostAcceptFinding(t *testing.T) {
	encComment := EciesEncryptedComment{
		EncryptedComment: []byte("comment"),
	}
	commentAny, _ := codectypes.NewAnyWithValue(&encComment)

	testCases := []struct {
		findingId  uint64
		hostAddr   sdk.AccAddress
		comment    *codectypes.Any
		expectPass bool
	}{
		{0, addrs[0], commentAny, false},
		{1, sdk.AccAddress{}, commentAny, false},
		{1, addrs[0], commentAny, true},
		{1, addrs[0], commentAny, true},
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
	encComment := EciesEncryptedComment{
		EncryptedComment: []byte("comment"),
	}
	commentAny, _ := codectypes.NewAnyWithValue(&encComment)

	testCases := []struct {
		findingId  uint64
		hostAddr   sdk.AccAddress
		comment    *codectypes.Any
		expectPass bool
	}{
		{0, addrs[0], commentAny, false},
		{1, sdk.AccAddress{}, commentAny, false},
		{1, addrs[0], commentAny, true},
		{1, addrs[0], nil, true},
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