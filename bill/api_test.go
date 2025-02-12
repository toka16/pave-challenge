package bill

import (
	"context"
	"encore.app/bill/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/mocks"
	"testing"
)

func TestCreateBill(t *testing.T) {
	clientMock := mocks.NewClient(t)
	clientMock.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		opts := args.Get(1).(client.StartWorkflowOptions)
		assert.Equal(t, "BILL-", opts.ID[:5])

		bill := args.Get(3).(workflow.BillState)
		assert.Equal(t, opts.ID, bill.ID)
		assert.Equal(t, workflow.BillStatusOpen, bill.Status)
		assert.Equal(t, "GEL", bill.Currency)
		assert.Len(t, bill.Items, 0)
	}).Return(nil, nil)
	s := Service{client: clientMock}

	res, err := s.CreateBill(context.Background(), CreateBillRequest{Currency: "GEL"})
	assert.NoError(t, err)
	assert.Equal(t, string(workflow.BillStatusOpen), res.Status)
	assert.Equal(t, "GEL", res.Currency)
	assert.Equal(t, 0.0, res.Sum)
}

func TestQueryBill(t *testing.T) {
	clientMock := mocks.NewClient(t)
	val := &mocks.Value{}
	val.On("Get", mock.Anything).Run(func(args mock.Arguments) {
		ptr := args.Get(0).(*workflow.BillState)
		*ptr = workflow.BillState{ID: "BILL-123", Currency: "USD", Items: make([]workflow.BillItem, 0), Status: workflow.BillStatusOpen}
	}).Return(nil)
	clientMock.On("QueryWorkflow", mock.Anything, "BILL-123", "", "getBill").Return(val, nil)
	s := Service{client: clientMock}

	res, err := s.QueryBill(context.Background(), "BILL-123")
	assert.NoError(t, err)
	assert.Equal(t, string(workflow.BillStatusOpen), res.Status)
	assert.Equal(t, "USD", res.Currency)
	assert.Equal(t, 0.0, res.Sum)
}

func TestCloseBill(t *testing.T) {
	clientMock := mocks.NewClient(t)
	clientMock.On("SignalWorkflow", mock.Anything, "BILL-123", "", workflow.CHANNEL_CLOSE, nil).Return(nil)
	s := Service{client: clientMock}

	err := s.CloseBill(context.Background(), "BILL-123")
	assert.NoError(t, err)
}

func TestAddItem(t *testing.T) {
	clientMock := mocks.NewClient(t)
	handle := &mocks.WorkflowUpdateHandle{}
	handle.On("Get", mock.Anything, mock.Anything).Return(nil)
	clientMock.On("UpdateWorkflow", mock.Anything, client.UpdateWorkflowOptions{
		WorkflowID:   "BILL-123",
		UpdateName:   workflow.UPDATE_NAME_MODIFY_ITEMS,
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args: []interface{}{
			workflow.ACTION_ADD_ITEM,
			workflow.ModifyItemData{
				Item: workflow.BillItem{
					ItemID: "ITEM-123",
					Name:   "item name",
					Price:  2.0,
				},
			},
		},
	}).Return(handle, nil)
	s := Service{client: clientMock}

	err := s.AddItem(context.Background(), "BILL-123", BillItemDTO{ItemID: "ITEM-123", Name: "item name", Price: 2.0})
	assert.NoError(t, err)
}

func TestRemoveItem(t *testing.T) {
	clientMock := mocks.NewClient(t)
	handle := &mocks.WorkflowUpdateHandle{}
	handle.On("Get", mock.Anything, mock.Anything).Return(nil)
	clientMock.On("UpdateWorkflow", mock.Anything, client.UpdateWorkflowOptions{
		WorkflowID:   "BILL-123",
		UpdateName:   workflow.UPDATE_NAME_MODIFY_ITEMS,
		WaitForStage: client.WorkflowUpdateStageCompleted,
		Args: []interface{}{
			workflow.ACTION_REMOVE_ITEM,
			workflow.ModifyItemData{
				Item: workflow.BillItem{
					ItemID: "ITEM-123",
				},
			},
		},
	}).Return(handle, nil)
	s := Service{client: clientMock}

	err := s.RemoveItem(context.Background(), "BILL-123", "ITEM-123")
	assert.NoError(t, err)
}
