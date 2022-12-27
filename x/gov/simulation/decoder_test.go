package simulation_test

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	sim "github.com/cosmos/cosmos-sdk/types/simulation"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	shentuapp "github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/x/gov/keeper"
	. "github.com/shentufoundation/shentu/v2/x/gov/simulation"
)

func TestDecodeStore(t *testing.T) {
	cdc := shentuapp.MakeEncodingConfig()
	dec := NewDecodeStore(cdc.Marshaler)

	rand.Seed(time.Now().UnixNano())

	endTime := time.Now().UTC()

	content := govtypes.ContentFromProposalType("test", "test", govtypes.ProposalTypeText)
	proposalID := rand.Uint64()
	proposal, _ := govtypes.NewProposal(content, proposalID, endTime, endTime.Add(24*time.Hour))

	proposalIDBz := make([]byte, 8)
	binary.LittleEndian.PutUint64(proposalIDBz, proposalID)

	depositor := RandomAccount()
	deposit := govtypes.NewDeposit(proposalID, depositor.Address, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())))
	voter := RandomAccount()
	options := govtypes.NewNonSplitVoteOption(govtypes.OptionYes)
	vote := govtypes.NewVote(proposalID, voter.Address, options)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: keeper.ProposalKey(proposalID), Value: cdc.Marshaler.MustMarshal(&proposal)},
			{Key: govtypes.InactiveProposalQueueKey(proposalID, endTime), Value: proposalIDBz},
			{Key: govtypes.DepositKey(proposalID, depositor.Address), Value: cdc.Marshaler.MustMarshal(&deposit)},
			{Key: govtypes.VoteKey(proposalID, voter.Address), Value: cdc.Marshaler.MustMarshal(&vote)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"proposals", fmt.Sprintf("%v\n%v", proposal, proposal)},
		{"proposal IDs", fmt.Sprintf("%d\n%d", proposalID, proposalID)},
		{"deposits", fmt.Sprintf("%v\n%v", deposit, deposit)},
		{"votes", fmt.Sprintf("%v\n%v", vote, vote)},
		{"other", ""},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch i { // nolint
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name) // nolint
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name) // nolint
			}
		})
	}
}

func RandomAccount() sim.Account {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	address := sdk.AccAddress(pubKey.Address())

	return sim.Account{
		PrivKey: privKey,
		PubKey:  pubKey,
		Address: address,
	}
}
