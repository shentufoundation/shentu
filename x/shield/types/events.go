package types

const (
	EventTypeCreatePool             = "create_pool"
	EventTypeUpdatePool             = "update_pool"
	EventTypePausePool              = "pause_pool"
	EventTypeResumePool             = "resume_pool"
	EventTypeDepositCollateral      = "deposit_collateral"
	EventTypeWithdrawCollateral     = "withdraw_collateral"
	EventTypePurchaseShield         = "purchase_shield"
	EventTypeWithdrawRewards        = "withdraw_rewards"
	EventTypeWithdrawForeignRewards = "withdraw_foreign_rewards"
	EventTypeClearPayouts           = "clear_payouts"
	EventTypeCreateReimbursement    = "create_reimbursement"
	EventTypeWithdrawReimbursement  = "withdraw_reimbursement"

	AttributeKeyShield             = "shield"
	AttributeKeyDeposit            = "deposit"
	AttributeKeySponsor            = "sponsor"
	AttributeKeyPoolID             = "pool_id"
	AttributeKeyCollateral         = "collateral"
	AttributeKeyDenom              = "denom"
	AttributeKeyToAddr             = "to_address"
	AttributeKeyAccountAddress     = "account_address"
	AttributeKeyAmount             = "amount"
	AttributeKeyPurchaseID         = "purchase_id"
	AttributeKeyCompensationAmount = "compensation_amount"
	AttributeKeyBeneficiary        = "beneficiary"
	AttributeValueCategory         = ModuleName
)
