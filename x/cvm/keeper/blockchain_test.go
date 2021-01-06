package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBlockChain(t *testing.T) {
	testInput := CreateTestInput(t)

	ctx := testInput.Ctx
	cvmk := testInput.CvmKeeper

	bc := NewBlockChain(ctx, cvmk)
	require.NotNil(t, bc.LastBlockTime())
}

func TestBlockchain_BlockHash(t *testing.T) {
	testInput := CreateTestInput(t)

	ctx := testInput.Ctx
	cvmk := testInput.CvmKeeper

	bc := NewBlockChain(ctx, cvmk)
	ctxHash := ctx.BlockHeader().LastBlockId.Hash
	bz, err := bc.BlockHash(0)
	require.Nil(t, err)
	require.Equal(t, ctxHash, bz)

	bz, err = bc.BlockHash(10)
	require.NotNil(t, err)
	require.Nil(t, bz)
}

func TestBlockchain_LastBlockHeight(t *testing.T) {
	testInput := CreateTestInput(t)

	ctx := testInput.Ctx
	cvmk := testInput.CvmKeeper
	bc := NewBlockChain(ctx, cvmk)

	ctxHeight := ctx.BlockHeader().Height
	require.Equal(t, uint64(ctxHeight), bc.LastBlockHeight())
}

func TestBlockchain_LastBlockTime(t *testing.T) {
	testInput := CreateTestInput(t)

	ctx := testInput.Ctx
	cvmk := testInput.CvmKeeper
	bc := NewBlockChain(ctx, cvmk)

	ctxTime := ctx.BlockHeader().Time
	require.Equal(t, ctxTime, bc.LastBlockTime())
}
