package types

const (
	EventTypeCreateProgram   = "create_program"
	EventTypeEditProgram     = "edit_program"
	EventTypeActivateProgram = "activate_program"
	EventTypeCloseProgram    = "close_program"

	EventTypeSubmitFinding          = "submit_finding"
	EventTypeEditFinding            = "edit_finding"
	EventTypeEditFindingPaymentHash = "edit_finding_payment_hash"
	EventTypeActivateFinding        = "activate_finding"
	EventTypeConfirmFinding         = "confirm_finding"
	EventTypeConfirmFindingPaid     = "confirm_finding_paid"
	EventTypeCloseFinding           = "close_finding"
	EventTypePublishFinding         = "publish_finding"

	EventTypeCreateTheorem = "create_theorem"

	AttributeKeyProgramID        = "program_id"
	AttributeKeyFindingID        = "finding_id"
	AttributeKeyTheoremID        = "theorem_id"
	AttributeKeyProofPeriodStart = "proof_period_start"

	AttributeValueCategory = ModuleName
)
