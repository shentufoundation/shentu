package types

const (
	EventTypeCreateProgram     = "create_program"
	EventTypeSubmitFinding     = "submit_finding"
	EventTypeWithdrawalFinding = "withdrawal_finding"
	EventTypeReactivateFinding = "reactivate_finding"

	AttributeKeyProgramID = "program_id"
	AttributeKeyDeposit   = "deposit"

	AttributeKeyFindingID = "finding_id"

	AttributeValueCategory = ModuleName
)
