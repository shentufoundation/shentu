package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/certikfoundation/shentu/x/cert"
	"github.com/certikfoundation/shentu/x/shield"
)

// Proposal defines a struct used by the governance module to allow for voting
// on network changes.
type Proposal struct {
	// proposal content interface
	types.Content `json:"content" yaml:"content"`

	// ID of the proposal
	ProposalID uint64 `json:"id" yaml:"id"`
	// status of the Proposal {Pending, Active, Passed, Rejected}
	Status ProposalStatus `json:"proposal_status" yaml:"proposal_status"`
	// whether or not the proposer is a council member (validator or certifier)
	IsProposerCouncilMember bool `json:"is_proposer_council_member" yaml:"is_proposer_council_member"`
	// proposer address
	ProposerAddress sdk.AccAddress `json:"proposer_address" yaml:"proposer_address"`
	// result of Tally
	FinalTallyResult gov.TallyResult `json:"final_tally_result" yaml:"final_tally_result"`

	// time of the block where TxGovSubmitProposal was included
	SubmitTime time.Time `json:"submit_time" yaml:"submit_time"`
	// time that the Proposal would expire if deposit amount isn't met
	DepositEndTime time.Time `json:"deposit_end_time" yaml:"deposit_end_time"`
	// current deposit on this proposal
	TotalDeposit sdk.Coins `json:"total_deposit" yaml:"total_deposit"`

	// VotingStartTime is the time of the block where MinDeposit was reached.
	// It is set to -1 if MinDeposit is not reached.
	VotingStartTime time.Time `json:"voting_start_time" yaml:"voting_start_time"`
	// time that the VotingPeriodString for this proposal will end and votes will be tallied
	VotingEndTime time.Time `json:"voting_end_time" yaml:"voting_end_time"`
}

// NewProposal returns a new proposal.
func NewProposal(content types.Content,
	id uint64,
	proposerAddress sdk.AccAddress,
	isProposerCouncilMember bool,
	submitTime time.Time,
	depositEndTime time.Time) Proposal {
	return Proposal{
		Content:                 content,
		ProposalID:              id,
		Status:                  StatusDepositPeriod,
		IsProposerCouncilMember: isProposerCouncilMember,
		ProposerAddress:         proposerAddress,
		FinalTallyResult:        gov.EmptyTallyResult(),
		TotalDeposit:            sdk.NewCoins(),
		SubmitTime:              submitTime,
		DepositEndTime:          depositEndTime,
	}
}

// String returns proposal struct in string type.
func (p Proposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Status:             %s
  Is Council Member:  %t
  Proposer Address:   %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s
  Description:        %s`,
		p.ProposalID, p.GetTitle(), p.ProposalType(), p.Status, p.IsProposerCouncilMember, p.ProposerAddress, p.SubmitTime, p.DepositEndTime,
		p.TotalDeposit, p.VotingStartTime, p.VotingEndTime, p.GetDescription(),
	)
}

// HasSecurityVoting returns true if the proposal needs to go through security
// (certifier) voting before stake (validator) voting.
func (p Proposal) HasSecurityVoting() bool {
	switch p.Content.(type) {
	case upgrade.SoftwareUpgradeProposal, cert.CertifierUpdateProposal, shield.ClaimProposal:
		return true
	default:
		return false
	}
}

// Proposals is an array of proposals.
type Proposals []Proposal

// String implements stringer interface.
func (p Proposals) String() string {
	out := "ID - (Status) [Type] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] %s\n",
			prop.ProposalID, prop.Status,
			prop.ProposalType(), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

type (
	// ProposalQueue is a type alias that represents a list of proposal IDs.
	ProposalQueue []uint64

	// ProposalStatus is a type alias that represents a proposal status as a byte
	ProposalStatus byte
)

const (
	// StatusNil is the nil status.
	StatusNil ProposalStatus = iota
	// StatusDepositPeriod is the deposit period status.
	StatusDepositPeriod
	// StatusCertifierVotingPeriod is the certifier voting period status.
	StatusCertifierVotingPeriod
	// StatusValidatorVotingPeriod is the validator voting period status.
	StatusValidatorVotingPeriod
	// StatusPassed is the passed status.
	StatusPassed
	// StatusRejected is the rejected status.
	StatusRejected
	// StatusFailed is the failed status.
	StatusFailed

	// DepositPeriod is the string of DepositPeriod.
	DepositPeriod = "DepositPeriod"
	// CertifierVotingPeriod is the string of CertifierVotingPeriod.
	CertifierVotingPeriod = "CertifierVotingPeriod"
	// ValidatorVotingPeriod is the string of ValidatorVotingPeriod.
	ValidatorVotingPeriod = "ValidatorVotingPeriod"
	// Passed is the string of Passed.
	Passed = "Passed"
	// Rejected is the string of Rejected.
	Rejected = "Rejected"
	// Failed is the string of Failed.
	Failed = "Failed"
)

// ProposalStatusFromString turns a string into a ProposalStatus.
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case DepositPeriod:
		return StatusDepositPeriod, nil

	case CertifierVotingPeriod:
		return StatusCertifierVotingPeriod, nil

	case Passed:
		return StatusPassed, nil

	case Rejected:
		return StatusRejected, nil

	case Failed:
		return StatusFailed, nil

	case "":
		return StatusNil, nil

	case ValidatorVotingPeriod:
		return StatusValidatorVotingPeriod, nil

	default:
		return ProposalStatus(0xff), fmt.Errorf("'%s' is not a valid proposal status", str)
	}
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

// String implements the Stringer interface.
func (status ProposalStatus) String() string {
	switch status {
	case StatusDepositPeriod:
		return DepositPeriod

	case StatusCertifierVotingPeriod:
		return CertifierVotingPeriod

	case StatusPassed:
		return Passed

	case StatusRejected:
		return Rejected

	case StatusFailed:
		return Failed

	case StatusValidatorVotingPeriod:
		return ValidatorVotingPeriod

	default:
		return ""
	}
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

// Proposal types
const (
	ProposalTypeText            string = "Text"
	ProposalTypeSoftwareUpgrade string = "SoftwareUpgrade"
)

// ProposalHandler implements the Handler interface for governance module-based
// proposals (ie. TextProposal and SoftwareUpgradeProposal). Since these are
// merely signaling mechanisms at the moment and do not affect state, it
// performs a no-op.
func ProposalHandler(_ sdk.Context, c types.Content) error {
	switch c.ProposalType() {
	case ProposalTypeText, ProposalTypeSoftwareUpgrade:
		// both proposal types do not change state so this performs a no-op
		return nil

	default:
		errMsg := fmt.Sprintf("unrecognized gov proposal type: %s", c.ProposalType())
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
