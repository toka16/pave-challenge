package bill

type BillStatus string

const (
	BillStatusOpen   BillStatus = "open"
	BillStatusClosed BillStatus = "closed"
)

type BillState struct {
	ID       string
	Currency string
	Items    []BillItem
	Status   BillStatus
}

type BillItem struct {
	ItemID string
	Name   string
	Price  float64
}

func (b *BillState) AddItem(item BillItem) {
	for i, v := range b.Items {
		if v.ItemID == item.ItemID {
			b.Items[i] = item // override the existing item
			return
		}
	}
	b.Items = append(b.Items, item)
}

func (b *BillState) RemoveItem(itemID string) {
	for i, v := range b.Items {
		if v.ItemID == itemID {
			b.Items = append(b.Items[:i], b.Items[i+1:]...)
			return
		}
	}
}
