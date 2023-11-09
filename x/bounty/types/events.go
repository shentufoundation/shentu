package types

const (
	EventTypeCreateProgram = "create_program"
	EventTypeEditProgram   = "edit_program"
	//todo
	EventTypeActivateProgram = "activate_program"
	EventTypeCloseProgram    = "close_program"

	EventTypeSubmitFinding  = "submit_finding"
	EventTypeEditFinding    = "edit_finding"
	EventTypeConfirmFinding = "confirm_finding"
	EventTypeCloseFinding   = "close_finding"
	EventTypeReleaseFinding = "release_finding"

	AttributeKeyProgramID = "program_id"
	AttributeKeyFindingID = "finding_id"

	AttributeValueCategory = ModuleName
)
