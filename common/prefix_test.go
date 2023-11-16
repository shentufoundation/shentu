package common

import (
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func TestPrefixToShentu(t *testing.T) {
	sum := sha256.Sum256([]byte("hello world\n"))
	ss := "certik"

	address, err := bech32.ConvertAndEncode(ss, sum[:])
	require.NoError(t, err)

	shentuAddr, err := PrefixToShentu(address)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(shentuAddr, "shentu"))

	address, err = bech32.ConvertAndEncode("shentu", sum[:])
	require.NoError(t, err)
	_, err = PrefixToShentu(address)
	require.NoError(t, err)
}
