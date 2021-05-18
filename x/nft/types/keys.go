package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	AdminKeyPrefix = []byte{0x10}
)

func AdminKey(addr sdk.AccAddress) []byte {
	return append(AdminKeyPrefix, addr...)
}
