package v1

import (
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
)

// IsCertifierUpdateProposalMsg reports whether msg is either a
// MsgUpdateCertifier or a legacy CertifierUpdateProposal wrapped in
// MsgExecLegacyContent. Kept in the types package so both genesis
// validation and the keeper's submission guard share one definition.
func IsCertifierUpdateProposalMsg(msg sdk.Msg) (bool, error) {
	if sdk.MsgTypeURL(msg) == sdk.MsgTypeURL(&certtypes.MsgUpdateCertifier{}) {
		return true, nil
	}
	legacyMsg, ok := msg.(*govtypesv1.MsgExecLegacyContent)
	if !ok {
		return false, nil
	}
	content, err := govtypesv1.LegacyContentFromMessage(legacyMsg)
	if err != nil {
		return false, err
	}
	//nolint:staticcheck // legacy proposal path retained for in-flight proposals
	_, isCertUpdate := content.(*certtypes.CertifierUpdateProposal)
	return isCertUpdate, nil
}

// ValidateCertifierUpdateSoloMessage rejects proposals that bundle a
// certifier-update message with any other message. CertifierUpdate is
// the only proposal type that traverses the certifier round; bundling
// it with ordinary messages would let those messages ride the
// head-count tally and bypass validator stake voting entirely.
func ValidateCertifierUpdateSoloMessage(messages []sdk.Msg) error {
	if len(messages) <= 1 {
		return nil
	}
	for _, m := range messages {
		isCertUpdate, err := IsCertifierUpdateProposalMsg(m)
		if err != nil {
			return err
		}
		if isCertUpdate {
			return errors.Wrap(
				sdkerrors.ErrInvalidRequest,
				"proposals containing a certifier-update message must contain exactly one message",
			)
		}
	}
	return nil
}
