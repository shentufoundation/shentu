package types

const (
	EventTypeCreateProgram  = "create_program"
	EventTypeEditProgram    = "edit_program"
	EventTypeSubmitFinding  = "submit_finding"
	EventTypeAcceptFinding  = "accept_finding"
	EventTypeRejectFinding  = "reject_finding"
	EventTypeCancelFinding  = "cancel_finding"
	EventTypeReleaseFinding = "release_finding"
	EventTypeOpenProgram    = "open_program"
	EventTypeEndProgram     = "end_program"

	AttributeKeyProgramID = "program_id"
	AttributeKeyDeposit   = "deposit"

	AttributeKeyFindingID = "finding_id"

	AttributeValueCategory = ModuleName
)
