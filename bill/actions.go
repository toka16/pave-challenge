package bill

var ActionChannels = struct {
	ADD_ITEM_CHANNEL    string
	REMOVE_ITEM_CHANNEL string
	CLOSE_CHANNEL       string
}{
	ADD_ITEM_CHANNEL:    "ADD_ITEM_CHANNEL",
	REMOVE_ITEM_CHANNEL: "REMOVE_ITEM_CHANNEL",
	CLOSE_CHANNEL:       "CLOSE_CHANNEL",
}

type AddItemSignal struct {
	Item BillItem
}

type RemoveItemSignal struct {
	ItemID string
}
