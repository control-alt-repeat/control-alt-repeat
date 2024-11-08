package internal

import (
	"context"
	"fmt"

	"github.com/control-alt-repeat/control-alt-repeat/internal/labels"
	"github.com/control-alt-repeat/control-alt-repeat/internal/logger"
	"github.com/control-alt-repeat/control-alt-repeat/internal/warehouse"
)

type ItemPrintShelfLabelOptions struct {
	ItemID string
}

func ItemPrintShelfLabel(ctx context.Context, opts ItemPrintShelfLabelOptions) error {
	log := logger.Instance
	log.Info().Fields(opts).Msg("Loading item from warehouse")

	item, exists, err := warehouse.LoadItem(ctx, opts.ItemID)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("item does not exist for ID '%s'", opts.ItemID)
	}

	log.Info().Msg("Generating label")
	log.Debug().Fields(item).Msg("")
	label, name, err := labels.CreateShelfLabelFromItem(ctx, item)
	if err != nil {
		return err
	}

	log.Info().Msgf("Sending label to the label printer as '%s'", name)
	return labels.UploadFileFromBytes(label, name)
}
