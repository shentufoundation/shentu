package types

const (
	// ModuleName defines the module name.
	ModuleName = "blacklist"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for the blacklist module.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName
)

const (
	// EntriesKeyPrefix is the prefix byte for blacklisted address entries.
	EntriesKeyPrefix byte = 0x01
)
