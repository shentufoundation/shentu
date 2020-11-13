package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	// ProposalTypeShieldClaim defines the type for a ShieldClaimProposal.
	ProposalTypeShieldClaim = "ShieldClaim"
)

// Assert ShieldClaimProposal implements govTypes.Content at compile-time.
var _ govTypes.Content = ShieldClaimProposal{}

func init() {
	govTypes.RegisterProposalType(ProposalTypeShieldClaim)
	govTypes.RegisterProposalTypeCodec(ShieldClaimProposal{}, "shield/ShieldClaimProposal")
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

// NewShieldClaimProposal creates a new shield claim proposal.
func NewShieldClaimProposal(poolID uint64, loss sdk.Coins, purchaseID uint64, evidence, description string, proposer sdk.AccAddress) ShieldClaimProposal {
	return ShieldClaimProposal{
		PoolID:      poolID,
		Loss:        loss,
		Evidence:    evidence,
		PurchaseID:  purchaseID,
		Description: description,
		Proposer:    proposer,
	}
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
	return RouterKey
}

// ProposalType returns the type of a shield claim proposal.
func (scp ShieldClaimProposal) ProposalType() string {
	return ProposalTypeShieldClaim
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

// LockedCollateral defines the data type of locked collateral for a claim proposal.
type LockedCollateral struct {
	ProposalID uint64  `json:"proposal_id" yaml:"proposal_id"`
	Amount     sdk.Int `json:"locked_coins" yaml:"locked_coins"`
}

// NewLockedCollateral returns a new LockedCollateral instance.
func NewLockedCollateral(proposalID uint64, lockedAmt sdk.Int) LockedCollateral {
	return LockedCollateral{
		ProposalID: proposalID,
		Amount:     lockedAmt,
	}
}

// NewUnbondingDelegation returns a new UnbondingDelegation instance.
func NewUnbondingDelegation(delAddr sdk.AccAddress, valAddr sdk.ValAddress, entry stakingTypes.UnbondingDelegationEntry) staking.UnbondingDelegation {
	return staking.UnbondingDelegation{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Entries:          []stakingTypes.UnbondingDelegationEntry{entry},
	}
}

// Reimbursement stores information of a reimbursement.
type Reimbursement struct {
	Amount      sdk.Coins      `json:"amount" yaml:"amount"`
	Beneficiary sdk.AccAddress `json:"beneficiary" yaml:"beneficiary"`
	PayoutTime  time.Time      `json:"payout_time" yaml:"payout_time"`
}

// NewReimbursement returns a new Reimbursement instance.
func NewReimbursement(amount sdk.Coins, beneficiary sdk.AccAddress, payoutTime time.Time) Reimbursement {
	return Reimbursement{
		Amount:      amount,
		Beneficiary: beneficiary,
		PayoutTime:  payoutTime,
	}
}

// ProposalIDReimbursementPair stores information of a reimbursement and corresponding proposal ID.
type ProposalIDReimbursementPair struct {
	ProposalID    uint64
	Reimbursement Reimbursement
}

// NewProposalIDReimbursementPair returns a new ProposalIDReimbursementPair instance.
func NewProposalIDReimbursementPair(proposalID uint64, reimbursement Reimbursement) ProposalIDReimbursementPair {
	return ProposalIDReimbursementPair{
		ProposalID:    proposalID,
		Reimbursement: reimbursement,
	}
}
