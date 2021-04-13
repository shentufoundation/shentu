package migrate

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	shieldtypes "github.com/certikfoundation/shentu/x/shield/types"

	certtypes "github.com/certikfoundation/shentu/x/cert/types"

	proto "github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v036distr "github.com/cosmos/cosmos-sdk/x/distribution/legacy/v036"
	v040distr "github.com/cosmos/cosmos-sdk/x/distribution/types"
	v034gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v034"
	v036gov "github.com/cosmos/cosmos-sdk/x/gov/legacy/v036"
	cosmosgov "github.com/cosmos/cosmos-sdk/x/gov/types"
	v040gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	v036params "github.com/cosmos/cosmos-sdk/x/params/legacy/v036"
	v040params "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	v038upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/legacy/v038"
	v040upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	govtypes "github.com/certikfoundation/shentu/x/gov/types"
)

type Deposit struct {
	v034gov.Deposit
	TxHash string `json:"txhash" yaml:"txhash"`
}

type Deposits []Deposit

// Vote wraps a vote and corresponding txhash.
type Vote struct {
	v034gov.Vote
	TxHash string `json:"txhash" yaml:"txhash"`
}

type Votes []Vote

type ProposalStatus byte

// Proposal defines a struct used by the governance module to allow for voting
// on network changes.
type Proposal struct {
	// proposal content interface
	v036gov.Content `json:"content" yaml:"content"`

	// ID of the proposal
	ProposalID uint64 `json:"id" yaml:"id"`
	// status of the Proposal {Pending, Active, Passed, Rejected}
	Status string `json:"proposal_status" yaml:"proposal_status"`
	// whether or not the proposer is a council member (validator or certifier)
	IsProposerCouncilMember bool `json:"is_proposer_council_member" yaml:"is_proposer_council_member"`
	// proposer address
	ProposerAddress sdk.AccAddress `json:"proposer_address" yaml:"proposer_address"`
	// result of Tally
	FinalTallyResult v034gov.TallyResult `json:"final_tally_result" yaml:"final_tally_result"`

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

type Proposals []Proposal

// DepositParams struct around deposits for governance
type DepositParams struct {
	// Minimum initial deposit when users submitting a proposal
	MinInitialDeposit sdk.Coins `json:"min_initial_deposit,omitempty" yaml:"min_initial_deposit,omitempty"`
	// Minimum deposit for a proposal to enter voting period.
	MinDeposit sdk.Coins `json:"min_deposit,omitempty" yaml:"min_deposit,omitempty"`
	// Maximum period for CTK holders to deposit on a proposal. Initial value: 2 months
	MaxDepositPeriod time.Duration `json:"max_deposit_period,omitempty" yaml:"max_deposit_period,omitempty"`
}

type TallyParams struct {
	DefaultTally                     v034gov.TallyParams
	CertifierUpdateSecurityVoteTally v034gov.TallyParams
	CertifierUpdateStakeVoteTally    v034gov.TallyParams
}

// GenesisState defines the governance genesis state.
type v131GovGenesisState struct {
	StartingProposalID uint64               `json:"starting_proposal_id" yaml:"starting_proposal_id"`
	Deposits           Deposits             `json:"deposits" yaml:"deposits"`
	Votes              Votes                `json:"votes" yaml:"votes"`
	Proposals          Proposals            `json:"proposals" yaml:"proposals"`
	DepositParams      DepositParams        `json:"deposit_params" yaml:"deposit_params"`
	VotingParams       v034gov.VotingParams `json:"voting_params" yaml:"voting_params"`
	TallyParams        TallyParams          `json:"tally_params" yaml:"tally_params"`
}

func migrateVoteOption(oldVoteOption v034gov.VoteOption) v040gov.VoteOption {
	switch oldVoteOption {
	case v034gov.OptionEmpty:
		return v040gov.OptionEmpty

	case v034gov.OptionYes:
		return v040gov.OptionYes

	case v034gov.OptionAbstain:
		return v040gov.OptionAbstain

	case v034gov.OptionNo:
		return v040gov.OptionNo

	case v034gov.OptionNoWithVeto:
		return v040gov.OptionNoWithVeto

	default:
		panic(fmt.Errorf("'%s' is not a valid vote option", oldVoteOption))
	}
}

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

func migrateProposalStatus(oldProposalStatus ProposalStatus) govtypes.ProposalStatus {
	switch oldProposalStatus {
	case StatusNil:
		return govtypes.StatusNil

	case StatusDepositPeriod:
		return govtypes.StatusDepositPeriod

	case StatusCertifierVotingPeriod:
		return govtypes.StatusCertifierVotingPeriod

	case StatusValidatorVotingPeriod:
		return govtypes.StatusValidatorVotingPeriod

	case StatusPassed:
		return govtypes.StatusPassed

	case StatusRejected:
		return govtypes.StatusRejected

	case StatusFailed:
		return govtypes.StatusFailed

	default:
		panic(fmt.Errorf("'%b' is not a valid proposal status", oldProposalStatus))
	}
}

// CertifierUpdateProposal adds or removes a certifier
type CertifierUpdateProposal struct {
	Title       string         `json:"title" yaml:"title"`
	Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
	Alias       string         `json:"alias" yaml:"alias"`
	Certifier   sdk.AccAddress `json:"certifier" yaml:"certifier"`
	Description string         `json:"description" yaml:"description"`
	AddOrRemove string         `json:"add_or_remove" yaml:"add_or_remove"`
}

// GetTitle returns the title of a certifier update proposal.
func (cup CertifierUpdateProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a certifier update proposal.
func (cup CertifierUpdateProposal) GetDescription() string { return cup.Description }

// GetDescription returns the routing key of a certifier update proposal.
func (cup CertifierUpdateProposal) ProposalRoute() string { return certtypes.RouterKey }

// ProposalType returns the type of a certifier update proposal.
func (cup CertifierUpdateProposal) ProposalType() string {
	return certtypes.ProposalTypeCertifierUpdate
}

// ValidateBasic runs basic stateless validity checks
func (cup CertifierUpdateProposal) ValidateBasic() error {
	return nil
}

// String implements the Stringer interface.
func (cup CertifierUpdateProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Certifier Update Proposal:
  Title:       %s
  Description: %s
  Certifier:   %s
  AddOrRemove: %s
`, cup.Title, cup.Description, cup.Certifier, cup.AddOrRemove))
	return b.String()
}

func updateOptiontoBool(opt string) bool {
	return opt == "add"
}

// ShieldClaimProposal defines the data structure of a shield claim proposal.
type ShieldClaimProposal struct {
	ProposalID  uint64         `json:"proposal_id" yaml:"proposal_id"`
	PoolID      uint64         `json:"pool_id" yaml:"pool_id"`
	PurchaseID  uint64         `json:"purchase_id" yaml:"purchase_id"`
	Loss        sdk.Coins      `json:"loss" yaml:"loss"`
	Evidence    string         `json:"evidence" yaml:"evidence"`
	Description string         `json:"description" yaml:"description"`
	Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
}

// GetTitle returns the title of a shield claim proposal.
func (scp ShieldClaimProposal) GetTitle() string {
	return fmt.Sprintf("%s:%s", strconv.FormatUint(scp.PoolID, 10), scp.Loss)
}

// GetDescription returns the description of a shield claim proposal.
func (scp ShieldClaimProposal) GetDescription() string {
	return scp.Description
}

// GetDescription returns the routing key of a shield claim proposal.
func (scp ShieldClaimProposal) ProposalRoute() string {
	return shieldtypes.RouterKey
}

// ProposalType returns the type of a shield claim proposal.
func (scp ShieldClaimProposal) ProposalType() string {
	return shieldtypes.ProposalTypeShieldClaim
}

// ValidateBasic runs basic stateless validity checks.
func (scp ShieldClaimProposal) ValidateBasic() error {
	// TODO
	return nil
}

// String implements the Stringer interface.
func (scp ShieldClaimProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Shield Claim Proposal:
  PoolID:         %d
  Loss:           %s
  Evidence:       %s
  PurchaseID:     %d
  Description:    %s
  Proposer:       %s
`, scp.PoolID, scp.Loss, scp.Evidence, scp.PurchaseID, scp.Description, scp.Proposer))
	return b.String()
}

func migrateContent(oldContent v036gov.Content) *codectypes.Any {
	var protoProposal proto.Message

	switch oldContent := oldContent.(type) {
	case v036gov.TextProposal:
		{
			protoProposal = &v040gov.TextProposal{
				Title:       oldContent.Title,
				Description: oldContent.Description,
			}
			// Convert the content into Any.
			contentAny, err := codectypes.NewAnyWithValue(protoProposal)
			if err != nil {
				panic(err)
			}

			return contentAny
		}
	case v036distr.CommunityPoolSpendProposal:
		{
			protoProposal = &v040distr.CommunityPoolSpendProposal{
				Title:       oldContent.Title,
				Description: oldContent.Description,
				Recipient:   oldContent.Recipient.String(),
				Amount:      oldContent.Amount,
			}
		}
	case v038upgrade.CancelSoftwareUpgradeProposal:
		{
			protoProposal = &v040upgrade.CancelSoftwareUpgradeProposal{
				Description: oldContent.Description,
				Title:       oldContent.Title,
			}
		}
	case v038upgrade.SoftwareUpgradeProposal:
		{
			protoProposal = &v040upgrade.SoftwareUpgradeProposal{
				Description: oldContent.Description,
				Title:       oldContent.Title,
				Plan: v040upgrade.Plan{
					Name:   oldContent.Plan.Name,
					Time:   oldContent.Plan.Time,
					Height: oldContent.Plan.Height,
					Info:   oldContent.Plan.Info,
				},
			}
		}
	case v036params.ParameterChangeProposal:
		{
			newChanges := make([]v040params.ParamChange, len(oldContent.Changes))
			for i, oldChange := range oldContent.Changes {
				newChanges[i] = v040params.ParamChange{
					Subspace: oldChange.Subspace,
					Key:      oldChange.Key,
					Value:    oldChange.Value,
				}
			}

			protoProposal = &v040params.ParameterChangeProposal{
				Description: oldContent.Description,
				Title:       oldContent.Title,
				Changes:     newChanges,
			}
		}

	// TODO: cert proposal
	case CertifierUpdateProposal:
		{
			protoProposal = &certtypes.CertifierUpdateProposal{
				Title:       oldContent.Title,
				Proposer:    oldContent.Proposer.String(),
				Alias:       oldContent.Alias,
				Certifier:   oldContent.Certifier.String(),
				Description: oldContent.Description,
				AddOrRemove: certtypes.AddOrRemove(updateOptiontoBool(oldContent.AddOrRemove)),
			}
		}
	// TODO: shield proposal
	case ShieldClaimProposal:
		{
			protoProposal = &shieldtypes.ShieldClaimProposal{
				ProposalId:  oldContent.ProposalID,
				PoolId:      oldContent.PoolID,
				PurchaseId:  oldContent.PurchaseID,
				Loss:        oldContent.Loss,
				Evidence:    oldContent.Evidence,
				Description: oldContent.Description,
				Proposer:    oldContent.Proposer.String(),
			}
		}
	default:
		panic(fmt.Errorf("%T is not a valid proposal content type", oldContent))
	}

	// Convert the content into Any.
	contentAny, err := codectypes.NewAnyWithValue(protoProposal)
	if err != nil {
		panic(err)
	}

	return contentAny
}

// Migrate accepts exported v0.36 x/gov genesis state and migrates it to
// v0.40 x/gov genesis state. The migration includes:
//
// - Convert vote option & proposal status from byte to enum.
// - Migrate proposal content to Any.
// - Convert addresses from bytes to bech32 strings.
// - Re-encode in v0.40 GenesisState.
func govMigrate(oldGovState v131GovGenesisState) *govtypes.GenesisState {
	newDeposits := make([]govtypes.Deposit, len(oldGovState.Deposits))
	for i, oldDeposit := range oldGovState.Deposits {
		newDeposits[i] = govtypes.Deposit{
			Deposit: &cosmosgov.Deposit{
				ProposalId: oldDeposit.ProposalID,
				Depositor:  oldDeposit.Depositor.String(),
				Amount:     oldDeposit.Amount,
			},
			TxHash: oldDeposit.TxHash,
		}
	}

	newVotes := make([]govtypes.Vote, len(oldGovState.Votes))
	for i, oldVote := range oldGovState.Votes {
		newVotes[i] = govtypes.Vote{
			Vote: &cosmosgov.Vote{
				ProposalId: oldVote.ProposalID,
				Voter:      oldVote.Voter.String(),
				Option:     migrateVoteOption(oldVote.Option),
			},
			TxHash: oldVote.TxHash,
		}
	}

	newProposals := make([]govtypes.Proposal, len(oldGovState.Proposals))
	for i, oldProposal := range oldGovState.Proposals {
		stat, err := ProposalStatusFromString(oldProposal.Status)
		if err != nil {
			panic(err)
		}
		newProposals[i] = govtypes.Proposal{
			Content:                 migrateContent(oldProposal.Content),
			ProposalId:              oldProposal.ProposalID,
			Status:                  migrateProposalStatus(stat),
			IsProposerCouncilMember: oldProposal.IsProposerCouncilMember,
			ProposerAddress:         oldProposal.ProposerAddress.String(),
			FinalTallyResult: v040gov.TallyResult{
				Yes:        oldProposal.FinalTallyResult.Yes,
				Abstain:    oldProposal.FinalTallyResult.Abstain,
				No:         oldProposal.FinalTallyResult.No,
				NoWithVeto: oldProposal.FinalTallyResult.NoWithVeto,
			},
			SubmitTime:      oldProposal.SubmitTime,
			DepositEndTime:  oldProposal.DepositEndTime,
			TotalDeposit:    oldProposal.TotalDeposit,
			VotingStartTime: oldProposal.VotingStartTime,
			VotingEndTime:   oldProposal.VotingEndTime,
		}
	}

	newTallyParams := govtypes.TallyParams{
		DefaultTally: &cosmosgov.TallyParams{
			Quorum:        oldGovState.TallyParams.DefaultTally.Quorum,
			Threshold:     oldGovState.TallyParams.DefaultTally.Threshold,
			VetoThreshold: oldGovState.TallyParams.DefaultTally.Veto,
		},
		CertifierUpdateSecurityVoteTally: &cosmosgov.TallyParams{
			Quorum:        oldGovState.TallyParams.CertifierUpdateSecurityVoteTally.Quorum,
			Threshold:     oldGovState.TallyParams.CertifierUpdateSecurityVoteTally.Threshold,
			VetoThreshold: oldGovState.TallyParams.CertifierUpdateSecurityVoteTally.Veto,
		},
		CertifierUpdateStakeVoteTally: &cosmosgov.TallyParams{
			Quorum:        oldGovState.TallyParams.CertifierUpdateStakeVoteTally.Quorum,
			Threshold:     oldGovState.TallyParams.CertifierUpdateStakeVoteTally.Threshold,
			VetoThreshold: oldGovState.TallyParams.CertifierUpdateStakeVoteTally.Veto,
		},
	}

	return &govtypes.GenesisState{
		StartingProposalId: oldGovState.StartingProposalID,
		Deposits:           newDeposits,
		Votes:              newVotes,
		Proposals:          newProposals,
		DepositParams: govtypes.DepositParams{
			MinDeposit:       oldGovState.DepositParams.MinDeposit,
			MaxDepositPeriod: oldGovState.DepositParams.MaxDepositPeriod,
		},
		VotingParams: v040gov.VotingParams{
			VotingPeriod: oldGovState.VotingParams.VotingPeriod,
		},
		TallyParams: newTallyParams,
	}
}
