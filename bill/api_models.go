package bill

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
