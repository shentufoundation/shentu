package cmd

import (
	"fmt"
	"sort"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	captypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	// crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidexported "github.com/cosmos/cosmos-sdk/x/evidence/exported"
	evidtypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	certtypes "github.com/shentufoundation/shentu/v2/x/cert/types"
	cvmtypes "github.com/shentufoundation/shentu/v2/x/cvm/types"
	stgovtypes "github.com/shentufoundation/shentu/v2/x/gov/types"
	oracletypes "github.com/shentufoundation/shentu/v2/x/oracle/types"
	shieldtypes "github.com/shentufoundation/shentu/v2/x/shield/types"

	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	icahosttypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	ibcexported "github.com/cosmos/ibc-go/v4/modules/core/exported"
	ibcconntypes "github.com/cosmos/ibc-go/v4/modules/core/03-connection/types"
	ibcchantypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
)

type OneKey struct {
	prefix []byte
	ptr interface{}
	marshalWay int //1: Marshal; 2: MarshalLengthPrefixed; 3: MarshalInterface
}

var (
	ai authtypes.AccountI
	evidi evidexported.Evidence
	pki cryptotypes.PubKey
	ti  oracletypes.TaskI
	ibcCS ibcexported.ClientState
	ibcConsS    ibcexported.ConsensusState
)
type OKs []OneKey
type KeyTypes map[string]OKs

func (kt KeyTypes) SortedModuleNames() []string {
	allModules := make([]string, 0, len(kt))
	for k := range kt {
		allModules = append(allModules, k)
	}
	sort.Strings(allModules)
	return allModules
}

func (ks OKs) MatchKey(key []byte) (OneKey, error) {
	for _, i := range ks {
		//Todo: cache the compiled exp
		m := regexp.MustCompile(string(i.prefix))
		if m.FindIndex(key) != nil {
			return i, nil
		}
	}
	return OneKey{}, fmt.Errorf("key not matched")
}

var allKeys = KeyTypes {
	certtypes.StoreKey: {
		{certtypes.CertifiersStoreKey(),        &certtypes.Certifier{}, 2},
		{certtypes.CertifierAliasesStoreKey(),  &certtypes.Certifier{}, 2},
		{certtypes.PlatformsStoreKey(),         &certtypes.Platform{}, 1},
		{certtypes.CertificatesStoreKey(),      &certtypes.Certificate{}, 1},
		{certtypes.LibrariesStoreKey(),         &certtypes.Library{}, 2},
		{certtypes.NextCertificateIDStoreKey(), nil, 4}, //binary.LittleEndian.Uint64
		{certtypes.ValidatorsStoreKey(),        nil, 4}, //this key is not used
	},
	cvmtypes.StoreKey: {
		{cvmtypes.StorageStoreKeyPrefix,         nil, 4}, //raw []byte
		{cvmtypes.BlockHashStoreKeyPrefix,       nil, 4}, //hash []byte
		{cvmtypes.CodeStoreKeyPrefix,            &cvmtypes.CVMCode{}, 1},
		{cvmtypes.AbiStoreKeyPrefix,             nil, 4}, //string
		{cvmtypes.MetaHashStoreKeyPrefix,        nil, 4}, //string
		{cvmtypes.AddressMetaHashStoreKeyPrefix, &cvmtypes.ContractMetas{}, 1},
	},
	oracletypes.StoreKey: {
		{oracletypes.OperatorStoreKeyPrefix,         &oracletypes.Operator{}, 2},
		{oracletypes.WithdrawStoreKeyPrefix,         &oracletypes.Withdraw{}, 2},
		{oracletypes.TotalCollateralKeyPrefix,       &oracletypes.CoinsProto{}, 2},
		{oracletypes.TaskStoreKeyPrefix,             &ti, 3},
		{oracletypes.ClosingTaskStoreKeyPrefix,      &oracletypes.TaskIDs{}, 2},
		{oracletypes.ClosingTaskStoreKeyTimedPrefix, &oracletypes.TaskIDs{}, 2},
		{oracletypes.ExpireTaskStoreKeyPrefix,       &oracletypes.TaskIDs{}, 2},
		{oracletypes.ShortcutTasksKeyPrefix,         &oracletypes.TaskIDs{}, 2},
	},
	shieldtypes.StoreKey: {
		{shieldtypes.ShieldAdminKey,          nil, 4}, //sdk.AccAddress
		{shieldtypes.TotalCollateralKey,      &sdk.IntProto{}, 2},
		{shieldtypes.TotalWithdrawingKey,     &sdk.IntProto{}, 2},
		{shieldtypes.TotalShieldKey,          &sdk.IntProto{}, 2},
		{shieldtypes.TotalClaimedKey,         &sdk.IntProto{}, 2},
		{shieldtypes.ServiceFeesKey,          &shieldtypes.Fees{}, 2},
		{shieldtypes.RemainingServiceFeesKey, &shieldtypes.Fees{}, 2},
		{shieldtypes.PoolKey,                 &shieldtypes.Pool{}, 2},
		{shieldtypes.NextPoolIDKey,           nil, 4}, //binary.LittleEndian.PutUint64(bz, id)
		{shieldtypes.NextPurchaseIDKey,       nil, 4}, //binary.LittleEndian.PutUint64(bz, id)
		{shieldtypes.PurchaseListKey,         &shieldtypes.PurchaseList{}, 2},
		{shieldtypes.PurchaseQueueKey,        &shieldtypes.PoolPurchaserPairs{}, 2},
		{shieldtypes.ProviderKey,             &shieldtypes.Provider{}, 2},
		{shieldtypes.WithdrawQueueKey,        &shieldtypes.Withdraws{}, 2},
		{shieldtypes.LastUpdateTimeKey,       &shieldtypes.LastUpdateTime{}, 2},
		{shieldtypes.GlobalStakeForShieldPoolKey, &sdk.IntProto{}, 2},
		{shieldtypes.StakeForShieldKey,       &shieldtypes.ShieldStaking{}, 2},
		{shieldtypes.BlockServiceFeesKey,     &shieldtypes.Fees{}, 2},
		{shieldtypes.OriginalStakingKey,      &sdk.IntProto{}, 2},
		{shieldtypes.ReimbursementKey,        &shieldtypes.Reimbursement{}, 2},
	},
	authtypes.StoreKey: {
		{authtypes.AddressStoreKeyPrefix,  &ai, 3},
		{authtypes.GlobalAccountNumberKey, &gogotypes.UInt64Value{}, 1},
	},
	authztypes.StoreKey: {
		{authztypes.GrantKey, &authz.Grant{}, 1},
	},
	banktypes.StoreKey: {
		{banktypes.BalancesPrefix,      &sdk.Coin{}, 1},
		{banktypes.SupplyKey,           nil, 4}, //sdk.Int
		{banktypes.DenomMetadataPrefix, &banktypes.Metadata{}, 1},
	},
	captypes.StoreKey: {
		{captypes.KeyIndex,                 nil, 4}, //binary.BigEndian.PutUint64
		{captypes.KeyPrefixIndexCapability, &captypes.CapabilityOwners{}, 1},
		{captypes.KeyMemInitialized,        nil, 4}, //[]byte{1}
	},
	// crisistypes: {}
	disttypes.StoreKey: {
		{disttypes.FeePoolKey,  &disttypes.FeePool{}, 1},
		{disttypes.ProposerKey, &gogotypes.BytesValue{}, 1},
		{disttypes.ValidatorOutstandingRewardsPrefix,    &disttypes.ValidatorOutstandingRewards{}, 1},
		{disttypes.DelegatorWithdrawAddrPrefix,          nil, 4}, //sdk.AccAddress
		{disttypes.DelegatorStartingInfoPrefix,          &disttypes.DelegatorStartingInfo{}, 1},
		{disttypes.ValidatorHistoricalRewardsPrefix,     &disttypes.ValidatorHistoricalRewards{}, 1},
		{disttypes.ValidatorCurrentRewardsPrefix,        &disttypes.ValidatorCurrentRewards{}, 1},
		{disttypes.ValidatorAccumulatedCommissionPrefix, &disttypes.ValidatorAccumulatedCommission{}, 1},
		{disttypes.ValidatorSlashEventPrefix,            &disttypes.ValidatorSlashEvent{}, 1},
	},
	evidtypes.StoreKey: {
		{evidtypes.KeyPrefixEvidence, &evidi, 3},
	},
	feegrant.StoreKey: {
		{feegrant.FeeAllowanceKeyPrefix, &feegrant.Grant{}, 1},
	},
	govtypes.StoreKey: {
		{govtypes.ProposalsKeyPrefix,          &govtypes.Proposal{}, 1},
		{govtypes.ActiveProposalQueuePrefix,   nil, 4}, //binary.BigEndian.PutUint64
		{govtypes.InactiveProposalQueuePrefix, nil, 4}, //binary.BigEndian.PutUint64
		{govtypes.ProposalIDKey,               nil, 4}, //binary.BigEndian.PutUint64
		{govtypes.DepositsKeyPrefix,           &govtypes.Deposit{}, 1},
		{govtypes.VotesKeyPrefix,              &govtypes.Vote{}, 1},
		{stgovtypes.CertVotesKeyPrefix,        nil, 4}, // managed by shentu/gov, binary.BigEndian.PutUint64
	},
	minttypes.StoreKey: {
		{minttypes.MinterKey, &minttypes.Minter{}, 1},
	},
	// params //check if it's possible to dump all params managed by param module
	slashingtypes.StoreKey: {
		{slashingtypes.ValidatorSigningInfoKeyPrefix, &slashingtypes.ValidatorSigningInfo{}, 1},
		{slashingtypes.ValidatorMissedBlockBitArrayKeyPrefix, &gogotypes.BoolValue{}, 1},
		{slashingtypes.AddrPubkeyRelationKeyPrefix, &pki, 3},
	},
	stakingtypes.StoreKey: {
		{stakingtypes.LastValidatorPowerKey, &gogotypes.Int64Value{}, 1},
		{stakingtypes.LastTotalPowerKey,     &sdk.IntProto{}, 1},
		{stakingtypes.ValidatorsKey,         &stakingtypes.Validator{}, 1},
		{stakingtypes.ValidatorsByConsAddrKey,   nil, 4}, //sdk.ValAddress
		{stakingtypes.ValidatorsByPowerIndexKey, nil, 4}, //sdk.ValAddress
		{stakingtypes.DelegationKey,          &stakingtypes.Delegation{}, 1},
		{stakingtypes.UnbondingDelegationKey, &stakingtypes.UnbondingDelegation{}, 1},
		{stakingtypes.UnbondingDelegationByValIndexKey, nil, 4}, //[]byte
		{stakingtypes.RedelegationKey,      &stakingtypes.Redelegation{}, 1},
		{stakingtypes.RedelegationByValSrcIndexKey, nil, 4}, //[]byte
		{stakingtypes.RedelegationByValDstIndexKey, nil, 4}, //[]byte
		{stakingtypes.UnbondingQueueKey,    &stakingtypes.DVPairs{}, 1},
		{stakingtypes.RedelegationQueueKey, &stakingtypes.DVVTriplets{}, 1},
		{stakingtypes.ValidatorQueueKey,    &stakingtypes.ValAddresses{}, 1},
		{stakingtypes.HistoricalInfoKey,    &stakingtypes.HistoricalInfo{}, 1},
	},
	upgradetypes.StoreKey: {
		{[]byte{upgradetypes.PlanByte},            &upgradetypes.Plan{}, 1},
		{[]byte{upgradetypes.DoneByte},            nil, 4}, //binary.BigEndian.PutUint64
		{[]byte{upgradetypes.VersionMapByte},      nil, 4}, //binary.BigEndian.PutUint64
		{[]byte{upgradetypes.ProtocolVersionByte}, nil, 4}, //binary.BigEndian.PutUint64
		{[]byte(upgradetypes.KeyUpgradedIBCState), nil, 4}, //[]byte
		// {[]byte{upgradetypes.KeyUpgradedClient}, },
		// {[]byte{upgradetypes.KeyUpgradedConsState}, },
	},
	ibctransfertypes.StoreKey: {
		{ibctransfertypes.PortKey,       nil, 4},
		{ibctransfertypes.DenomTraceKey, &ibctransfertypes.DenomTrace{}, 1},
	},
	icahosttypes.StoreKey: {
		{[]byte(icatypes.ActiveChannelKeyPrefix), nil, 4},
		{[]byte(icatypes.OwnerKeyPrefix),         nil, 4},
		{[]byte(icatypes.PortKeyPrefix),          nil, 4},
	},
}

var pathKeys = KeyTypes {
	ibchost.StoreKey: {
		//sample: clients/07-tendermint-9/clientState
		{[]byte("clients/[^/]+/clientState$"), &ibcCS, 3},
		//sample: clients/07-tendermint-9/connections
		{[]byte("clients/[^/]+/connections$"), &ibcconntypes.ClientPaths{}, 1},
		//sample: clients/07-tendermint-9/consensusStates/1-11532270
		{[]byte("clients/[^/]+/consensusStates/[^/]+[0-9]+$"), &ibcConsS, 3},
		//sample: clients/07-tendermint-9/consensusStates/1-11504960/processedHeight
		{[]byte("clients/[^/]+/consensusStates/[^/]+[0-9]+/processedHeight$"), nil, 4}, //[]byte(exported.Height.String())
		//sample: clients/07-tendermint-9/consensusStates/1-11504960/processedTime
		{[]byte("clients/[^/]+/consensusStates/[^/]+[0-9]+/processedTime$"), nil, 4},  //sdk.Uint64ToBigEndian(uint64 timeNs)
		//sample: clients/07-tendermint-9/iterateConsensusStates[00][00][00][00][00][00][00][01][00][00][00][00][00][ad][02][85]
		{[]byte("(?s)clients/[^/]+/iterateConsensusStates.+$"), nil, 4},  //[]byte(ConsensusStatePath(height))
		//sample: receipts/ports/transfer/channels/channel-8/sequences/88
		{[]byte("receipts/.+$"),         nil, 4},  //[]byte{byte(1)}
		//sample: nextSequenceSend/ports/transfer/channels/channel-9
		{[]byte("nextSequenceSend/.+$"), nil, 4},  //sdk.Uint64ToBigEndian(uint64)
		//sample: connections/connection-18
		{[]byte("connections/.+$"), &ibcconntypes.ConnectionEnd{}, 1},
		//sample: nextSequenceAck/ports/transfer/channels/channel-12
		{[]byte("nextSequenceAck/.+$"),  nil, 4},  //sdk.Uint64ToBigEndian(uint64)
		//sample: nextSequenceRecv/ports/transfer/channels/channel-11
		{[]byte("nextSequenceRecv/.+$"), nil, 4},  // sdk.Uint64ToBigEndian(uint64)
		//sample: acks/ports/transfer/channels/channel-1/sequences/28
		{[]byte("acks/.+$"),             nil, 4},  // []byte
		//sample: channelEnds/ports/transfer/channels/channel-15
		{[]byte("channelEnds/.+$"),         &ibcchantypes.Channel{}, 1},
		//sample: commitments/ports/transfer/channels/channel-11/sequences/1
		{[]byte("commitments/.+$"),         nil, 4}, //sha256.Sum256(buf)
		//sample: nextChannelSequence
		{[]byte("nextChannelSequence$"),    nil, 4}, //sdk.Uint64ToBigEndian(sequence)
		//sample: nextClientSequence
		{[]byte("nextClientSequence$"),     nil, 4}, //sdk.Uint64ToBigEndian(sequence)
		//sample: nextConnectionSequence
		{[]byte("nextConnectionSequence$"), nil, 4}, //sdk.Uint64ToBigEndian(sequence)
	},
}
