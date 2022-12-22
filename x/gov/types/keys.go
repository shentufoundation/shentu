package types

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// CertVotesKeyPrefix "cert"
var CertVotesKeyPrefix = []byte{0x63, 0x65, 0x72, 0x74}

// CertVotesKey gets the first part of the cert votes key based on the proposalID
func CertVotesKey(proposalID uint64) []byte {
	return append(CertVotesKeyPrefix, govtypes.GetProposalIDBytes(proposalID)...)
}
