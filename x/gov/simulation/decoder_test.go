package simulation

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/certikfoundation/shentu/x/gov/internal/types"
)

func makeTestCodec() *codec.Codec {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	govTypes.RegisterCodec(cdc)
	return cdc
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()

	rand.Seed(time.Now().UnixNano())

	endTime := time.Now().UTC()

	content := govTypes.ContentFromProposalType("test", "test", types.ProposalTypeText)
	proposalID := rand.Uint64()
	proposer := RandomAccount()
	isMember := 1 == rand.Intn(2)
	proposal := types.NewProposal(content, proposalID, proposer.Address, isMember, endTime, endTime.Add(24*time.Hour))

	proposalIDBz := make([]byte, 8)
	binary.LittleEndian.PutUint64(proposalIDBz, proposalID)

	depositor := RandomAccount()
	txhash := "2300092389009f098099"
	deposit := types.NewDeposit(proposalID, depositor.Address, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())), txhash)
	voter := RandomAccount()
	vote := types.NewVote(proposalID, voter.Address, govTypes.OptionYes, txhash)

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: govTypes.ProposalKey(proposalID), Value: cdc.MustMarshalBinaryLengthPrefixed(&proposal)},
		tmkv.Pair{Key: govTypes.InactiveProposalQueueKey(proposalID, endTime), Value: proposalIDBz},
		tmkv.Pair{Key: govTypes.DepositKey(proposalID, depositor.Address), Value: cdc.MustMarshalBinaryLengthPrefixed(&deposit)},
		tmkv.Pair{Key: govTypes.VoteKey(proposalID, voter.Address), Value: cdc.MustMarshalBinaryLengthPrefixed(&vote)},
		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
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
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name) // nolint
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name) // nolint
			}
		})
	}
}

func RandomAccount() simulation.Account {
	privkeySeed := make([]byte, 15)
	rand.Read(privkeySeed)

	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	pubKey := privKey.PubKey()
	address := sdk.AccAddress(pubKey.Address())

	return simulation.Account{
		PrivKey: privKey,
		PubKey:  pubKey,
		Address: address,
	}
}
