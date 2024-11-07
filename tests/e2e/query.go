package e2e

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	sdk "github.com/cosmos/cosmos-sdk/types"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	sdkgovtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	bountytypes "github.com/shentufoundation/shentu/v2/x/bounty/types"
	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	govtypes "github.com/shentufoundation/shentu/v2/x/gov/types/v1"
)

func connectGrpc(grpcEndpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(grpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect %s: %v", grpcEndpoint, err)
	}
	return conn, nil
}

func queryShentuAllBalances(endpoint, address string) (sdk.Coins, error) {
	body, err := httpGet(fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", endpoint, address))
	if err != nil {
		return nil, err
	}
	var balancesResp banktypes.QueryAllBalancesResponse
	if err := cdc.UnmarshalJSON(body, &balancesResp); err != nil {
		return nil, err
	}
	return balancesResp.Balances, nil
}

func queryShentuBalance(endpoint, address, denom string) (sdk.Coin, error) {
	var balance sdk.Coin
	body, err := httpGet(fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s/by_denom?denom=%s", endpoint, address, denom))
	if err != nil {
		return balance, err
	}
	var balanceResp banktypes.QuerySpendableBalanceByDenomResponse
	if err := cdc.UnmarshalJSON(body, &balanceResp); err != nil {
		return balance, err
	}
	return *balanceResp.Balance, nil
}

func queryShentuTx(endpoint, txHash string) error {
	body, err := httpGet(fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", endpoint, txHash))
	if err != nil {
		return err
	}
	var resp txtypes.GetTxResponse
	if err := cdc.UnmarshalJSON(body, &resp); err != nil {
		return err
	}
	if resp.TxResponse.Code != 0 {
		return fmt.Errorf("tx %s failed with status code %v", txHash, resp.TxResponse.Code)
	}
	return nil
}

func queryDelegatorWithdrawalAddress(endpoint, delegatorAddr string) (distribtypes.QueryDelegatorWithdrawAddressResponse, error) {
	var res distribtypes.QueryDelegatorWithdrawAddressResponse

	body, err := httpGet(fmt.Sprintf("%s/cosmos/distribution/v1beta1/delegators/%s/withdraw_address", endpoint, delegatorAddr))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryDelegation(endpoint, validator, delegator string) (stakingtypes.QueryDelegationResponse, error) {
	var res stakingtypes.QueryDelegationResponse

	body, err := httpGet(fmt.Sprintf("%s/cosmos/staking/v1beta1/validators/%s/delegations/%s", endpoint, validator, delegator))
	if err != nil {
		return res, err
	}

	if err = cdc.UnmarshalJSON(body, &res); err != nil {
		return res, err
	}
	return res, nil
}

func queryProgram(grpcEndpoint, programID string) (bountytypes.Program, error) {
	conn, err := connectGrpc(grpcEndpoint)
	defer conn.Close()

	client := bountytypes.NewQueryClient(conn)

	res, err := client.Program(context.Background(), &bountytypes.QueryProgramRequest{
		ProgramId: programID,
	})
	if err != nil {
		return bountytypes.Program{}, err
	}
	return res.Program, nil
}

func queryFinding(grpcEndpoint, findingID string) (bountytypes.Finding, error) {
	conn, err := connectGrpc(grpcEndpoint)
	defer conn.Close()

	client := bountytypes.NewQueryClient(conn)

	res, err := client.Finding(context.Background(), &bountytypes.QueryFindingRequest{
		FindingId: findingID,
	})
	if err != nil {
		return bountytypes.Finding{}, err
	}
	return res.Finding, nil
}

func queryFindingFingerprint(grpcEndpoint, findingID string) (string, error) {
	conn, err := connectGrpc(grpcEndpoint)
	defer conn.Close()

	client := bountytypes.NewQueryClient(conn)

	res, err := client.FindingFingerprint(context.Background(), &bountytypes.QueryFindingFingerprintRequest{
		FindingId: findingID,
	})
	if err != nil {
		return "", err
	}
	return res.Fingerprint, nil
}

func queryCertificate(grpcEndpoint, content, certificate string) (bool, error) {
	conn, _ := connectGrpc(grpcEndpoint)
	defer conn.Close()

	client := certtypes.NewQueryClient(conn)
	res, err := client.Certificates(context.Background(), &certtypes.QueryCertificatesRequest{
		CertificateType: certificate,
		Pagination: &querytypes.PageRequest{
			Limit:  5000,
			Offset: 0,
		},
	})
	if err != nil {
		return false, err
	}
	for _, item := range res.Certificates {
		tmp := certtypes.AssembleContent(certificate, "")
		err := cdc.UnpackAny(item.Certificate.Content, &tmp)
		if err != nil {
			return false, err
		}
		if tmp.GetContent() == content {
			return true, nil
		}
	}
	return false, nil
}

func queryCertVoted(grpcEndpoint string, proposalID uint64) (bool, error) {
	conn, err := connectGrpc(grpcEndpoint)
	defer conn.Close()

	client := govtypes.NewCustomQueryClient(conn)

	res, err := client.CertVoted(context.Background(), &govtypes.QueryCertVotedRequest{
		ProposalId: proposalID,
	})
	if err != nil {
		return false, err
	}
	return res.CertVoted, nil
}

func queryProposal(grpcEndpoint string, proposalID uint64) (*sdkgovtypes.Proposal, error) {
	conn, err := connectGrpc(grpcEndpoint)
	defer conn.Close()

	client := govtypesv1.NewQueryClient(conn)

	res, err := client.Proposal(context.Background(), &sdkgovtypes.QueryProposalRequest{
		ProposalId: proposalID,
	})
	if err != nil {
		return nil, err
	}
	return res.Proposal, nil
}
