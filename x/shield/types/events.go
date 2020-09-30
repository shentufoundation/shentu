package types

const (
	EventTypeCreatePool             = "create_pool"
	EventTypeUpdatePool             = "update_pool"
	EventTypePausePool              = "pause_pool"
	EventTypeResumePool             = "resume_pool"
	EventTypeDepositCollateral      = "deposit_collateral"
	EventTypePurchase               = "purchase_shield"
	EventTypeWithdrawForeignRewards = "withdraw_foreign_rewards"
	EventTypeClearPayouts           = "clear_payouts"

	AttributeKeyShield           = "shield"
	AttributeKeyDeposit          = "deposit"
	AttributeKeySponsor          = "sponsor"
	AttributeKeyPoolID           = "pool_id"
	AttributeKeyAdditionalTime   = "additional_time"
	AttributeKeyTimeOfCoverage   = "time_of_coverage"
	AttributeKeyBlocksOfCoverage = "blocks_of_coverage"
	AttributeKeyCollateral       = "collateral"
	AttributeKeyDenom            = "denom"
	AttributeKeyToAddr           = "to_addr"
	AttributeKeyAmount           = "amount"
	AttributeValueCategory       = ModuleName
)
