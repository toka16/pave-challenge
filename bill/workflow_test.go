package bill

import (
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
	"testing"
)

func Test_BillWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	onAccept := func() {}
	shouldNotAccept := func() { require.Fail(t, "Should not accept") }
	shouldNotReject := func(err error) { require.Fail(t, "Should not reject") }
	shouldNotComplete := func(i interface{}, err error) { require.Fail(t, "Should not complete") }

	// add first item
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: onAccept,
			OnReject: shouldNotReject,
			OnComplete: func(i interface{}, err error) {
				require.NoError(t, err)
				billState, ok := i.(BillState)
				if !ok {
					require.Fail(t, "Invalid return type")
				}
				require.Len(t, billState.Items, 1)
				require.Equal(t, "it1", billState.Items[0].ItemID)
			},
		}, ACTION_ADD_ITEM, ModifyItemData{Item: BillItem{ItemID: "it1", Name: "apple", Price: 1.0}})
	}, 0)

	// add item with same id
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: onAccept,
			OnReject: shouldNotReject,
			OnComplete: func(i interface{}, err error) {
				require.NoError(t, err)
				billState, ok := i.(BillState)
				if !ok {
					require.Fail(t, "Invalid return type")
				}
				require.Len(t, billState.Items, 1)
				require.Equal(t, "it1", billState.Items[0].ItemID)
			},
		}, ACTION_ADD_ITEM, ModifyItemData{Item: BillItem{ItemID: "it1", Name: "apple", Price: 1.0}})
	}, 0)

	// add item with empty id
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: shouldNotAccept,
			OnReject: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "item id cannot be empty")
			},
			OnComplete: shouldNotComplete,
		}, ACTION_ADD_ITEM, ModifyItemData{Item: BillItem{ItemID: "", Name: "apple", Price: 1.0}})
	}, 0)

	// add item with empty name
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: shouldNotAccept,
			OnReject: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "name cannot be empty")
			},
			OnComplete: shouldNotComplete,
		}, ACTION_ADD_ITEM, ModifyItemData{Item: BillItem{ItemID: "it2", Name: "", Price: 1.0}})
	}, 0)

	// add item with price 0
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: shouldNotAccept,
			OnReject: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "price must be greater than 0")
			},
			OnComplete: shouldNotComplete,
		}, ACTION_ADD_ITEM, ModifyItemData{Item: BillItem{ItemID: "it2", Name: "apple", Price: 0.0}})
	}, 0)

	// add item with negative price
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: shouldNotAccept,
			OnReject: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "price must be greater than 0")
			},
			OnComplete: shouldNotComplete,
		}, ACTION_ADD_ITEM, ModifyItemData{Item: BillItem{ItemID: "it2", Name: "apple", Price: -5}})
	}, 0)

	// add second item
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: onAccept,
			OnReject: shouldNotReject,
			OnComplete: func(i interface{}, err error) {
				require.NoError(t, err)
				billState, ok := i.(BillState)
				if !ok {
					require.Fail(t, "Invalid return type")
				}
				require.Len(t, billState.Items, 2)
				require.Equal(t, "it2", billState.Items[1].ItemID)
			},
		}, ACTION_ADD_ITEM, ModifyItemData{Item: BillItem{ItemID: "it2", Name: "banana", Price: 3.0}})
	}, 0)

	// query bill
	env.RegisterDelayedCallback(func() {
		handle, err := env.QueryWorkflow("getBill")
		require.NoError(t, err)
		billState := BillState{Items: make([]BillItem, 0)}
		err = handle.Get(&billState)
		require.NoError(t, err)
		require.Len(t, billState.Items, 2)
	}, 1)

	// remove first item
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: onAccept,
			OnReject: shouldNotReject,
			OnComplete: func(i interface{}, err error) {
				require.NoError(t, err)
				billState, ok := i.(BillState)
				if !ok {
					require.Fail(t, "Invalid return type")
				}
				require.Len(t, billState.Items, 1)
				require.Equal(t, "it2", billState.Items[0].ItemID)
			},
		}, ACTION_REMOVE_ITEM, ModifyItemData{Item: BillItem{ItemID: "it1"}})
	}, 0)

	// remove item with empty id
	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept: shouldNotAccept,
			OnReject: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "item id cannot be empty")
			},
			OnComplete: shouldNotComplete,
		}, ACTION_REMOVE_ITEM, ModifyItemData{Item: BillItem{ItemID: ""}})
	}, 0)

	// close bill
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(CHANNEL_CLOSE, nil)
	}, 0)

	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept:   shouldNotAccept,
			OnReject:   shouldNotReject,
			OnComplete: shouldNotComplete,
		}, ACTION_REMOVE_ITEM, ModifyItemData{Item: BillItem{ItemID: "it2"}})
	}, 10)

	env.RegisterDelayedCallback(func() {
		env.UpdateWorkflow(UPDATE_NAME_MODIFY_ITEMS, "", &testsuite.TestUpdateCallback{
			OnAccept:   shouldNotAccept,
			OnReject:   shouldNotReject,
			OnComplete: shouldNotComplete,
		}, ACTION_ADD_ITEM, ModifyItemData{Item: BillItem{ItemID: "it3"}})
	}, 10)

	env.ExecuteWorkflow(BillWorkflow, BillState{})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

}
