package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	runnerTypes "github.com/certikfoundation/shentu/toolsets/oracle-operator/runner/types"
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

		var rawJson []byte
		if viper.GetBool(flags.FlagIndentResponse) {
			rawJson, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				return sdk.TxResponse{}, err
			}
		} else {
			rawJson = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", rawJson)

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

// BuildQueryInsight queries insight from Security Primitive Interface Contract from certik-chain.
func BuildQueryInsight(cliCtx context.CLIContext, calleeString string, function string, args []string) (bool, string, error) {
	calleeAddr, err := sdk.AccAddressFromBech32(calleeString)
	if err != nil {
		return false, "", err
	}
	accGetter := authtypes.NewAccountRetriever(cliCtx)
	if err := accGetter.EnsureExists(calleeAddr); err != nil {
		return false, "", err
	}

	abiSpec, err := queryAbi(cliCtx, "cvm", calleeString)
	if err != nil {
		return false, "", err
	}
	logger := logging.NewNoopLogger()

	data, err := parseData(function, abiSpec, args, logger)
	if err != nil {
		return false, "", err
	}

	value := viper.GetUint64(runnerTypes.FlagValue)
	msg := cvm.NewMsgCall(cliCtx.GetFromAddress(), calleeAddr, value, data)
	if err := msg.ValidateBasic(); err != nil {
		return false, "", err
	}
	// Decode abiSpec to check if the called function's type is view or pure.
	// If it is, reroute to query.
	var abiEntries []runnerTypes.AbiEntry
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
		return queryInsightAndReturn(cliCtx, queryPath, function, abiSpec, data)
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

// queryAbi queries ABI from certik chain
func queryAbi(cliCtx context.CLIContext, queryRoute string, addr string) ([]byte, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/abi/%s", queryRoute, addr), nil)
	if err != nil {
		return nil, err
	}

	var out cvm.QueryResAbi
	cliCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Abi, nil
}

// queryInsightAndReturn queries Security Primitive Interface and get insight for primitive
func queryInsightAndReturn(cliCtx context.CLIContext, queryPath, fname string, abiSpec, data []byte,
) (bool, string, error) {
	res, _, err := cliCtx.QueryWithData(queryPath, data)
	if err != nil {
		return false, "", fmt.Errorf("querying CVM contract code: %v", err)
	}
	var out cvm.QueryResView
	err = json.Unmarshal(res, &out)
	if err != nil {
		return false, "", err
	}
	ret, err := abi.DecodeFunctionReturn(string(abiSpec), fname, out.Ret)
	if err != nil {
		return false, "", fmt.Errorf("decoding function return: %v", err)
	}
	if len(ret) != 2 {
		return false, "", fmt.Errorf("mismatch return length: %v", ret)
	}

	retBool, err := strconv.ParseBool(ret[0].Value)
	if err != nil {
		return false, "", fmt.Errorf("decoding function return: %v", err)
	}

	return retBool, ret[1].Value, nil
}
