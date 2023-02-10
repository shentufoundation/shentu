package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shentufoundation/shentu/v2/x/bounty/client/cli"
	"github.com/shentufoundation/shentu/v2/x/bounty/types"
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

func (k Keeper) DeleteFinding(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetFindingKey(id))
}

func (k Keeper) GetNextFindingID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	Bz := store.Get(types.GetNextFindingIDKey())
	return binary.LittleEndian.Uint64(Bz)
}

func (k Keeper) SetNextFindingID(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	store.Set(types.GetNextFindingIDKey(), bz)
}

func (k Keeper) SetPidFindingIDList(ctx sdk.Context, pid uint64, findingIds []uint64) error {
	findingIDList, err := Uint64sToBytes(findingIds)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetProgramIDFindingListKey(pid), findingIDList)
	return nil
}

func (k Keeper) GetPidFindingIDList(ctx sdk.Context, pid uint64) ([]uint64, error) {
	store := ctx.KVStore(k.storeKey)
	findingIDs := store.Get(types.GetProgramIDFindingListKey(pid))

	if findingIDs == nil {
		return nil, types.ErrProgramFindingListEmpty
	}

	findingIDList, err := BytesToUint64s(findingIDs)
	if err != nil {
		return nil, err
	}
	return findingIDList, nil
}

func (k Keeper) AppendFidToFidList(ctx sdk.Context, pid, fid uint64) error {
	fids, err := k.GetPidFindingIDList(ctx, pid)
	if err != nil {
		if err == types.ErrProgramFindingListEmpty {
			fids = []uint64{}
		} else {
			return err
		}
	}

	fids = append(fids, fid)
	err = k.SetPidFindingIDList(ctx, pid, fids)
	return err
}

func (k Keeper) DeleteFidFromFidList(ctx sdk.Context, pid, fid uint64) error {
	fids, err := k.GetPidFindingIDList(ctx, pid)
	if err != nil {
		return err
	}
	for idx, id := range fids {
		if id == fid {
			if len(fids) == 1 {
				// Delete fid list if empty
				store := ctx.KVStore(k.storeKey)
				store.Delete(types.GetProgramIDFindingListKey(pid))
				return nil
			}
			fids = append(fids[:idx], fids[idx+1:]...)
			return k.SetPidFindingIDList(ctx, pid, fids)
		}
	}
	return types.ErrFindingNotExists
}

func Uint64sToBytes(list []uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, list)
	if err != nil {
		return nil, types.ErrProgramFindingListMarshal
	}
	return buf.Bytes(), nil
}

func BytesToUint64s(list []byte) ([]uint64, error) {
	buf := bytes.NewBuffer(list)
	r64 := make([]uint64, (len(list)+7)/8)
	err := binary.Read(buf, binary.LittleEndian, &r64)
	if err != nil {
		return nil, types.ErrProgramFindingListUnmarshal
	}
	return r64, nil
}

func CheckPlainText(pubKey *ecies.PublicKey, plainText string, anyData *codectypes.Any) (bool, error) {
	var randBytes []byte
	var encryptData []byte

	if anyData.TypeUrl == "/shentu.bounty.v1.EciesEncryptedDesc" {
		var encrypted types.EciesEncryptedDesc
		err := proto.Unmarshal(anyData.GetValue(), &encrypted)
		if err != nil {
			return false, err
		}

		if len(encrypted.FindingDesc) < cli.RandBytesLen {
			return false, types.ErrFindingEncryptedDataInvalid
		}
		randBytesStart := len(encrypted.FindingDesc) - cli.RandBytesLen
		encryptData = encrypted.FindingDesc[:randBytesStart]
		randBytes = encrypted.FindingDesc[randBytesStart:]

	} else if anyData.TypeUrl == "/shentu.bounty.v1.EciesEncryptedPoc" {
		var encrypted types.EciesEncryptedPoc
		err := proto.Unmarshal(anyData.GetValue(), &encrypted)
		if err != nil {
			return false, err
		}

		if len(encrypted.FindingPoc) < cli.RandBytesLen {
			return false, types.ErrFindingEncryptedDataInvalid
		}
		randBytesStart := len(encrypted.FindingPoc) - cli.RandBytesLen
		encryptData = encrypted.FindingPoc[:randBytesStart]
		randBytes = encrypted.FindingPoc[randBytesStart:]
	} else if anyData.TypeUrl == "/shentu.bounty.v1.EciesEncryptedComment" {
		var encrypted types.EciesEncryptedComment
		err := proto.Unmarshal(anyData.GetValue(), &encrypted)
		if err != nil {
			return false, err
		}

		if len(encrypted.FindingComment) < cli.RandBytesLen {
			return false, types.ErrFindingEncryptedDataInvalid
		}
		randBytesStart := len(encrypted.FindingComment) - cli.RandBytesLen
		encryptData = encrypted.FindingComment[:randBytesStart]
		randBytes = encrypted.FindingComment[randBytesStart:]

	} else {
		return false, fmt.Errorf("invalid any data")
	}

	encryptedBytes, err := ecies.Encrypt(bytes.NewReader(randBytes), pubKey, []byte(plainText), nil, nil)
	if err != nil {
		return false, types.ErrProgramPubKey
	}

	if !bytes.Equal(encryptedBytes, encryptData) {
		return false, types.ErrFindingPlainTextDataInvalid
	}

	return true, nil
}
