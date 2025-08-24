package shared

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/constants"
)

// UIInterface defines minimal interface needed by BaseView to avoid circular dependency
type UIInterface interface {
	ShowError(error)
	ShowDetails(any)
	ShowCurrentView()
	ShowConfirm(string, func(bool))
	GetServicesAny() any
	ShowInfo(string)
	ShowNodeAvailabilityModal(string, string, func(string))
	ShowRetryDialog(string, error, func() error, func())
	ShowFallbackDialog(string, error, []string, func(string))
	ShowContextualHelp(string, string)
	ShowServiceScaleModal(string, uint64, func(int))
	GetSwarmServiceService() any
	GetSwarmNodeService() any
	IsContainerServiceAvailable() bool
	GetContainerService() any
}

// ServiceFactoryInterface defines minimal interface for services
type ServiceFactoryInterface interface {
	GetContainerService() any
	GetImageService() any
	GetVolumeService() any
	GetNetworkService() any
	GetDockerInfoService() any
	GetLogsService() any
	GetSwarmServiceService() any
	GetSwarmNodeService() any
	GetCurrentService() any
	SetCurrentService(serviceName string)
	IsServiceAvailable(serviceName string) bool
	IsContainerServiceAvailable() bool
}

// BaseView provides common functionality for all Docker resource views
type BaseView[T any] struct {
	ui       UIInterface
	view     *tview.Flex
	table    *tview.Table
	items    []T
	headers  []string
	viewName string
	log      *slog.Logger

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
func NewBaseView[T any](ui UIInterface, viewName string, headers []string) *BaseView[T] {
	view := builders.NewViewBuilder().CreateView()
	table := builders.NewTableBuilder().CreateTable()
	builders.NewTableBuilder().SetupHeaders(table, headers)

	bv := &BaseView[T]{
		ui:       ui,
		view:     view,
		table:    table,
		headers:  headers,
		viewName: viewName,
		log:      logger.GetLogger(),
	}

	bv.setupKeyBindings()
	view.AddItem(table, 0, 1, true)

	return bv
}

// GetView returns the underlying tview.Primitive view component
func (bv *BaseView[T]) GetView() tview.Primitive {
	return bv.view
}

// GetUI returns the UI interface for this view
func (bv *BaseView[T]) GetUI() UIInterface {
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

// GetTitle returns the view title for testing purposes
func (bv *BaseView[T]) GetTitle() string {
	return bv.viewName
}

// GetHeaders returns the table headers for testing purposes
func (bv *BaseView[T]) GetHeaders() []string {
	return bv.headers
}

// GetSelectedItem returns the currently selected item from the table
func (bv *BaseView[T]) GetSelectedItem() *T {
	row, _ := bv.table.GetSelection()
	if row <= 0 || row > len(bv.items) {
		return nil
	}
	item := bv.items[row-1]
	return &item
}

// Refresh updates the view by fetching and displaying the latest items
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

// ShowItemDetails displays detailed information about a selected item
func (bv *BaseView[T]) ShowItemDetails(item T, inspectData map[string]any, err error) {
	actions := bv.getActionsForItem()

	if bv.shouldShowErrorDetails(err) {
		bv.showErrorDetails(item, err, actions)
		return
	}

	bv.showSuccessDetails(item, inspectData, actions)
}

// ShowConfirmDialog displays a confirmation dialog with the given message and callback
func (bv *BaseView[T]) ShowConfirmDialog(message string, onConfirm func()) {
	bv.ui.ShowConfirm(message, func(confirmed bool) {
		if confirmed {
			onConfirm()
		}
	})
}

func (bv *BaseView[T]) setupKeyBindings() {
	// Set up key bindings on the table
	bv.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		bv.log.Info("Table InputCapture called", "key", string(event.Rune()), "keyType", event.Key())

		row, _ := bv.table.GetSelection()
		if row <= 0 || row > len(bv.items) {
			bv.log.Info("No valid selection", "row", row, "itemCount", len(bv.items))
			return event
		}

		item := bv.items[row-1]

		if event.Key() == tcell.KeyEnter {
			bv.log.Info("Enter key pressed")
			if bv.ShowDetails != nil {
				bv.ShowDetails(item)
			}
			return nil
		}

		// Handle action keys
		if event.Key() == tcell.KeyRune && bv.HandleKeyPress != nil {
			// Get available actions for this view
			actions := bv.getActionsForItem()

			bv.log.Info("Checking action key", "key", string(event.Rune()), "actions", actions, "GetActions_nil", bv.GetActions == nil)

			// If this key is a valid action, handle it and consume the event
			if _, isValidAction := actions[event.Rune()]; isValidAction {
				bv.log.Info("Valid action key, handling", "key", string(event.Rune()))
				bv.HandleKeyPress(event.Rune(), item)
				return nil // Consume the event so it doesn't propagate to global handlers
			}
			bv.log.Info("Key is not a valid action", "key", string(event.Rune()))
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
	bv.log.Info("getActionsForItem called", "GetActions_nil", bv.GetActions == nil)

	actions := bv.GetActions()
	bv.log.Info("GetActions returned", "actions", actions)

	if actions == nil {
		bv.log.Info("Using fallback actions")
		actions = map[rune]string{'d': "Delete", 'i': "Inspect"}
	}

	bv.log.Info("Final actions", "actions", actions)
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
