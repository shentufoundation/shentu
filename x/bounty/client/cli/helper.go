package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// GetEncryptionKey get the key in the program for information encryption
func GetEncryptionKey(cmd *cobra.Command, programId uint64) (*codectypes.Any, error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return nil, err
	}
	queryClient := types.NewQueryClient(clientCtx)

	// Query the program
	res, err := queryClient.Program(
		cmd.Context(),
		&types.QueryProgramRequest{
			ProgramId: programId,
		})

	if err != nil {
		return nil, err
	}
	return res.GetProgram().EncryptionKey, nil
}

// GetFinding get finding details
func GetFinding(cmd *cobra.Command, findingId uint64) (*types.Finding, error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return nil, err
	}
	queryClient := types.NewQueryClient(clientCtx)

	// Query the program
	res, err := queryClient.Finding(
		cmd.Context(),
		&types.QueryFindingRequest{
			FindingId: findingId,
		})

	if err != nil {
		return nil, err
	}
	return &res.Finding, nil
}
