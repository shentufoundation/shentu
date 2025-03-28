package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
)

func (k Keeper) DeleteTheorem(ctx context.Context, theoremID uint64) error {
	theorem, err := k.Theorems.Get(ctx, theoremID)
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return status.Errorf(codes.NotFound, "theorem %d doesn't exist", theoremID)
		}
		return err
	}

	err = k.ActiveTheoremsQueue.Remove(ctx, collections.Join(*theorem.EndTime, theorem.Id))
	if err != nil {
		return err
	}
	err = k.Theorems.Remove(ctx, theorem.Id)
	if err != nil {
		return err
	}

	return nil
}
