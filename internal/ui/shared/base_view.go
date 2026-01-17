package shared

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/ui/builders"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
	"github.com/wikczerski/whaletui/internal/ui/utils"
)

// BaseView provides common functionality for all Docker resource views
type BaseView[T any] struct {
	ui       interfaces.SharedUIInterface
	view     *tview.Flex
	table    *tview.Table
	items    []T
	headers  []string
	viewName string
	log      *slog.Logger

	// Search functionality
	searchTerm     string
	filteredItems  []T
	originalItems  []T // Store original items separately from current display
	isSearchActive bool

	// Callbacks for specific behavior
	ListItems           func(ctx context.Context) ([]T, error)
	FormatRow           func(item T) []string
	GetRowColor         func(item T) tcell.Color // Optional: custom row colors
	GetItemID           func(item T) string
	GetItemName         func(item T) string
	HandleKeyPress      func(key rune, item T)
	ShowDetailsCallback func(item T)
	GetActions          func() map[rune]string

	// Character limits support
	columnTypes []string
	formatter   *utils.TableFormatter
}

// NewBaseView creates a new base view with common functionality
func NewBaseView[T any](ui interfaces.SharedUIInterface, viewName string, headers []string) *BaseView[T] {
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
func (bv *BaseView[T]) GetUI() interfaces.SharedUIInterface {
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

// SetItems sets the items for testing purposes
func (bv *BaseView[T]) SetItems(items []T) {
	bv.items = items
	bv.originalItems = make([]T, len(items))
	copy(bv.originalItems, items)
}

// GetTitle returns the view title for testing purposes
func (bv *BaseView[T]) GetTitle() string {
	return bv.viewName
}

// GetHeaders returns the table headers for testing purposes
func (bv *BaseView[T]) GetHeaders() []string {
	return bv.headers
}

// SetColumnTypes sets the column types for character limit formatting
func (bv *BaseView[T]) SetColumnTypes(columnTypes []string) {
	bv.columnTypes = columnTypes
}

// SetFormatter sets the table formatter for character limits
func (bv *BaseView[T]) SetFormatter(formatter *utils.TableFormatter) {
	bv.formatter = formatter
}

// RefreshFormatter updates the formatter with the latest theme configuration
func (bv *BaseView[T]) RefreshFormatter(ui interfaces.SharedUIInterface) {
	if ui != nil {
		themeManager := ui.GetThemeManager()
		if themeManager != nil {
			// Create a new formatter with the updated limits
			bv.formatter = utils.NewTableFormatterFromTheme(themeManager)
		}
	}
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
	bv.log.Info("BaseView.Refresh called",
		"isSearchActive", bv.isSearchActive,
		"searchTerm", bv.searchTerm)

	if bv.ListItems == nil {
		return
	}

	// Refresh the formatter to get latest theme configuration
	bv.RefreshFormatter(bv.ui)

	items, err := bv.ListItems(context.Background())
	if err != nil {
		bv.ui.ShowError(err)
		return
	}

	bv.log.Info("Refresh got items from ListItems", "count", len(items))

	// Store the original items
	bv.originalItems = make([]T, len(items))
	copy(bv.originalItems, items)

	// If we're in search mode, apply local filtering
	if bv.isSearchActive && bv.searchTerm != "" {
		bv.log.Info("Refresh applying local filter")
		filteredItems := bv.filterItems(bv.searchTerm)
		bv.updateItemsAndTable(filteredItems)
	} else {
		bv.log.Info("Refresh showing all items")
		bv.updateItemsAndTable(items)
	}
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
	if ui, ok := bv.ui.(interface{ ShowConfirm(string, func(bool)) }); ok {
		ui.ShowConfirm(message, func(confirmed bool) {
			if confirmed {
				onConfirm()
			}
		})
	}
}

func (bv *BaseView[T]) setupKeyBindings() {
	// Set up key bindings on the table
	bv.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		bv.log.Info(
			"Table InputCapture called",
			"key",
			string(event.Rune()),
			"keyType",
			event.Key(),
		)

		row, _ := bv.table.GetSelection()
		if row <= 0 || row > len(bv.items) {
			bv.log.Info("No valid selection", "row", row, "itemCount", len(bv.items))
			return event
		}

		item := bv.items[row-1]

		if event.Key() == tcell.KeyEnter {
			bv.log.Info("Enter key pressed")
			if bv.ShowDetailsCallback != nil {
				bv.ShowDetailsCallback(item)
			}
			return nil
		}

		// Handle action keys
		if event.Key() == tcell.KeyRune && bv.HandleKeyPress != nil {
			return bv.handleActionKey(event, item)
		}

		return event
	})
}

// handleActionKey handles action key presses
func (bv *BaseView[T]) handleActionKey(event *tcell.EventKey, item T) *tcell.EventKey {
	// Get available actions for this view
	actions := bv.getActionsForItem()

	bv.log.Info("Checking action key", "key", string(event.Rune()),
		"actions", actions, "GetActions_nil", bv.GetActions == nil)

	// If this key is a valid action, handle it and consume the event
	if _, isValidAction := actions[event.Rune()]; isValidAction {
		bv.log.Info("Valid action key, handling", "key", string(event.Rune()))
		bv.HandleKeyPress(event.Rune(), item)
		return nil // Consume the event so it doesn't propagate to global handlers
	}
	bv.log.Info("Key is not a valid action", "key", string(event.Rune()))
	return event
}

// updateItemsAndTable updates the items and refreshes the table display
func (bv *BaseView[T]) updateItemsAndTable(items []T) {
	// Capture current selection before updating items
	selectedID := bv.saveSelection()

	// Store items first to ensure they're available for selection
	bv.items = items

	// Always store as original items when not in search mode
	if !bv.isSearchActive {
		bv.originalItems = make([]T, len(items))
		copy(bv.originalItems, items)
	}

	// Completely clear the table to prevent leftover content
	bv.table.Clear()

	// Reset table selection to prevent display issues
	bv.table.Select(0, 0)

	// Ensure the table is properly initialized
	bv.table.SetFixed(1, 0)

	// Ensure headers are always visible
	bv.setupTableHeaders()

	// Only populate rows if we have items and a formatter
	if bv.FormatRow != nil && len(items) > 0 {
		bv.populateTableRows(items)
	}

	// Restore selection
	bv.restoreSelection(selectedID, items)

	// Note: Don't call table.Draw(nil) manually - tview handles drawing automatically
	// This prevents the panic and ensures proper UI updates
}

// saveSelection saves the current selection ID
func (bv *BaseView[T]) saveSelection() string {
	if bv.GetItemID != nil {
		row, _ := bv.table.GetSelection()
		// Check against OLD items
		if row > 0 && row <= len(bv.items) {
			return bv.GetItemID(bv.items[row-1])
		}
	}
	return ""
}

// restoreSelection restores the selection based on saved ID
func (bv *BaseView[T]) restoreSelection(selectedID string, items []T) {
	restored := false
	if selectedID != "" {
		for i, item := range items {
			if bv.GetItemID != nil && bv.GetItemID(item) == selectedID {
				// i+1 because of header
				bv.table.Select(i+1, 0)
				restored = true
				break
			}
		}
	}

	// Ensure proper selection state if not restored
	if !restored {
		bv.selectFirstRowIfAvailable(items)
	}
}

// setupTableHeaders sets up the table headers
func (bv *BaseView[T]) setupTableHeaders() {
	if bv.formatter != nil && len(bv.columnTypes) > 0 {
		builders.NewTableBuilder().SetupHeadersWithConfigForView(
			bv.table, bv.headers, bv.columnTypes, bv.formatter, bv.viewName)
	} else {
		builders.NewTableBuilder().SetupHeaders(bv.table, bv.headers)
	}
}

// populateTableRows populates the table with formatted row data
func (bv *BaseView[T]) populateTableRows(items []T) {
	for i, item := range items {
		cells := bv.FormatRow(item)
		rowColor := bv.getRowColor(item)

		// Use character limits if formatter and column types are available
		if bv.formatter != nil && len(bv.columnTypes) > 0 {
			builders.NewTableBuilder().SetupRowWithLimitsForView(
				bv.table, i+1, cells, bv.columnTypes, rowColor, bv.formatter, bv.viewName)
		} else {
			builders.NewTableBuilder().SetupRow(bv.table, i+1, cells, rowColor)
		}
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
	detailsView := builders.CreateDetailsView(
		itemName,
		details,
		actions,
		bv.handleAction,
		bv.showTable,
	)
	bv.ui.ShowDetails(detailsView)
}

// showSuccessDetails shows successful inspection details
func (bv *BaseView[T]) showSuccessDetails(
	item T,
	inspectData map[string]any,
	actions map[rune]string,
) {
	itemName := bv.GetItemName(item)
	detailsView := builders.CreateInspectDetailsView(
		itemName,
		inspectData,
		actions,
		bv.handleAction,
		bv.showTable,
	)
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

// ShowDetails shows details for the selected item
func (bv *BaseView[T]) ShowDetails(item T) {
	if bv.ShowDetailsCallback != nil {
		bv.ShowDetailsCallback(item)
	} else {
		// Default implementation
		detailsView := bv.createDetailsView(item)
		bv.ui.ShowDetails(detailsView)
	}
}

// createDetailsView creates a default details view for the item
func (bv *BaseView[T]) createDetailsView(item T) tview.Primitive {
	details := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	details.SetTitle(fmt.Sprintf("Details: %s", bv.GetItemName(item)))
	details.SetBorder(true)

	// Format the item details
	details.SetText(fmt.Sprintf("ID: %s\nName: %s", bv.GetItemID(item), bv.GetItemName(item)))

	return details
}

// Search filters the items based on the search term
func (bv *BaseView[T]) Search(searchTerm string) {
	bv.log.Info("BaseView.Search called", "term", searchTerm, "viewName", bv.viewName)
	bv.searchTerm = searchTerm
	bv.isSearchActive = searchTerm != ""
	bv.log.Info("Search state set",
		"isSearchActive", bv.isSearchActive,
		"originalItemsCount", len(bv.originalItems),
	)

	if bv.isSearchActive {
		bv.log.Info("Using local filtering")
		// Use local filtering
		bv.filteredItems = bv.filterItems(searchTerm)
		bv.log.Info("Filtered items", "filteredCount", len(bv.filteredItems))
		bv.updateItemsAndTable(bv.filteredItems)
	} else {
		bv.log.Info("Clearing search")
		// Show all original items
		bv.updateItemsAndTable(bv.originalItems)
	}
}

// ClearSearch clears the search and shows all items
func (bv *BaseView[T]) ClearSearch() {
	bv.searchTerm = ""
	bv.isSearchActive = false
	bv.filteredItems = nil
	bv.updateItemsAndTable(bv.originalItems)
}

// IsSearchActive returns whether search is currently active
func (bv *BaseView[T]) IsSearchActive() bool {
	return bv.isSearchActive
}

// GetSearchTerm returns the current search term
func (bv *BaseView[T]) GetSearchTerm() string {
	return bv.searchTerm
}

// filterItems filters items based on the search term
func (bv *BaseView[T]) filterItems(searchTerm string) []T {
	if searchTerm == "" {
		return bv.originalItems
	}

	searchLower := strings.ToLower(searchTerm)
	var filtered []T

	bv.log.Info("Filtering items", "searchTerm", searchTerm, "originalCount", len(bv.originalItems))

	for i, item := range bv.originalItems {
		matches := bv.matchesSearch(item, searchLower)
		if matches {
			filtered = append(filtered, item)
			bv.log.Info("Item matches", "index", i, "item", fmt.Sprintf("%+v", item))
		}
	}

	bv.log.Info("Filtering complete", "filteredCount", len(filtered))
	return filtered
}

// matchesSearch checks if an item matches the search term
func (bv *BaseView[T]) matchesSearch(item T, searchLower string) bool {
	// Search in all visible columns by formatting the row and checking each cell
	if bv.FormatRow != nil {
		cells := bv.FormatRow(item)
		for _, cell := range cells {
			if strings.Contains(strings.ToLower(cell), searchLower) {
				return true
			}
		}
	}

	// Also search in item ID and name if available
	if bv.GetItemID != nil {
		if strings.Contains(strings.ToLower(bv.GetItemID(item)), searchLower) {
			return true
		}
	}

	if bv.GetItemName != nil {
		if strings.Contains(strings.ToLower(bv.GetItemName(item)), searchLower) {
			return true
		}
	}

	return false
}
