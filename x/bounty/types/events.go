package types

const (
	// Module name as value category
	AttributeValueCategory = ModuleName

	// Program related events
	EventTypeCreateProgram   = "create_program"
	EventTypeEditProgram     = "edit_program"
	EventTypeActivateProgram = "activate_program"
	EventTypeCloseProgram    = "close_program"

	// Finding related events
	EventTypeSubmitFinding          = "submit_finding"
	EventTypeEditFinding            = "edit_finding"
	EventTypeEditFindingPaymentHash = "edit_finding_payment_hash"
	EventTypeActivateFinding        = "activate_finding"
	EventTypeConfirmFinding         = "confirm_finding"
	EventTypeConfirmFindingPaid     = "confirm_finding_paid"
	EventTypeCloseFinding           = "close_finding"
	EventTypePublishFinding         = "publish_finding"

	// Program/Finding attributes
	AttributeKeyProgramID = "program_id"
	AttributeKeyFindingID = "finding_id"

	// Theorem related events
	EventTypeCreateTheorem = "create_theorem"
	EventTypeTheoremGrant  = "theorem_grant"
	EventTypeDeleteTheorem = "delete_theorem"
	// Proof related events
	EventTypeSubmitProofHash         = "submit_proof_hash"
	EventTypeSubmitProofDetail       = "submit_proof_detail"
	EventTypeSubmitProofVerification = "submit_proof_verification"
	EventTypeProofDeposit            = "proof_deposit"
	EventTypeDeleteProof             = "delete_proof"
	EventTypeWithdrawReward          = "withdraw_reward"

	// Theorem/Proof attributes
	AttributeKeyTheoremID                = "theorem_id"
	AttributeKeyProofID                  = "proof_id"
	AttributeKeyTheoremProposer          = "proposer"
	AttributeKeyTheoremGrantor           = "grantor"
	AttributeKeyProofDepositor           = "depositor"
	AttributeKeyTheoremResult            = "theorem_result"
	AttributeValueTheoremDropped         = "theorem_dropped"
	AttributeKeyTheoremProofPeriodStart  = "theorem_proof_period_start"
	AttributeKeyProofHashLockPeriodStart = "proof_hash_lock_period_start"
	AttributeKeyChecker                  = "checker"
	AttributeKeyProofStatus              = "proof_status"
	AttributeKeyReward                   = "reward"
	AttributeKeyAddress                  = "address"
)
