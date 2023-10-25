package types

const (
	EventTypeCreateProgram = "create_program"
	EventTypeEditProgram   = "edit_program"
	EventTypeOpenProgram   = "open_program"
	EventTypeEndProgram    = "end_program"

	EventTypeSubmitFinding  = "submit_finding"
	EventTypeRejectFinding  = "reject_finding"
	EventTypeAcceptFinding  = "accept_finding"
	EventTypeCancelFinding  = "cancel_finding"
	EventTypeReleaseFinding = "release_finding"

	AttributeKeyProgramID = "program_id"
	AttributeKeyFindingID = "finding_id"

	AttributeValueCategory = ModuleName
)
