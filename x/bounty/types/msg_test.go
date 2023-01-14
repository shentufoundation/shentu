package types

import (
	"crypto/rand"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/shentufoundation/shentu/v2/common"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
