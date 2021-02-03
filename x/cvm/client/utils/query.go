package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	auth_types "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/certikfoundation/shentu/x/cvm/types"
)

// QueryCVMAccount is to query the cvm contract related info by addresss
func QueryCVMAccount(cliCtx context.CLIContext, address string, account exported.Account) (*types.CVMAccount, error) {
	cvmCodeRes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cvm/code/%s", address), nil)
	if err != nil {
		return nil, err
	}
	var cvmCodeOut types.QueryResCode
	cliCtx.Codec.MustUnmarshalJSON(cvmCodeRes, &cvmCodeOut)
	cvmCode := cvmCodeOut.Code
	if cvmCode.String() == "" {
		return nil, ErrEmptyCVMCode
	}

	cvmAbiRes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/cvm/abi/%s", address), nil)
	if err != nil {
		return nil, err
	}
	var cvmAbiOut types.QueryResAbi
	cliCtx.Codec.MustUnmarshalJSON(cvmAbiRes, &cvmAbiOut)
	cvmAbi := string(cvmAbiOut.Abi)
	if cvmAbi == "" {
		return nil, ErrEmptyCVMAbi
	}

	baseAcc, ok := account.(*auth_types.BaseAccount)
	if !ok {
		return nil, ErrBaseAccount
	}

	cvmAcc := types.NewCVMAccount(
		baseAcc,
		cvmCodeOut.Code,
		cvmAbi,
	)

	return cvmAcc, nil
}
