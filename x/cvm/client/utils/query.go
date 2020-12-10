package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/certikfoundation/shentu/x/cvm/internal/types"
)

// QueryCVMAccount is to query the cvm contract related info by addresss
func QueryCVMAccount(cliCtx client.Context, address string, account *authtypes.BaseAccount) (*types.CVMAccount, error) {
	cvmCodeRes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cvm/code/%s", address), nil)
	if err != nil {
		return nil, err
	}
	var cvmCodeOut types.QueryResCode
	cliCtx.LegacyAmino.MustUnmarshalJSON(cvmCodeRes, &cvmCodeOut)
	cvmCode := cvmCodeOut.Code
	if cvmCode.String() == "" {
		return nil, ErrEmptyCVMCode
	}

	cvmAbiRes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cvm/abi/%s", address), nil)
	if err != nil {
		return nil, err
	}
	var cvmAbiOut types.QueryResAbi
	cliCtx.LegacyAmino.MustUnmarshalJSON(cvmAbiRes, &cvmAbiOut)
	cvmAbi := string(cvmAbiOut.Abi)
	if cvmAbi == "" {
		return nil, ErrEmptyCVMAbi
	}

	cvmAcc := types.CVMAccount{
		BaseAccount: account,
		Code:        string(cvmCodeOut.Code),
		Abi:         cvmAbi,
	}

	return &cvmAcc, nil
}
