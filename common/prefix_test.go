package common

import (
	"crypto/sha256"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

func TestPrefixToCertik(t *testing.T) {
	sum := sha256.Sum256([]byte("hello world\n"))
	ss := "shentu"

	address, err := bech32.ConvertAndEncode(ss, sum[:])
	require.NoError(t, err)

	certikAddr, err := PrefixToCertik(address)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(certikAddr, "certik"))

	address, err = bech32.ConvertAndEncode("certik", sum[:])
	require.NoError(t, err)
	_, err = PrefixToCertik(address)
	require.Error(t, err)
}

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
	require.Error(t, err)
}

func TestAddress(t *testing.T) {
	address := "certik1fdqkh6v4d09hq7hu8xq8gl9zr2zm0azydcpq0d"
	shentuAddr, err := PrefixToShentu(address)
	if err != nil {
		t.Log(err.Error())
	}
	t.Log("address:", shentuAddr)
}
