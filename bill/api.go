package bill

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"time"
)

type BillResponse struct {
	ID     string        `json:"id"`
	Sum    float64       `json:"sum"`
	Items  []BillItemDTO `json:"items"`
	Status string        `json:"status"`
}
type BillItemDTO struct {
	ItemID string  `json:"item_id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

func BillResponseFromState(state BillState) *BillResponse {
	res := &BillResponse{
		ID:     state.ID,
		Sum:    0,
		Items:  make([]BillItemDTO, 0, len(state.Items)),
		Status: string(state.Status),
	}
	for _, item := range state.Items {
		res.Items = append(res.Items, BillItemDTO{
			ItemID: item.ItemID,
			Name:   item.Name,
			Price:  item.Price,
		})
		res.Sum += item.Price
	}
	return res
}

//encore:api public method=POST path=/bill
func (s *Service) CreateBill(ctx context.Context) (*BillResponse, error) {
	fmt.Println("CreateBill")
	billID := "BILL-" + fmt.Sprintf("%d", time.Now().Unix())

	options := client.StartWorkflowOptions{
		ID:        billID,
		TaskQueue: TaskQueueName,
	}

	bill := BillState{ID: billID, Items: make([]BillItem, 0), Status: BillStatusOpen}
	_, err := s.client.ExecuteWorkflow(context.Background(), options, BillWorkflow, bill)
	if err != nil {
		return nil, err
	}

	return BillResponseFromState(bill), nil
}

//encore:api public method=GET path=/bill/:billID
func (s *Service) QueryBill(ctx context.Context, billID string) (*BillResponse, error) {
	fmt.Println("QueryBill: " + billID)
	response, err := s.client.QueryWorkflow(context.Background(), billID, "", "getBill")
	if err != nil {
		return nil, err
	}
	var res BillState
	if err := response.Get(&res); err != nil {
		return nil, err
	}
	return BillResponseFromState(res), nil
}

//encore:api public method=POST path=/bill/:billID/close
func (s *Service) CloseBill(ctx context.Context, billID string) error {
	fmt.Println("CloseBill: " + billID)
	err := s.client.SignalWorkflow(context.Background(), billID, "", ActionChannels.CLOSE_CHANNEL, nil)
	if err != nil {
		return err
	}
	return nil
}

//encore:api public method=POST path=/bill/:billID/items
func (s *Service) AddItem(ctx context.Context, billID string, item BillItemDTO) error {
	fmt.Println("AddItem: "+billID, " item: ", item)

	payload := AddItemSignal{Item: BillItem{
		ItemID: item.ItemID,
		Name:   item.Name,
		Price:  item.Price,
	}}
	err := s.client.SignalWorkflow(context.Background(), billID, "", ActionChannels.ADD_ITEM_CHANNEL, payload)
	if err != nil {
		return err
	}
	return nil
}

//encore:api public method=DELETE path=/bill/:billID/items/:itemID
func (s *Service) RemoveItem(ctx context.Context, billID string, itemID string) error {
	fmt.Println("RemoveItem: "+billID, " item: ", itemID)

	payload := RemoveItemSignal{ItemID: itemID}
	err := s.client.SignalWorkflow(context.Background(), billID, "", ActionChannels.REMOVE_ITEM_CHANNEL, payload)
	if err != nil {
		return err
	}

	return nil
}
