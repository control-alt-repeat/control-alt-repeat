package internal

import (
	"context"

	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

func MoveItem(ctx context.Context, itemID string, newShelf string) error {
	return warehouse.UpdateShelf(ctx, itemID, newShelf)
}
