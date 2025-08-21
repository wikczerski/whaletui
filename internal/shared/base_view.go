package shared

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

func (bv *BaseView[T]) GetView() tview.Primitive {
	return bv.view
}

func (bv *BaseView[T]) GetUI() interfaces.UIInterface {
	return bv.ui
}

// GetTable returns the table for testing purposes
func (bv *BaseView[T]) GetTable() *tview.Table {
	return bv.table
}

// GetItems returns the items for testing purposes
func (bv *BaseView[T]) GetItems() []T {
	return bv.items
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

func (bv *BaseView[T]) ShowItemDetails(item T, inspectData map[string]any, err error) {
	actions := bv.getActionsForItem()

	if bv.shouldShowErrorDetails(err) {
		bv.showErrorDetails(item, err, actions)
		return
	}

	bv.showSuccessDetails(item, inspectData, actions)
}

func (bv *BaseView[T]) ShowConfirmDialog(message string, onConfirm func()) {
	bv.ui.ShowConfirm(message, func(confirmed bool) {
		if confirmed {
			onConfirm()
		}
	})
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

// updateItemsAndTable updates the items and refreshes the table display
func (bv *BaseView[T]) updateItemsAndTable(items []T) {
	// Store items first to ensure they're available for selection
	bv.items = items

	// Completely clear the table to prevent leftover content
	bv.table.Clear()

	// Reset table selection to prevent display issues
	bv.table.Select(0, 0)

	// Ensure the table is properly initialized
	bv.table.SetFixed(1, 0)

	// Ensure headers are always visible
	builders.NewTableBuilder().SetupHeaders(bv.table, bv.headers)

	// Only populate rows if we have items and a formatter
	if bv.FormatRow != nil && len(items) > 0 {
		bv.populateTableRows(items)
	}

	// Ensure proper selection state
	bv.selectFirstRowIfAvailable(items)

	// Note: Don't call table.Draw(nil) manually - tview handles drawing automatically
	// This prevents the panic and ensures proper UI updates
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
		// Ensure we're selecting a valid row and column
		bv.table.Select(1, 0)
	} else {
		// If no items, clear selection to prevent display issues
		bv.table.Select(0, 0)
	}
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
