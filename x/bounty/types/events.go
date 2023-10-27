package types

const (
	EventTypeCreateProgram = "create_program"
	EventTypeEditProgram   = "edit_program"
	EventTypeOpenProgram   = "open_program"
	EventTypeCloseProgram  = "close_program"

	EventTypeSubmitFinding  = "submit_finding"
	EventTypeEditFinding    = "edit_finding"
	EventTypeRejectFinding  = "reject_finding"
	EventTypeAcceptFinding  = "accept_finding"
	EventTypeCloseFinding   = "close_finding"
	EventTypeReleaseFinding = "release_finding"

	AttributeKeyProgramID = "program_id"
	AttributeKeyFindingID = "finding_id"

	AttributeValueCategory = ModuleName
)
