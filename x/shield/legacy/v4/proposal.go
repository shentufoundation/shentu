package v4

import (
	"fmt"
	"strings"
)

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