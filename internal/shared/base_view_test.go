package shared

import (
	"context"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mockui "github.com/wikczerski/whaletui/internal/mocks/ui"
	"github.com/wikczerski/whaletui/internal/ui/constants"
)

// Test data structure
type TestItem struct {
	ID   string
	Name string
}

func TestNewBaseView(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	headers := []string{"ID", "Name"}

	bv := NewBaseView[TestItem](mockUI, "TestView", headers)

	assert.NotNil(t, bv)
	assert.Equal(t, mockUI, bv.GetUI())
	assert.Equal(t, "TestView", bv.viewName)
	assert.Equal(t, headers, bv.headers)
	assert.NotNil(t, bv.GetView())
	assert.NotNil(t, bv.GetTable())
}

func TestBaseView_GetView(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	view := bv.GetView()
	assert.NotNil(t, view)
	assert.IsType(t, &tview.Flex{}, view)
}

func TestBaseView_GetUI(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	ui := bv.GetUI()
	assert.Equal(t, mockUI, ui)
}

func TestBaseView_GetTable(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	table := bv.GetTable()
	assert.NotNil(t, table)
	assert.IsType(t, &tview.Table{}, table)
}

func TestBaseView_GetItems(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	items := bv.GetItems()
	assert.Empty(t, items)
}

func TestBaseView_Refresh_NoListItems(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	// Should not panic when ListItems is nil
	assert.NotPanics(t, func() {
		bv.Refresh()
	})
}

func TestBaseView_Refresh_WithListItems(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	testItems := []TestItem{
		{ID: "1", Name: "Item1"},
		{ID: "2", Name: "Item2"},
	}

	bv.ListItems = func(ctx context.Context) ([]TestItem, error) {
		return testItems, nil
	}

	bv.Refresh()

	items := bv.GetItems()
	assert.Equal(t, testItems, items)
}

func TestBaseView_Refresh_WithError(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	mockUI.EXPECT().ShowError(assert.AnError).Once()

	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	bv.ListItems = func(ctx context.Context) ([]TestItem, error) {
		return nil, assert.AnError
	}

	bv.Refresh()
}

func TestBaseView_ShowConfirmDialog(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	mockUI.EXPECT().ShowConfirm("Test message", mock.AnythingOfType("func(bool)")).Once()

	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	bv.ShowConfirmDialog("Test message", func() {})
}

func TestBaseView_ShowItemDetails_WithError(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	mockUI.EXPECT().ShowDetails(mock.AnythingOfType("*tview.Flex")).Once()

	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	bv.GetItemName = func(item TestItem) string {
		return "TestItem"
	}

	bv.GetActions = func() map[rune]string {
		return map[rune]string{'d': "Delete"}
	}

	item := TestItem{ID: "1", Name: "TestItem"}
	bv.ShowItemDetails(item, nil, assert.AnError)
}

func TestBaseView_ShowItemDetails_WithSuccess(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	mockUI.EXPECT().ShowDetails(mock.AnythingOfType("*tview.Flex")).Once()

	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	bv.GetItemName = func(item TestItem) string {
		return "TestItem"
	}

	bv.GetActions = func() map[rune]string {
		return map[rune]string{'d': "Delete"}
	}

	item := TestItem{ID: "1", Name: "TestItem"}
	inspectData := map[string]any{"status": "running"}
	bv.ShowItemDetails(item, inspectData, nil)
}

func TestBaseView_ShowItemDetails_DefaultActions(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	mockUI.EXPECT().ShowDetails(mock.AnythingOfType("*tview.Flex")).Once()

	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	bv.GetItemName = func(item TestItem) string {
		return "TestItem"
	}

	// Ensure GetActions returns nil to test default actions fallback
	bv.GetActions = func() map[rune]string {
		return nil
	}

	item := TestItem{ID: "1", Name: "TestItem"}
	bv.ShowItemDetails(item, nil, assert.AnError)
}

func TestBaseView_UpdateItemsAndTable(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	testItems := []TestItem{
		{ID: "1", Name: "Item1"},
		{ID: "2", Name: "Item2"},
	}

	bv.FormatRow = func(item TestItem) []string {
		return []string{item.ID, item.Name}
	}

	bv.updateItemsAndTable(testItems)

	items := bv.GetItems()
	assert.Equal(t, testItems, items)
}

func TestBaseView_UpdateItemsAndTable_NoFormatter(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	testItems := []TestItem{
		{ID: "1", Name: "Item1"},
		{ID: "2", Name: "Item2"},
	}

	// No FormatRow function set
	bv.updateItemsAndTable(testItems)

	items := bv.GetItems()
	assert.Equal(t, testItems, items)
}

func TestBaseView_UpdateItemsAndTable_EmptyItems(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	bv.updateItemsAndTable([]TestItem{})

	items := bv.GetItems()
	assert.Empty(t, items)
}

func TestBaseView_GetRowColor_Default(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	item := TestItem{ID: "1", Name: "TestItem"}
	color := bv.getRowColor(item)

	// Should return default color when GetRowColor is not set
	assert.Equal(t, constants.TableDefaultRowColor, color)
}

func TestBaseView_GetRowColor_Custom(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	customColor := tcell.ColorRed
	bv.GetRowColor = func(item TestItem) tcell.Color {
		return customColor
	}

	item := TestItem{ID: "1", Name: "TestItem"}
	color := bv.getRowColor(item)

	assert.Equal(t, customColor, color)
}

func TestBaseView_HandleAction(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	mockUI.EXPECT().ShowCurrentView().Once()

	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	// Set up items so we can select one
	testItems := []TestItem{{ID: "1", Name: "Item1"}}
	bv.items = testItems

	// Set up table selection
	bv.table.Select(1, 0)

	// Set up key handler
	keyHandled := false
	bv.HandleKeyPress = func(key rune, item TestItem) {
		keyHandled = true
		assert.Equal(t, 'd', key)
		assert.Equal(t, testItems[0], item)
	}

	bv.handleAction('d')

	assert.True(t, keyHandled)
}

func TestBaseView_HandleAction_NoSelection(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	// No items, so no selection possible
	bv.handleAction('d')

	// Should not panic and should not call HandleKeyPress
}

func TestBaseView_ShowTable(t *testing.T) {
	mockUI := mockui.NewMockUIInterface(t)
	mockUI.EXPECT().ShowCurrentView().Once()

	bv := NewBaseView[TestItem](mockUI, "TestView", []string{"ID", "Name"})

	bv.showTable()
}
