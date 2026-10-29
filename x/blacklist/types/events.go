package types

// Event types emitted by the x/blacklist module.
const (
	// EventTypeBanAddress is emitted when an address is added to the
	// blacklist.
	EventTypeBanAddress = "ban_address"

	// EventTypeUnbanAddress is emitted when an address is removed from the
	// blacklist.
	EventTypeUnbanAddress = "unban_address"
)

// Event attribute keys.
const (
	AttributeKeyAddress   = "address"
	AttributeKeyCategory  = "category"
	AttributeKeyExpiresAt = "expires_at"
)
