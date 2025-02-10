package bill

import (
	"go.temporal.io/sdk/workflow"
)

func BillWorkflow(ctx workflow.Context, bill BillState) error {
	logger := workflow.GetLogger(ctx)

	err := workflow.SetQueryHandler(ctx, "getBill", func(input []byte) (BillState, error) {
		return bill, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	addToCartChannel := workflow.GetSignalChannel(ctx, ActionChannels.ADD_ITEM_CHANNEL)
	removeFromCartChannel := workflow.GetSignalChannel(ctx, ActionChannels.REMOVE_ITEM_CHANNEL)
	closeChannel := workflow.GetSignalChannel(ctx, ActionChannels.CLOSE_CHANNEL)

	for {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(addToCartChannel, func(c workflow.ReceiveChannel, _ bool) {
			var message AddItemSignal
			c.Receive(ctx, &message)

			bill.AddItem(message.Item)
		})

		selector.AddReceive(removeFromCartChannel, func(c workflow.ReceiveChannel, _ bool) {
			var message RemoveItemSignal
			c.Receive(ctx, &message)

			bill.RemoveItem(message.ItemID)
		})

		selector.AddReceive(closeChannel, func(c workflow.ReceiveChannel, _ bool) {
			bill.Status = BillStatusClosed
		})

		selector.Select(ctx)

		if bill.Status == BillStatusClosed {
			break
		}
	}

	return nil

}
