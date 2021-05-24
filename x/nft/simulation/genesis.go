// copied from irisnet/irismod@875a1d1

package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"

	certtypes "github.com/certikfoundation/shentu/x/cert/types"
	"github.com/certikfoundation/shentu/x/nft/types"
)

const (
	kitties = "kitties"
	doggos  = "doggos"
)

// RandomizedGenState generates a random GenesisState for nft
func RandomizedGenState(simState *module.SimulationState) {
	collections := nfttypes.NewCollections(
		nfttypes.NewCollection(
			nfttypes.Denom{
				Id:      doggos,
				Name:    doggos,
				Schema:  "",
				Creator: "",
			},
			nfttypes.NFTs{},
		),
		nfttypes.NewCollection(
			nfttypes.Denom{
				Id:      kitties,
				Name:    kitties,
				Schema:  "",
				Creator: "",
			},
			nfttypes.NFTs{}),
	)
	for _, acc := range simState.Accounts {
		// 10% of accounts own an NFT
		if simState.Rand.Intn(100) < 10 {
			baseNFT := nfttypes.NewBaseNFT(
				RandnNFTID(simState.Rand, nfttypes.MinDenomLen, nfttypes.MaxDenomLen), // id
				simtypes.RandStringOfLength(simState.Rand, 10),
				acc.Address,
				simtypes.RandStringOfLength(simState.Rand, 45), // tokenURI
				simtypes.RandStringOfLength(simState.Rand, 10),
			)

			// 50% doggos and 50% kitties
			if simState.Rand.Intn(100) < 50 {
				collections[0].Denom.Creator = baseNFT.Owner
				collections[0] = collections[0].AddNFT(baseNFT)
			} else {
				collections[1].Denom.Creator = baseNFT.Owner
				collections[1] = collections[1].AddNFT(baseNFT)
			}
		}
	}

	certbz := simState.GenState[certtypes.ModuleName]
	var certGenState certtypes.GenesisState
	simState.Cdc.MustUnmarshalJSON(certbz, &certGenState)
	var admins []types.Admin
	for _, c := range certGenState.Certifiers {
		admins = append(admins, types.Admin{Address: c.Address})
	}

	nftGenesis := types.NewGenesisState(collections, admins)

	bz, err := json.MarshalIndent(nftGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", nfttypes.ModuleName, bz)

	simState.GenState[nfttypes.ModuleName] = simState.Cdc.MustMarshalJSON(nftGenesis)
}

func RandnNFTID(r *rand.Rand, min, max int) string {
	n := simtypes.RandIntBetween(r, min, max)
	id := simtypes.RandStringOfLength(r, n)
	return strings.ToLower(id)
}
