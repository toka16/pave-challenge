package bill

const (
	CHANNEL_CLOSE = "CHANNEL_CLOSE"

	UPDATE_NAME_MODIFY_ITEMS = "MODIFY_ITEMS"

	ACTION_ADD_ITEM    = "ADD_ITEM"
	ACTION_REMOVE_ITEM = "REMOVE_ITEM"
)

type ModifyItemData struct {
	Item BillItem
}
