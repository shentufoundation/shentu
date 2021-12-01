package simapp

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	log "github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// SimApp testing.
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func Setup(isCheckTx bool) *SimApp {
	db := dbm.NewMemDB()
	app := NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, MakeTestEncodingConfig(), simapp.EmptyAppOptions{})
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt
func AddTestAddrs(app *SimApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		testAddrs[i] = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	}

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		if err := simapp.FundAccount(app.BankKeeper, ctx, addr, initCoins); err != nil {
			panic(err)
		}
	}
	return testAddrs
}

// AddTestAddrsFromPubKeys adds the addresses into the SimApp providing only the public keys.
func AddTestAddrsFromPubKeys(app *SimApp, ctx sdk.Context, pubKeys []cryptotypes.PubKey, accAmt sdk.Int) {
	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, pubKey := range pubKeys {
		if err := simapp.FundAccount(app.BankKeeper, ctx, sdk.AccAddress(pubKey.Address()), initCoins); err != nil {
			panic(err)
		}
	}
}

// CreateTestPubKeys returns a total of numPubKeys public keys in ascending order.
func CreateTestPubKeys(numPubKeys int) []cryptotypes.PubKey {
	var publicKeys []cryptotypes.PubKey
	var buffer bytes.Buffer

	// start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") // base pubkey string
		buffer.WriteString(numString)                                                       // adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKeyFromHex(buffer.String()))
		buffer.Reset()
	}

	return publicKeys
}

// NewPubKeyFromHex returns a PubKey from a hex string.
func NewPubKeyFromHex(pk string) (res cryptotypes.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	if len(pkBytes) != ed25519.PubKeySize {
		panic(errors.Wrap(errors.ErrInvalidPubKey, "invalid pubkey size"))
	}
	return &ed25519.PubKey{Key: pkBytes}
}
