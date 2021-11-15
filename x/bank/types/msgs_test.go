package types

import (
	//"fmt"
	"testing"

	//"github.com/certikfoundation/shentu/v2/x/bank"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

)

func TestMsgSendRoute(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("from"))
	addr2 := sdk.AccAddress([]byte("to"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("uctk", 10))
	var msg = NewMsgLockedSend(addr1, addr2, "", coins)
	require.Equal(t, msg.Route(), bankTypes.RouterKey)
	require.Equal(t, msg.Type(), "locked_send")
}

