package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/shentufoundation/shentu/v2/x/cert/types"
)

// RandomizedGenState generates a random GenesisState for the cert module.
func RandomizedGenState(simState *module.SimulationState) {
	var certifiers []types.Certifier

	for _, acc := range simState.Accounts {
		if simState.Rand.Intn(100) < 10 {
			certifiers = append(certifiers, types.NewCertifier(acc.Address, acc.Address, ""))
		}
	}

	// Ensure at least one certifier exists.
	if len(certifiers) == 0 {
		idx := simState.Rand.Intn(len(simState.Accounts))
		acc := simState.Accounts[idx]
		certifiers = append(certifiers, types.NewCertifier(acc.Address, acc.Address, ""))
	}

	certGenesis := types.GenesisState{
		Certifiers:        certifiers,
		NextCertificateId: 1,
	}

	bz, err := json.MarshalIndent(&certGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated cert parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&certGenesis)
}

// randCertificateType returns a random certificate type string suitable for
// AssembleContent.
func randCertificateType(r *rand.Rand) string {
	certTypes := types.IssueableCertificateTypeNames()
	return certTypes[r.Intn(len(certTypes))]
}
