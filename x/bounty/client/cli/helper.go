package cli

import (
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

// GetEncryptionKey get the key in the program for information encryption
func GetEncryptionKey(cmd *cobra.Command, programID uint64) (*ecies.PublicKey, error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return nil, err
	}
	queryClient := types.NewQueryClient(clientCtx)

	// Query the program
	res, err := queryClient.Program(
		cmd.Context(),
		&types.QueryProgramRequest{
			ProgramId: programID,
		})

	if err != nil {
		return nil, err
	}

	var encryptionKey types.EciesPubKey
	if err = proto.Unmarshal(res.Program.EncryptionKey.GetValue(), &encryptionKey); err != nil {
		return nil, err
	}

	pubEcdsa, err := crypto.UnmarshalPubkey(encryptionKey.EncryptionKey)
	if err != nil {
		return nil, err
	}
	eciesEncKey := ecies.ImportECDSAPublic(pubEcdsa)

	return eciesEncKey, nil
}

// GetFinding get finding details
func GetFinding(cmd *cobra.Command, findingID uint64) (*types.Finding, error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return nil, err
	}
	queryClient := types.NewQueryClient(clientCtx)

	// Query the program
	res, err := queryClient.Finding(
		cmd.Context(),
		&types.QueryFindingRequest{
			FindingId: findingID,
		})

	if err != nil {
		return nil, err
	}
	return &res.Finding, nil
}

func GetNextProgramID(cmd *cobra.Command) (uint64, error) {
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		return 0, err
	}
	queryClient := types.NewQueryClient(clientCtx)
	res, err := queryClient.NextProgramID(cmd.Context(), &types.QueryNextProgramIDRequest{})
	if err != nil {
		return 0, err
	}
	return res.NextProgramId, nil
}