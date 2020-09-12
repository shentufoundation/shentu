package oracle

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/hyperledger/burrow/execution/evm/abi"

	"github.com/certikfoundation/shentu/x/cvm"
)

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

// queryContract queries contract on certik-chain.
func queryContract(
	cliCtx context.CLIContext,
	queryPath, fname string,
	abiSpec, data []byte,
) (bool, string, error) {
	res, _, err := cliCtx.QueryWithData(queryPath, data)
	if err != nil {
		return false, "", fmt.Errorf("querying security primitive contract: %v", err)
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
