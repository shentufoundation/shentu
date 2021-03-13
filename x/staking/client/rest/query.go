package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// Get validator count
	r.HandleFunc("/staking/all_validators", allValidatorsHandlerFn(cliCtx)).Methods("GET")
}

type allValidatorsResult struct {
	Count int
	types.Validators
}

// HTTP request handler to query complete list of validators
func allValidatorsHandlerFn(cliCtx client.Context) http.HandlerFunc {
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

		queryClient := types.NewQueryClient(cliCtx)

		result, err := queryClient.Validators(context.Background(), &types.QueryValidatorsRequest{
			// Leaving status and pageReq empty on purpose to query all validators.
		})
		if err != nil {
			panic(err)
		}

		res := allValidatorsResult{len(result.Validators), result.Validators}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
