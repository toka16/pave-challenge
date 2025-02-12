package bill

import (
	"fmt"
	"slices"
)

var allowedCurrencies = []string{"GEL", "USD"}

type CreateBillRequest struct {
	Currency string `json:"currency"`
}

func (r CreateBillRequest) Validate() error {
	if r.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if !slices.Contains(allowedCurrencies, r.Currency) {
		return fmt.Errorf("Invalid currency. Available currencies are: %v", allowedCurrencies)
	}
	return nil
}

type BillResponse struct {
	ID       string        `json:"id"`
	Currency string        `json:"currency"`
	Sum      float64       `json:"sum"`
	Items    []BillItemDTO `json:"items"`
	Status   string        `json:"status"`
}
type BillItemDTO struct {
	ItemID string  `json:"item_id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

func (b BillItemDTO) Validate() error {
	if b.ItemID == "" {
		return fmt.Errorf("item_id is required")
	}
	if b.Name == "" {
		return fmt.Errorf("name is required")
	}
	if b.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	return nil
}

func BillResponseFromState(state BillState) *BillResponse {
	res := &BillResponse{
		ID:       state.ID,
		Currency: state.Currency,
		Sum:      0,
		Items:    make([]BillItemDTO, 0, len(state.Items)),
		Status:   string(state.Status),
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
