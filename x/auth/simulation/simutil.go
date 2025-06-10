package simulation

import (
	"fmt"
	"math/rand"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sim "github.com/cosmos/cosmos-sdk/types/simulation"
)

func RandomReasonableFees(r *rand.Rand, _ sdk.Context, spendableCoins sdk.Coins) (sdk.Coins, error) {
	if spendableCoins.Empty() {
		return nil, nil
	}

	perm := r.Perm(len(spendableCoins))
	var randCoin sdk.Coin
	for _, index := range perm {
		randCoin = spendableCoins[index]
		if !randCoin.Amount.IsZero() {
			break
		}
	}

	if randCoin.Amount.IsZero() {
		return nil, fmt.Errorf("no coins found for random fees")
	}
	// To limit the fee not exceeding ninth of remaining balances.
	// Without this limitation, it's easy to run out of accounts balance
	//  given so many simulate operations introduced by shentu modules.
	amt, err := sim.RandPositiveInt(r, randCoin.Amount.Add(math.NewInt(8)).Quo(math.NewInt(9)))
	if err != nil {
		return nil, err
	}

	// Create a random fee and verify the fees are within the account's spendable
	// balance.
	fees := sdk.NewCoins(sdk.NewCoin(randCoin.Denom, amt))

	return fees, nil
}
