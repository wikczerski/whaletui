package views

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// BaseView provides common functionality for all Docker resource views
type BaseView[T any] struct {
	ui       interfaces.UIInterface
	view     *tview.Flex
	table    *tview.Table
	items    []T
	headers  []string
	viewName string

	// Callbacks for specific behavior
	ListItems      func(ctx context.Context) ([]T, error)
	FormatRow      func(item T) []string
	GetRowColor    func(item T) tcell.Color // Optional: custom row colors
	GetItemID      func(item T) string
	GetItemName    func(item T) string
	HandleKeyPress func(key rune, item T)
	ShowDetails    func(item T)
	GetActions     func() map[rune]string
}

// NewBaseView creates a new base view with common functionality
func NewBaseView[T any](ui interfaces.UIInterface, viewName string, headers []string) *BaseView[T] {
	view := builders.NewViewBuilder().CreateView()
	table := builders.NewTableBuilder().CreateTable()
	builders.NewTableBuilder().SetupHeaders(table, headers)

	bv := &BaseView[T]{
		ui:       ui,
		view:     view,
		table:    table,
		headers:  headers,
		viewName: viewName,
	}

	bv.setupKeyBindings()
	view.AddItem(table, 0, 1, true)

	return bv
}

func (bv *BaseView[T]) setupKeyBindings() {
	bv.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		row, _ := bv.table.GetSelection()
		if row <= 0 || row > len(bv.items) {
			return event
		}

		item := bv.items[row-1]

		if event.Key() == tcell.KeyEnter {
			if bv.ShowDetails != nil {
				bv.ShowDetails(item)
			}
			return nil
		}

		if bv.HandleKeyPress != nil {
			bv.HandleKeyPress(event.Rune(), item)
		}

		return event
	})
}

func (bv *BaseView[T]) GetView() tview.Primitive {
	return bv.view
}

func (bv *BaseView[T]) Refresh() {
	if bv.ListItems == nil {
		return
	}

	items, err := bv.ListItems(context.Background())
	if err != nil {
		bv.ui.ShowError(err)
		return
	}

	bv.updateItemsAndTable(items)
}

// updateItemsAndTable updates the items and refreshes the table display
func (bv *BaseView[T]) updateItemsAndTable(items []T) {
	bv.items = items
	bv.table.Clear()
	builders.NewTableBuilder().SetupHeaders(bv.table, bv.headers)

	if bv.FormatRow != nil {
		bv.populateTableRows(items)
	}

	bv.selectFirstRowIfAvailable(items)
}

// populateTableRows populates the table with formatted row data
func (bv *BaseView[T]) populateTableRows(items []T) {
	for i, item := range items {
		cells := bv.FormatRow(item)
		rowColor := bv.getRowColor(item)
		builders.NewTableBuilder().SetupRow(bv.table, i+1, cells, rowColor)
	}
}

// getRowColor gets the color for a table row, with fallback to default
func (bv *BaseView[T]) getRowColor(item T) tcell.Color {
	if bv.GetRowColor != nil {
		return bv.GetRowColor(item)
	}
	return constants.TableDefaultRowColor
}

// selectFirstRowIfAvailable selects the first row if items are available
func (bv *BaseView[T]) selectFirstRowIfAvailable(items []T) {
	if len(items) > 0 {
		bv.table.Select(1, 0)
	}
}

func (bv *BaseView[T]) ShowItemDetails(item T, inspectData map[string]any, err error) {
	actions := bv.getActionsForItem()

	if bv.shouldShowErrorDetails(err) {
		bv.showErrorDetails(item, err, actions)
		return
	}

	bv.showSuccessDetails(item, inspectData, actions)
}

// shouldShowErrorDetails determines if error details should be shown
func (bv *BaseView[T]) shouldShowErrorDetails(err error) bool {
	return err != nil
}

// getActionsForItem gets the actions for the item, with fallback to defaults
func (bv *BaseView[T]) getActionsForItem() map[rune]string {
	actions := bv.GetActions()
	if actions == nil {
		actions = map[rune]string{'d': "Delete", 'i': "Inspect"}
	}
	return actions
}

// showErrorDetails shows error details when inspection fails
func (bv *BaseView[T]) showErrorDetails(item T, err error, actions map[rune]string) {
	itemName := bv.GetItemName(item)
	details := fmt.Sprintf("Item: %s\nInspect error: %v", itemName, err)
	detailsView := builders.CreateDetailsView(itemName, details, actions, bv.handleAction, bv.showTable)
	bv.ui.ShowDetails(detailsView)
}

// showSuccessDetails shows successful inspection details
func (bv *BaseView[T]) showSuccessDetails(item T, inspectData map[string]any, actions map[rune]string) {
	itemName := bv.GetItemName(item)
	detailsView := builders.CreateInspectDetailsView(itemName, inspectData, actions, bv.handleAction, bv.showTable)
	bv.ui.ShowDetails(detailsView)
}

func (bv *BaseView[T]) handleAction(key rune) {
	row, _ := bv.table.GetSelection()
	if row <= 0 || row > len(bv.items) {
		return
	}
	item := bv.items[row-1]

	if bv.HandleKeyPress != nil {
		bv.HandleKeyPress(key, item)
	}
	bv.ui.ShowCurrentView()
}

func (bv *BaseView[T]) showTable() {
	bv.ui.ShowCurrentView()
}

func (bv *BaseView[T]) ShowConfirmDialog(message string, onConfirm func()) {
	bv.ui.ShowConfirm(message, func(confirmed bool) {
		if confirmed {
			onConfirm()
		}
	})
}
