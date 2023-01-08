package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
)

const (
	ErrorEmptyProgramIDFindingList = "empty finding id list"
)

func (k Keeper) GetFinding(ctx sdk.Context, id uint64) (types.Finding, bool) {
	store := ctx.KVStore(k.storeKey)

	pBz := store.Get(types.GetFindingKey(id))
	if pBz == nil {
		return types.Finding{}, false
	}

	var finding types.Finding
	k.cdc.MustUnmarshal(pBz, &finding)
	return finding, true
}

func (k Keeper) SetFinding(ctx sdk.Context, finding types.Finding) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&finding)
	store.Set(types.GetFindingKey(finding.FindingId), bz)
}

func (k Keeper) GetNextFindingID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	Bz := store.Get(types.GetNextFindingIDKey())
	if Bz == nil {
		return 1
	}
	return binary.LittleEndian.Uint64(Bz)
}

func (k Keeper) SetNextFindingID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.GetNextFindingIDKey(), bz)
}

func (k Keeper) SetPidFindingIDList(ctx sdk.Context, pid uint64, findingIds []uint64) error {
	findingIdList, err := Uint64sToBytes(findingIds)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetProgramIDFindingListKey(pid), findingIdList)
	return nil
}

func (k Keeper) GetPidFindingIDList(ctx sdk.Context, pid uint64) ([]uint64, error) {
	store := ctx.KVStore(k.storeKey)
	findingIDs := store.Get(types.GetProgramIDFindingListKey(pid))

	if findingIDs == nil {
		return nil, fmt.Errorf(ErrorEmptyProgramIDFindingList)
	}

	findingIDList, err := BytesToUint64s(findingIDs)
	if err != nil {
		return nil, err
	}
	return findingIDList, nil
}

func (k Keeper) AppendFidToFidList(ctx sdk.Context, pid, fid uint64) error {
	fids, err := k.GetPidFindingIDList(ctx, pid)
	if err.Error() == ErrorEmptyProgramIDFindingList {
		fids = []uint64{}
	} else if err != nil {
		return err
	}

	fids = append(fids, fid)
	err = k.SetPidFindingIDList(ctx, pid, fids)
	return err
}

func Uint64sToBytes(list []uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, list)
	if err != nil {
		return nil, fmt.Errorf("convert uint64 to byte list error")
	}
	return buf.Bytes(), nil
}

func BytesToUint64s(list []byte) ([]uint64, error) {
	buf := bytes.NewBuffer(list)
	r64 := make([]uint64, (len(list)+7)/8)
	err := binary.Read(buf, binary.LittleEndian, &r64)
	if err != nil {
		return nil, fmt.Errorf("convert to uint64 list error")
	}
	return r64, nil
}
