package types

const (
	EventTypeCreateProgram   = "create_program"
	EventTypeEditProgram     = "edit_program"
	EventTypeActivateProgram = "activate_program"
	EventTypeCloseProgram    = "close_program"

	EventTypeSubmitFinding          = "submit_finding"
	EventTypeEditFinding            = "edit_finding"
	EventTypeEditFindingPaymentHash = "edit_finding_payment_hash"
	EventTypeConfirmFinding         = "confirm_finding"
	EventTypeConfirmFindingPaid     = "confirm_finding_paid"
	EventTypeCloseFinding           = "close_finding"
	EventTypeReleaseFinding         = "release_finding"

	AttributeKeyProgramID = "program_id"
	AttributeKeyFindingID = "finding_id"

	AttributeValueCategory = ModuleName
)
