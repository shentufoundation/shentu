package types

var (
	EventTypeCreateAdmin       = "create_admin"
	EventTypeRevokeAdmin       = "revoke_admin"
	EventTypeIssueCertificate  = "issue_certificate"
	EventTypeEditCertificate   = "edit_certificate"
	EventTypeRevokeCertificate = "revoke_certificate"

	AttributeKeyAdminCreator = "admin_creator"
	AttributeKeyAdminRevoker = "admin_revoker"
	AttributeKeyCreated      = "created"
	AttributeKeyRevoked      = "revoked"
)
