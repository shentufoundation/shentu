package gov_test

//
//import (
//	"bytes"
//	"log"
//	"sort"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
//	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
//	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
//	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
//	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
//	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
//)
//
//var (
//	valTokens           = sdk.TokensFromConsensusPower(1e3, sdk.DefaultPowerReduction)
//	TestProposal        = govtypesv1beta1.NewTextProposal("Test", "description")
//	TestDescription     = stakingtypes.NewDescription("T", "E", "S", "T", "Z")
//	TestCommissionRates = stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
//)
//
//// mkTestLegacyContent creates a MsgExecLegacyContent for testing purposes.
//func mkTestLegacyContent(t *testing.T) *govtypesv1.MsgExecLegacyContent {
//	msgContent, err := govtypesv1.NewLegacyContent(TestProposal, authtypes.NewModuleAddress(govtypes.ModuleName).String())
//	require.NoError(t, err)
//
//	return msgContent
//}
//
//// SortAddresses - Sorts Addresses
//func SortAddresses(addrs []sdk.AccAddress) {
//	byteAddrs := make([][]byte, len(addrs))
//
//	for i, addr := range addrs {
//		byteAddrs[i] = addr.Bytes()
//	}
//
//	SortByteArrays(byteAddrs)
//
//	for i, byteAddr := range byteAddrs {
//		addrs[i] = byteAddr
//	}
//}
//
//// implement `Interface` in sort package.
//type sortByteArrays [][]byte
//
//func (b sortByteArrays) Len() int {
//	return len(b)
//}
//
//func (b sortByteArrays) Less(i, j int) bool {
//	// bytes package already implements Comparable for []byte.
//	switch bytes.Compare(b[i], b[j]) {
//	case -1:
//		return true
//	case 0, 1:
//		return false
//	default:
//		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
//		return false
//	}
//}
//
//func (b sortByteArrays) Swap(i, j int) {
//	b[j], b[i] = b[i], b[j]
//}
//
//// SortByteArrays - sorts the provided byte array
//func SortByteArrays(src [][]byte) [][]byte {
//	sorted := sortByteArrays(src)
//	sort.Sort(sorted)
//	return sorted
//}
//
//var pubkeys = []cryptotypes.PubKey{
//	ed25519.GenPrivKey().PubKey(),
//	ed25519.GenPrivKey().PubKey(),
//	ed25519.GenPrivKey().PubKey(),
//}
//
//const contextKeyBadProposal = "contextKeyBadProposal"
//
//func createValidators(t *testing.T, stakingHandler sdk.Handler, ctx sdk.Context, addrs []sdk.ValAddress, powerAmt []int64) {
//	require.True(t, len(addrs) <= len(pubkeys), "Not enough pubkeys specified at top of file.")
//
//	for i := 0; i < len(addrs); i++ {
//		valTokens := sdk.TokensFromConsensusPower(powerAmt[i], sdk.DefaultPowerReduction)
//		valCreateMsg, err := stakingtypes.NewMsgCreateValidator(
//			addrs[i], pubkeys[i], sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
//			TestDescription, TestCommissionRates, sdk.OneInt(),
//		)
//		require.NoError(t, err)
//		handleAndCheck(t, stakingHandler, ctx, valCreateMsg)
//	}
//}
//
//func handleAndCheck(t *testing.T, h sdk.Handler, ctx sdk.Context, msg sdk.Msg) {
//	res, err := h(ctx, msg)
//	require.NoError(t, err)
//	require.NotNil(t, res)
//}
