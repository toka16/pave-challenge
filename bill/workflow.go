package bill

import (
	"errors"
	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, bill BillState) error {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, "getBill", func(_ any) (BillState, error) {
		return bill, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	err = workflow.SetUpdateHandlerWithOptions(ctx, UPDATE_NAME_MODIFY_ITEMS, func(ctx workflow.Context, action string, item ModifyItemData) (BillState, error) {
		logger.Info("Received update, ", action, item)
		switch action {
		case ACTION_ADD_ITEM:
			bill.AddItem(item.Item)
		case ACTION_REMOVE_ITEM:
			bill.RemoveItem(item.Item.ItemID)
		default:
			return bill, errors.New("unknown action")
		}
		return bill, nil
	}, workflow.UpdateHandlerOptions{
		Validator: func(ctx workflow.Context, action string, item ModifyItemData) error {
			logger.Info("Validating update, ", action, item)
			if bill.Status == BillStatusClosed {
				return errors.New("bill is closed")
			}
			if item.Item.ItemID == "" {
				return errors.New("item id cannot be empty")
			}
			if action == ACTION_ADD_ITEM {
				if item.Item.Name == "" {
					return errors.New("name cannot be empty")
				}
				if item.Item.Price <= 0 {
					return errors.New("price must be greater than 0")
				}
			}
			return nil
		},
	})
	if err != nil {
		logger.Error("SetUpdateHandler failed.", "Error", err)
		return err
	}

	closeChannel := workflow.GetSignalChannel(ctx, CHANNEL_CLOSE)
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(closeChannel, func(c workflow.ReceiveChannel, _ bool) {
		bill.Status = BillStatusClosed
	})
	for {
		selector.Select(ctx)

		if bill.Status == BillStatusClosed {
			break
		}
	}

	return nil

}
