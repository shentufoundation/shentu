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

	"github.com/certikfoundation/shentu/simapp"
	. "github.com/certikfoundation/shentu/x/gov/simulation"
	"github.com/certikfoundation/shentu/x/gov/types"
)

func TestDecodeStore(t *testing.T) {
	cdc := simapp.MakeTestEncodingConfig()
	dec := NewDecodeStore(cdc.Marshaler)

	rand.Seed(time.Now().UnixNano())

	endTime := time.Now().UTC()

	content := govtypes.ContentFromProposalType("test", "test", govtypes.ProposalTypeText)
	proposalID := rand.Uint64()
	proposer := RandomAccount()
	isMember := 1 == rand.Intn(2)
	proposal, _ := types.NewProposal(content, proposalID, proposer.Address, isMember, endTime, endTime.Add(24*time.Hour))

	proposalIDBz := make([]byte, 8)
	binary.LittleEndian.PutUint64(proposalIDBz, proposalID)

	depositor := RandomAccount()
	txhash := "2300092389009f098099"
	deposit := types.NewDeposit(proposalID, depositor.Address, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())), txhash)
	voter := RandomAccount()
	vote := types.NewVote(proposalID, voter.Address, govtypes.OptionYes, txhash)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: govtypes.ProposalKey(proposalID), Value: cdc.Marshaler.MustMarshalBinaryBare(&proposal)},
			{Key: govtypes.InactiveProposalQueueKey(proposalID, endTime), Value: proposalIDBz},
			{Key: govtypes.DepositKey(proposalID, depositor.Address), Value: cdc.Marshaler.MustMarshalBinaryBare(&deposit)},
			{Key: govtypes.VoteKey(proposalID, voter.Address), Value: cdc.Marshaler.MustMarshalBinaryBare(&vote)},
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
