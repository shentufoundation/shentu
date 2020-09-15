package oracle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/logging"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
	"github.com/certikfoundation/shentu/x/cvm"
	"github.com/certikfoundation/shentu/x/cvm/compile"
)

// CompleteAndBroadcastTx is adopted from auth.CompleteAndBroadcastTxCLI. The original function prints out response.
func CompleteAndBroadcastTx(cliCtx context.CLIContext, txBldr authtypes.TxBuilder, msgs []sdk.Msg) (sdk.TxResponse, error) {
	txBldr, err := utils.PrepareTxBuilder(txBldr, cliCtx)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	fromName := cliCtx.GetFromName()

	if txBldr.SimulateAndExecute() || cliCtx.Simulate {
		txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		if err != nil {
			return sdk.TxResponse{}, err
		}

		gasEst := utils.GasEstimateResponse{GasEstimate: txBldr.Gas()}
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", gasEst.String())
	}

	if cliCtx.Simulate {
		return sdk.TxResponse{}, nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return sdk.TxResponse{}, err
		}

		var rawJSON []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			rawJSON, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				return sdk.TxResponse{}, err
			}
		} else {
			rawJSON = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", rawJSON)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "canceled transaction")
			return sdk.TxResponse{}, err
		}
	}

	// build and sign the transaction
	txBytes, err := txBldr.BuildAndSign(fromName, keys.DefaultKeyPass, msgs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return res, nil
}

// callContract calls contract on certik-chain.
func callContract(ctx types.Context, calleeString string, function string, args []string) (bool, string, error) {
	cliCtx := ctx.ClientContext()

	calleeAddr, err := sdk.AccAddressFromBech32(calleeString)
	if err != nil {
		return false, "", err
	}
	accGetter := authtypes.NewAccountRetriever(cliCtx)
	if err := accGetter.EnsureExists(calleeAddr); err != nil {
		return false, "", err
	}

	abiSpec, err := queryAbi(cliCtx, cvm.QuerierRoute, calleeString)
	if err != nil {
		return false, "", err
	}

	data, err := parseData(function, abiSpec, args, logging.NewNoopLogger())
	if err != nil {
		return false, "", err
	}

	// Decode abiSpec to check if the called function's type is view or pure.
	// If it is, reroute to query.
	var abiEntries []types.ABIEntry
	err = json.Unmarshal(abiSpec, &abiEntries)
	if err != nil {
		return false, "", err
	}
	for _, entry := range abiEntries {
		if entry.Name != function {
			continue
		}
		if entry.Type != "view" && entry.Type != "pure" {
			return false, "", fmt.Errorf("getInsight function should be view or pure function")
		}
		queryPath := fmt.Sprintf("custom/%s/view/%s/%s", cvm.QuerierRoute, cliCtx.GetFromAddress(), calleeAddr)
		return queryContract(cliCtx, queryPath, function, abiSpec, data)
	}
	return false, "", fmt.Errorf("function %s was not found in abi", function)
}

// parseData parses Data for contract on certik chain
func parseData(function string, abiSpec []byte, args []string, logger *logging.Logger) ([]byte, error) {
	params := make([]interface{}, 0)

	if string(abiSpec) == compile.NoABI {
		panic("No ABI registered for this contract. Use --raw flag to submit raw bytecode.")
	}

	for _, arg := range args {
		var argi interface{}
		argi = arg
		for _, prefix := range []string{common.Bech32MainPrefix, common.Bech32PrefixConsAddr, common.Bech32PrefixAccAddr} {
			if strings.HasPrefix(arg, prefix) && ((len(arg) - len(prefix)) == 39) {
				data, _ := sdk.GetFromBech32(arg, prefix)
				var err error
				argi, err = crypto.AddressFromBytes(data)
				if err != nil {
					return nil, err
				}
				break
			}
		}
		params = append(params, argi)
	}

	data, _, err := abi.EncodeFunctionCall(string(abiSpec), function, logger, params...)
	return data, err
}
