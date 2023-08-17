package cmd

import (
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
)

var allKeys = map[string][]OneKey {
	certtypes.StoreKey: {
		{certtypes.CertifiersStoreKey(),       &certtypes.Certifier{}, 2},
		{certtypes.CertifierAliasesStoreKey(), &certtypes.Certifier{}, 2},
		{certtypes.PlatformsStoreKey(),        &certtypes.Platform{}, 1},
		{certtypes.CertificatesStoreKey(),     &certtypes.Certificate{}, 1},
		{certtypes.LibrariesStoreKey(),        &certtypes.Library{}, 2},
		{certtypes.NextCertificateIDStoreKey(), nil, 4}, //binary.LittleEndian.Uint64
		{certtypes.ValidatorsStoreKey(), nil, 4}, //this key is not used
	},
	cvmtypes.StoreKey: {
		{cvmtypes.StorageStoreKeyPrefix, nil, 4}, //raw []byte
		{cvmtypes.BlockHashStoreKeyPrefix, nil, 4}, //hash []byte
		{cvmtypes.CodeStoreKeyPrefix, &cvmtypes.CVMCode{}, 1},
		{cvmtypes.AbiStoreKeyPrefix, nil, 4}, //string
		{cvmtypes.MetaHashStoreKeyPrefix, nil, 4}, //string
		{cvmtypes.AddressMetaHashStoreKeyPrefix, &cvmtypes.ContractMetas{}, 1},
	},
	oracletypes.StoreKey: {
		{oracletypes.OperatorStoreKeyPrefix,   &oracletypes.Operator{}, 2},
		{oracletypes.WithdrawStoreKeyPrefix,   &oracletypes.Withdraw{}, 2},
		{oracletypes.TotalCollateralKeyPrefix, &oracletypes.CoinsProto{}, 2},
		{oracletypes.TaskStoreKeyPrefix,       &ti, 3},
		{oracletypes.ClosingTaskStoreKeyPrefix, &oracletypes.TaskIDs{}, 2},
		{oracletypes.ClosingTaskStoreKeyTimedPrefix, &oracletypes.TaskIDs{}, 2},
		{oracletypes.ExpireTaskStoreKeyPrefix, &oracletypes.TaskIDs{}, 2},
		{oracletypes.ShortcutTasksKeyPrefix,   &oracletypes.TaskIDs{}, 2},
	},
	shieldtypes.StoreKey: {
		{shieldtypes.ShieldAdminKey, nil, 4}, //sdk.AccAddress
		{shieldtypes.TotalCollateralKey, &sdk.IntProto{}, 2},
		{shieldtypes.TotalWithdrawingKey, &sdk.IntProto{}, 2},
		{shieldtypes.TotalShieldKey, &sdk.IntProto{}, 2},
		{shieldtypes.TotalClaimedKey, &sdk.IntProto{}, 2},
		{shieldtypes.ServiceFeesKey, &shieldtypes.Fees{}, 2},
		{shieldtypes.RemainingServiceFeesKey, &shieldtypes.Fees{}, 2},
		{shieldtypes.PoolKey,       &shieldtypes.Pool{}, 2},
		{shieldtypes.NextPoolIDKey, nil, 4}, //binary.LittleEndian.PutUint64(bz, id)
		{shieldtypes.NextPurchaseIDKey, nil, 4}, //binary.LittleEndian.PutUint64(bz, id)
		{shieldtypes.PurchaseListKey,   &shieldtypes.PurchaseList{}, 2},
		{shieldtypes.PurchaseQueueKey,  &shieldtypes.PoolPurchaserPairs{}, 2},
		{shieldtypes.ProviderKey,       &shieldtypes.Provider{}, 2},
		{shieldtypes.WithdrawQueueKey,  &shieldtypes.Withdraws{}, 2},
		{shieldtypes.LastUpdateTimeKey, &shieldtypes.LastUpdateTime{}, 2},
		{shieldtypes.GlobalStakeForShieldPoolKey, &sdk.IntProto{}, 2},
		{shieldtypes.StakeForShieldKey,   &shieldtypes.ShieldStaking{}, 2},
		{shieldtypes.BlockServiceFeesKey, &shieldtypes.Fees{}, 2},
		{shieldtypes.OriginalStakingKey,  &sdk.IntProto{}, 2},
		{shieldtypes.ReimbursementKey,    &shieldtypes.Reimbursement{}, 2},
	},
	authtypes.StoreKey: {
		{authtypes.AddressStoreKeyPrefix, &ai, 3},
		{authtypes.GlobalAccountNumberKey, &gogotypes.UInt64Value{}, 1},
	},
	authztypes.StoreKey: {
		{authztypes.GrantKey, &authz.Grant{}, 1},
	},
	banktypes.StoreKey: {
		{banktypes.BalancesPrefix, &sdk.Coin{}, 1},
		{banktypes.SupplyKey, nil, 4}, //sdk.Int
		{banktypes.DenomMetadataPrefix, &banktypes.Metadata{}, 1},
	},
	captypes.StoreKey: {
		{captypes.KeyIndex, nil, 4}, //binary.BigEndian.PutUint64
		{captypes.KeyPrefixIndexCapability, &captypes.CapabilityOwners{}, 1},
		{captypes.KeyMemInitialized, nil, 4}, //[]byte{1}
	},
	// crisistypes: {}
	disttypes.StoreKey: {
		{disttypes.FeePoolKey, &disttypes.FeePool{}, 1},
		{disttypes.ProposerKey, &gogotypes.BytesValue{}, 1},
		{disttypes.ValidatorOutstandingRewardsPrefix, &disttypes.ValidatorOutstandingRewards{}, 1},
		{disttypes.DelegatorWithdrawAddrPrefix, nil, 4}, //sdk.AccAddress
		{disttypes.DelegatorStartingInfoPrefix, &disttypes.DelegatorStartingInfo{}, 1},
		{disttypes.ValidatorHistoricalRewardsPrefix, &disttypes.ValidatorHistoricalRewards{}, 1},
		{disttypes.ValidatorCurrentRewardsPrefix, &disttypes.ValidatorCurrentRewards{}, 1},
		{disttypes.ValidatorAccumulatedCommissionPrefix, &disttypes.ValidatorAccumulatedCommission{}, 1},
		{disttypes.ValidatorSlashEventPrefix, &disttypes.ValidatorSlashEvent{}, 1},
	},
	evidtypes.StoreKey: {
		{evidtypes.KeyPrefixEvidence, &evidi, 3},
	},
	feegrant.StoreKey: {
		{feegrant.FeeAllowanceKeyPrefix, &feegrant.Grant{}, 1},
	},
	govtypes.StoreKey: {
		{govtypes.ProposalsKeyPrefix, &govtypes.Proposal{}, 1},
		{govtypes.ActiveProposalQueuePrefix, nil, 4}, //binary.BigEndian.PutUint64
		{govtypes.InactiveProposalQueuePrefix, nil, 4}, //binary.BigEndian.PutUint64
		{govtypes.ProposalIDKey, nil, 4}, //binary.BigEndian.PutUint64
		{govtypes.DepositsKeyPrefix, &govtypes.Deposit{}, 1},
		{govtypes.VotesKeyPrefix, &govtypes.Vote{}, 1},
		{stgovtypes.CertVotesKeyPrefix, nil, 4}, // managed by shentu/gov, binary.BigEndian.PutUint64
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
		{stakingtypes.ValidatorsByConsAddrKey, nil, 4}, //sdk.ValAddress
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
		{[]byte{upgradetypes.PlanByte}, &upgradetypes.Plan{}, 1},
		{[]byte{upgradetypes.DoneByte}, nil, 4}, //binary.BigEndian.PutUint64
		{[]byte{upgradetypes.VersionMapByte}, nil, 4}, //binary.BigEndian.PutUint64
		{[]byte{upgradetypes.ProtocolVersionByte}, nil, 4}, //binary.BigEndian.PutUint64
		{[]byte(upgradetypes.KeyUpgradedIBCState), nil, 4}, //[]byte
		// {[]byte{upgradetypes.KeyUpgradedClient}, },
		// {[]byte{upgradetypes.KeyUpgradedConsState}, },
	},
}