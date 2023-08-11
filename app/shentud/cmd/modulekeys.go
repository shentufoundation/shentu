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
)

var allKeys = map[string][]OneKey {
	certtypes.StoreKey: {
		{certtypes.CertifiersStoreKey(),       &certtypes.Certifier{}, 2},
		{certtypes.CertifierAliasesStoreKey(), &certtypes.Certifier{}, 2},
		{certtypes.PlatformsStoreKey(),        &certtypes.Platform{}, 1},
		{certtypes.CertificatesStoreKey(),     &certtypes.Certificate{}, 1},
		{certtypes.LibrariesStoreKey(),        &certtypes.Library{}, 2},
		// {certtypes.NextCertificateIDStoreKey(), &uint64(), 1}, //binary.LittleEndian.Uint64
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
		// {banktypes.SupplyKey, &sdk.Int{}, 1},
		{banktypes.DenomMetadataPrefix, &banktypes.Metadata{}, 1},
	},
	captypes.StoreKey: {
		// {captypes.KeyIndex, uint64(), 1},
		{captypes.KeyPrefixIndexCapability, &captypes.CapabilityOwners{}, 1},
		// {captypes.KeyMemInitialized, []byte, 1}
	},
	// crisistypes: {}
	disttypes.StoreKey: {
		{disttypes.FeePoolKey, &disttypes.FeePool{}, 1},
		{disttypes.ProposerKey, &gogotypes.BytesValue{}, 1},
		{disttypes.ValidatorOutstandingRewardsPrefix, &disttypes.ValidatorOutstandingRewards{}, 1},
		// {disttypes.DelegatorWithdrawAddrPrefix, sdk.AccAddress, 1},
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
		// {govtypes.ActiveProposalQueuePrefix, uint64, 1},
		// {govtypes.InactiveProposalQueuePrefix, uint64, 1},
		// {govtypes.ProposalIDKey, uint64, 1},
		{govtypes.DepositsKeyPrefix, &govtypes.Deposit{}, 1},
		{govtypes.VotesKeyPrefix, &govtypes.Vote{}, 1},
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
		// {stakingtypes.ValidatorsByConsAddrKey, sdk.ValAddress, 1},
		// {stakingtypes.ValidatorsByPowerIndexKey, sdk.ValAddress, 1},
		{stakingtypes.DelegationKey,          &stakingtypes.Delegation{}, 1},
		{stakingtypes.UnbondingDelegationKey, &stakingtypes.UnbondingDelegation{}, 1},
		// {stakingtypes.UnbondingDelegationByValIndexKey, []byte{}, 1},
		{stakingtypes.RedelegationKey,      &stakingtypes.Redelegation{}, 1},
		// {stakingtypes.RedelegationByValSrcIndexKey, []byte{}, 1},
		// {stakingtypes.RedelegationByValDstIndexKey, []byte{}, 1},
		{stakingtypes.UnbondingQueueKey,    &stakingtypes.DVPairs{}, 1},
		{stakingtypes.RedelegationQueueKey, &stakingtypes.DVVTriplets{}, 1},
		{stakingtypes.ValidatorQueueKey,    &stakingtypes.ValAddresses{}, 1},
		{stakingtypes.HistoricalInfoKey,    &stakingtypes.HistoricalInfo{}, 1},
	},
	upgradetypes.StoreKey: {
		{[]byte{upgradetypes.PlanByte}, &upgradetypes.Plan{}, 1},
		// {upgradetypes.DoneByte, uint64, 1},
		// {upgradetypes.VersionMapByte, uint64, 1},
		// {upgradetypes.ProtocolVersionByte, uint64, 1},
		// {upgradetypes.KeyUpgradedIBCState, []byte, 1},
		// {upgradetypes.KeyUpgradedClient, },
		// {upgradetypes.KeyUpgradedConsState, },
	},
}