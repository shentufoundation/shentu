package v260

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	proto "github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// Proposal types implement UnpackInterfaceMessages to unpack
// Content fields.
var _ types.UnpackInterfacesMessage = Proposal{}
var _ types.UnpackInterfacesMessage = Proposals{}

// NewProposal creates a new Proposal instance
func NewProposal(content govtypes.Content, id uint64, proposerAddress sdk.AccAddress, isProposerCouncilMember bool, submitTime time.Time, depositEndTime time.Time) (Proposal, error) {
	p := Proposal{
		ProposalId:              id,
		Status:                  StatusDepositPeriod,
		IsProposerCouncilMember: isProposerCouncilMember,
		ProposerAddress:         proposerAddress.String(),
		FinalTallyResult:        govtypes.EmptyTallyResult(),
		TotalDeposit:            sdk.NewCoins(),
		SubmitTime:              submitTime,
		DepositEndTime:          depositEndTime,
	}

	msg, ok := content.(proto.Message)
	if !ok {
		return Proposal{}, fmt.Errorf("%T does not implement proto.Message", content)
	}

	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return Proposal{}, err
	}

	p.Content = any

	return p, nil
}

// GetContent returns the proposal Content
func (p Proposal) GetContent() govtypes.Content {
	content, ok := p.Content.GetCachedValue().(govtypes.Content)
	if !ok {
		return nil
	}
	return content
}

// String implements stringer interface
func (p Proposal) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (p Proposal) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	var content govtypes.Content
	return unpacker.UnpackAny(p.Content, &content)
}

func (p Proposal) GetTitle() string {
	content := p.GetContent()
	if content == nil {
		return ""
	}
	return content.GetTitle()
}

func (p Proposal) ProposalType() string {
	content := p.GetContent()
	if content == nil {
		return ""
	}
	return content.ProposalType()
}

func (p Proposal) ProposalRoute() string {
	content := p.GetContent()
	if content == nil {
		return ""
	}
	return content.ProposalRoute()
}

// Proposals is an array of proposals.
type Proposals []Proposal

// String implements stringer interface.
func (p Proposals) String() string {
	out := "ID - (Status) [Type] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] %s\n",
			prop.ProposalId, prop.Status,
			prop.ProposalType(), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (p Proposals) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	for _, x := range p {
		err := x.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}

type (
	// ProposalQueue is a type alias that represents a list of proposal IDs.
	ProposalQueue []uint64
)

// ProposalStatusFromString turns a string into a ProposalStatus.
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	num, ok := ProposalStatus_value[str]
	if !ok {
		return StatusNil, fmt.Errorf("'%s' is not a valid proposal status", str)
	}
	return ProposalStatus(num), nil
}

// ValidProposalStatus returns true if the proposal status is valid and false otherwise.
func ValidProposalStatus(status ProposalStatus) bool {
	if status == StatusDepositPeriod ||
		status == StatusCertifierVotingPeriod ||
		status == StatusPassed ||
		status == StatusRejected ||
		status == StatusFailed ||
		status == StatusValidatorVotingPeriod {
		return true
	}
	return false
}

// Marshal implements the Marshal method for protobuf compatibility.
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal implements the Unmarshal method for protobuf compatibility.
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// MarshalJSON marshals to JSON using string.
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}

	*status = bz2
	return nil
}

// Format implements the fmt.Formatter interface.
func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

// ProposalHandler implements the Handler interface for governance module-based
// proposals (ie. TextProposal and SoftwareUpgradeProposal). Since these are
// merely signaling mechanisms at the moment and do not affect state, it
// performs a no-op.
func ProposalHandler(_ sdk.Context, c govtypes.Content) error {
	switch c.ProposalType() {
	case govtypes.ProposalTypeText:
		// both proposal types do not change state so this performs a no-op
		return nil

	default:
		errMsg := fmt.Sprintf("unrecognized gov proposal type: %s", c.ProposalType())
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
