package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	AdminPrefix = []byte{0x06}
)

func AdminKey(addr sdk.AccAddress) []byte{
	return append(AdminPrefix, addr...)
}