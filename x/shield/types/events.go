package types

const (
	EventTypeCreateReimbursement = "create_reimbursement"

	AttributeKeyShield              = "shield"
	AttributeKeyDeposit             = "deposit"
	AttributeKeySponsor             = "sponsor"
	AttributeKeySponsorAddress      = "sponsor_address"
	AttributeKeyPoolID              = "pool_id"
	AttributeKeyAdditionalTime      = "additional_time"
	AttributeKeyTimeOfCoverage      = "time_of_coverage"
	AttributeKeyBlocksOfCoverage    = "blocks_of_coverage"
	AttributeKeyCollateral          = "collateral"
	AttributeKeyDenom               = "denom"
	AttributeKeyToAddr              = "to_address"
	AttributeKeyAccountAddress      = "account_address"
	AttributeKeyAmount              = "amount"
	AttributeKeyPurchaseID          = "purchase_id"
	AttributeKeyCompensationAmount  = "compensation_amount"
	AttributeKeyBeneficiary         = "beneficiary"
	AttributeKeyPurchaseDescription = "purchase_description"
	AttributeKeyNativeServiceFee    = "native_service_fee"
	AttributeKeyForeignServiceFee   = "foreign_service_fee"
	AttributeKeyProtectionEndTime   = "protection_end_time"
	AttributeValueCategory          = ModuleName
)
