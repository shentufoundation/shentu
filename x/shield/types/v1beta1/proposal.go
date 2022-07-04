package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

// NewShieldClaimProposal creates a new shield claim proposal.
func NewShieldClaimProposal(poolID uint64, loss sdk.Coins, purchaseID uint64, evidence, description string, proposer sdk.AccAddress) *ShieldClaimProposal {
	return &ShieldClaimProposal{
		PoolId:      poolID,
		Loss:        loss,
		Evidence:    evidence,
		PurchaseId:  purchaseID,
		Description: description,
		Proposer:    proposer.String(),
	}
}

// GetTitle returns the title of a shield claim proposal.
func (scp ShieldClaimProposal) GetTitle() string {
	return fmt.Sprintf("%s:%s", strconv.FormatUint(scp.PoolId, 10), scp.Loss)
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
`, scp.PoolId, scp.Loss, scp.Evidence, scp.PurchaseId, scp.Description, scp.Proposer))
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
func NewUnbondingDelegation(delAddr, valAddr string, entry stakingTypes.UnbondingDelegationEntry) stakingTypes.UnbondingDelegation {
	return stakingTypes.UnbondingDelegation{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Entries:          []stakingTypes.UnbondingDelegationEntry{entry},
	}
}

// NewReimbursement returns a new Reimbursement instance.
func NewReimbursement(amount sdk.Coins, beneficiary sdk.AccAddress, payoutTime time.Time) Reimbursement {
	return Reimbursement{
		Amount:      amount,
		Beneficiary: beneficiary.String(),
		PayoutTime:  payoutTime,
	}
}

// NewProposalIDReimbursementPair returns a new ProposalIDReimbursementPair instance.
func NewProposalIDReimbursementPair(proposalID uint64, reimbursement Reimbursement) ProposalIDReimbursementPair {
	return ProposalIDReimbursementPair{
		ProposalId:    proposalID,
		Reimbursement: reimbursement,
	}
}
