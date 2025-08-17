package views

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/ui/builders"
	"github.com/wikczerski/D5r/internal/ui/constants"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
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

	bv.items = items
	bv.table.Clear()
	builders.NewTableBuilder().SetupHeaders(bv.table, bv.headers)

	if bv.FormatRow != nil {
		for i, item := range items {
			cells := bv.FormatRow(item)
			rowColor := constants.TableDefaultRowColor
			if bv.GetRowColor != nil {
				rowColor = bv.GetRowColor(item)
			}
			builders.NewTableBuilder().SetupRow(bv.table, i+1, cells, rowColor)
		}
	}

	if len(items) > 0 {
		bv.table.Select(1, 0)
	}
}

func (bv *BaseView[T]) ShowItemDetails(item T, inspectData map[string]any, err error) {
	actions := bv.GetActions()
	if actions == nil {
		actions = map[rune]string{'d': "Delete", 'i': "Inspect"}
	}

	if err != nil {
		itemName := bv.GetItemName(item)
		details := fmt.Sprintf("Item: %s\nInspect error: %v", itemName, err)
		detailsView := builders.CreateDetailsView(itemName, details, actions, bv.handleAction, bv.showTable)
		bv.ui.ShowDetails(detailsView)
		return
	}

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
