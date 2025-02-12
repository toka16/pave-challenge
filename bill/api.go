package bill

import (
	"context"
	"encore.app/bill/workflow"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"fmt"
	"go.temporal.io/sdk/client"
	"time"
)

//encore:api public method=POST path=/bill
func (s *Service) CreateBill(ctx context.Context, req CreateBillRequest) (*BillResponse, error) {
	billID := "BILL-" + fmt.Sprintf("%d", time.Now().Unix())

	options := client.StartWorkflowOptions{
		ID:        billID,
		TaskQueue: TaskQueueName,
	}

	bill := workflow.BillState{ID: billID, Currency: req.Currency, Items: make([]workflow.BillItem, 0), Status: workflow.BillStatusOpen}
	_, err := s.client.ExecuteWorkflow(ctx, options, workflow.Workflow, bill)
	if err != nil {
		rlog.Error("Error executing workflow", err)
		return nil, &errs.Error{Code: errs.Internal, Message: err.Error()}
	}
	rlog.Debug("Bill created", rlog.With("bill", bill))

	return BillResponseFromState(bill), nil
}

//encore:api public method=GET path=/bill/:billID
func (s *Service) QueryBill(ctx context.Context, billID string) (*BillResponse, error) {
	response, err := s.client.QueryWorkflow(ctx, billID, "", "getBill")
	if err != nil {
		rlog.Error("Error querying workflow", err)
		return nil, &errs.Error{Code: errs.Internal, Message: err.Error()}
	}
	var res workflow.BillState
	if err = response.Get(&res); err != nil {
		rlog.Error("Error getting workflow state", err)
		return nil, &errs.Error{Code: errs.Internal, Message: err.Error()}
	}
	rlog.Debug("Bill queried", rlog.With("bill", res))
	return BillResponseFromState(res), nil
}

//encore:api public method=POST path=/bill/:billID/close
func (s *Service) CloseBill(ctx context.Context, billID string) error {
	err := s.client.SignalWorkflow(ctx, billID, "", workflow.CHANNEL_CLOSE, nil)
	if err != nil {
		rlog.Error("Error signaling workflow", err)
		return &errs.Error{Code: errs.Internal, Message: err.Error()}
	}
	return nil
}

//encore:api public method=POST path=/bill/:billID/items
func (s *Service) AddItem(ctx context.Context, billID string, item BillItemDTO) error {
	payload := workflow.ModifyItemData{Item: workflow.BillItem{
		ItemID: item.ItemID,
		Name:   item.Name,
		Price:  item.Price,
	}}
	updateOptions := client.UpdateWorkflowOptions{
		WorkflowID:   billID,
		UpdateName:   workflow.UPDATE_NAME_MODIFY_ITEMS,
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args:         []interface{}{workflow.ACTION_ADD_ITEM, payload},
	}
	handle, err := s.client.UpdateWorkflow(ctx, updateOptions)
	if err != nil {
		rlog.Error("Error updating workflow", err)
		return &errs.Error{Code: errs.Internal, Message: err.Error()}
	}
	billState := workflow.BillState{Items: make([]workflow.BillItem, 0)}
	err = handle.Get(ctx, &billState)
	if err != nil {
		rlog.Error("Error getting workflow state", err)
		return &errs.Error{Code: errs.InvalidArgument, Message: err.Error()}
	}
	rlog.Debug("Item added", rlog.With("bill", billState))

	return nil
}

//encore:api public method=DELETE path=/bill/:billID/items/:itemID
func (s *Service) RemoveItem(ctx context.Context, billID string, itemID string) error {
	payload := workflow.ModifyItemData{Item: workflow.BillItem{
		ItemID: itemID,
	}}
	updateOptions := client.UpdateWorkflowOptions{
		WorkflowID:   billID,
		UpdateName:   workflow.UPDATE_NAME_MODIFY_ITEMS,
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args:         []interface{}{workflow.ACTION_REMOVE_ITEM, payload},
	}
	handle, err := s.client.UpdateWorkflow(ctx, updateOptions)
	if err != nil {
		rlog.Error("Error updating workflow", err)
		return &errs.Error{Code: errs.Internal, Message: err.Error()}
	}
	billState := workflow.BillState{Items: make([]workflow.BillItem, 0)}
	err = handle.Get(ctx, &billState)
	if err != nil {
		rlog.Error("Error getting workflow state", err)
		return &errs.Error{Code: errs.InvalidArgument, Message: err.Error()}
	}
	rlog.Debug("Item removed", rlog.With("bill", billState))

	return nil
}
