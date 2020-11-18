package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Get validator count
	r.HandleFunc("/staking/all_validators", allValidatorsHandlerFn(cliCtx)).Methods("GET")
}

type AllValidatorsResult struct {
	Count int
	types.Validators
}

// HTTP request handler to query complete list of validators
func allValidatorsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, _, err := rest.ParseHTTPArgsWithLimit(r, 100)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		resKVs, height, err := cliCtx.QuerySubspace(types.ValidatorsKey, types.StoreKey)
		if err != nil {
			return
		}

		var validators types.Validators
		for _, kv := range resKVs {
			validators = append(validators, types.MustUnmarshalValidator(cliCtx.Codec, kv.Value))
		}

		res := AllValidatorsResult{len(resKVs), validators}
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
