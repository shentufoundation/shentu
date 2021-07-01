package types

var (
	EventTypeCreateAdmin = "create_admin"
	EventTypeRevokeAdmin = "revoke_admin"

	AttributeKeyAdminCreator = "admin_creator"
	AttributeKeyAdminRevoker = "admin_revoker"
	AttributeKeyCreated      = "created"
	AttributeKeyRevoked      = "revoked"

	EventTypeIssueCertificate  = "issue_certificate"
	EventTypeEditCertificate   = "edit_certificate"
	EventTypeRevokeCertificate = "revoke_certificate"
)
