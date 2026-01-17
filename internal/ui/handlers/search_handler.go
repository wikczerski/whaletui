package handlers

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/logger"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// SearchHandler manages search functionality
type SearchHandler struct {
	*BaseInputHandler
	searchInput    *tview.InputField
	log            *slog.Logger
	lastSearchTerm string
	isExitingInput bool
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(ui interfaces.UIInterface) *SearchHandler {
	return &SearchHandler{
		BaseInputHandler: NewBaseInputHandler(ui),
		log:              logger.GetLogger(),
	}
}

// CreateSearchInput creates and configures the search input field
func (sh *SearchHandler) CreateSearchInput() *tview.InputField {
	sh.searchInput = tview.NewInputField()
	sh.configureSearchInput()
	sh.hideSearchInput()
	sh.SetInput(sh.searchInput)
	return sh.searchInput
}

// Enter activates search mode
func (sh *SearchHandler) Enter() {
	sh.SetActive(true)
	sh.showSearchInput()

	// Check if search is already active in the current view
	if sh.isSearchCurrentlyActive() {
		// Search is already active, just restore the input field
		if sh.lastSearchTerm != "" {
			sh.searchInput.SetText(sh.lastSearchTerm)
		}
	} else {
		// Search is not active, restore the last search term and apply it
		if sh.lastSearchTerm != "" {
			sh.searchInput.SetText(sh.lastSearchTerm)
			sh.ProcessSearch(sh.lastSearchTerm)
		}
	}

	sh.ShowInput()
}

// Exit deactivates search mode and clears search
func (sh *SearchHandler) Exit() {
	sh.SetActive(false)
	sh.hideSearchInput()
	sh.HideInput()
	sh.clearSearch() // Clear search when exiting
}

// ExitInputMode exits input mode but keeps search active
func (sh *SearchHandler) ExitInputMode() {
	sh.isExitingInput = true
	sh.SetActive(false)
	sh.hideSearchInput()
	// Don't hide the input completely, just make it invisible
	// This preserves the search state
	sh.HideInput()
	// Note: Search remains active, we don't call clearSearch()
	sh.isExitingInput = false
}

// HandleInput processes search input
func (sh *SearchHandler) HandleInput(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		// Store the current search term and keep search active
		currentTerm := sh.searchInput.GetText()
		sh.lastSearchTerm = currentTerm
		sh.log.Info("Enter pressed in search",
			"currentTerm", currentTerm,
			"lastSearchTerm", sh.lastSearchTerm)

		// Ensure the search is applied before exiting input mode
		sh.log.Info("About to call ProcessSearch", "currentTerm", currentTerm)
		sh.ProcessSearch(currentTerm)
		sh.log.Info("ProcessSearch completed", "currentTerm", currentTerm)

		sh.log.Info("About to call ExitInputMode")
		sh.ExitInputMode() // Exit input mode but keep search active
		sh.log.Info("ExitInputMode completed")
	case tcell.KeyEscape:
		// Clear search and exit search mode
		sh.searchInput.SetText("")
		sh.lastSearchTerm = "" // Clear the stored search term
		sh.clearSearch()
		sh.Exit()
	case tcell.KeyRune:
		// User is typing - clear any error message
		sh.ClearError()
		// Real-time search is handled by SetChangedFunc, no need to process here
	}
}

// ProcessSearch handles search functionality (public method)
func (sh *SearchHandler) ProcessSearch(searchTerm string) {
	sh.lastSearchTerm = searchTerm // Store the search term
	sh.processSearch(searchTerm)
}

// processSearch handles search functionality (private method)
func (sh *SearchHandler) processSearch(searchTerm string) {
	sh.log.Info("processSearch called", "searchTerm", searchTerm)
	// Get the current view and try to call its search method
	currentView := sh.getCurrentView()
	if currentView != nil {
		sh.log.Info("Current view found", "viewType", fmt.Sprintf("%T", currentView))
		// Use reflection to call the Search method
		viewValue := reflect.ValueOf(currentView)
		searchMethod := viewValue.MethodByName("Search")

		if searchMethod.IsValid() && !searchMethod.IsNil() {
			sh.log.Info("Calling Search method via reflection", "searchTerm", searchTerm)
			searchMethod.Call([]reflect.Value{reflect.ValueOf(searchTerm)})
			sh.log.Info("Search method called successfully")
		} else {
			sh.log.Warn("Search method not found or invalid on current view",
				"viewType", fmt.Sprintf("%T", currentView),
			)
		}
	} else {
		sh.log.Warn("No current view found for search")
	}
}

// isSearchCurrentlyActive checks if search is currently active in the view
func (sh *SearchHandler) isSearchCurrentlyActive() bool {
	currentView := sh.getCurrentView()
	if currentView == nil {
		return false
	}

	// Try to call IsSearchActive method using reflection
	viewValue := reflect.ValueOf(currentView)
	isActiveMethod := viewValue.MethodByName("IsSearchActive")

	if !isActiveMethod.IsValid() || isActiveMethod.IsNil() {
		return false
	}

	results := isActiveMethod.Call([]reflect.Value{})
	if len(results) == 0 || !results[0].CanInterface() {
		return false
	}

	isActive, ok := results[0].Interface().(bool)
	return ok && isActive
}

// clearSearch clears the current search
func (sh *SearchHandler) clearSearch() {
	currentView := sh.getCurrentView()
	if currentView != nil {
		// Use reflection to call the ClearSearch method
		viewValue := reflect.ValueOf(currentView)
		clearMethod := viewValue.MethodByName("ClearSearch")

		if clearMethod.IsValid() && !clearMethod.IsNil() {
			clearMethod.Call([]reflect.Value{})
		}
	}
}

// getCurrentView gets the currently active view
func (sh *SearchHandler) getCurrentView() any {
	// Try to get the actual view from the UI
	if uiWithViews, ok := sh.ui.(interface{ GetCurrentView() any }); ok {
		currentView := uiWithViews.GetCurrentView()
		sh.log.Info("getCurrentView called",
			"currentView", currentView,
			"viewType", fmt.Sprintf("%T", currentView),
		)
		return currentView
	}

	sh.log.Warn("UI does not implement GetCurrentView interface")
	return nil
}

// configureSearchInput sets up the search input styling and behavior
func (sh *SearchHandler) configureSearchInput() {
	themeManager := sh.ui.GetThemeManager()
	sh.setupSearchStyling(themeManager)
	sh.setupSearchBehavior()
}

// setupSearchStyling sets up the styling for the search input
func (sh *SearchHandler) setupSearchStyling(themeManager *config.ThemeManager) {
	sh.searchInput.SetLabel("/ ")
	sh.searchInput.SetTitle(" Search ")
	sh.searchInput.SetPlaceholder("Type to search... (ESC to clear)")
	sh.searchInput.SetFieldTextColor(themeManager.GetCommandModeTextColor())
	sh.searchInput.SetLabelColor(themeManager.GetCommandModeLabelColor())
	sh.searchInput.SetTitleColor(themeManager.GetCommandModeTitleColor())
	sh.searchInput.SetBackgroundColor(themeManager.GetCommandModeBackgroundColor())
	sh.searchInput.SetPlaceholderTextColor(themeManager.GetCommandModePlaceholderColor())
	sh.searchInput.SetBorder(true)
	sh.searchInput.SetBorderColor(themeManager.GetCommandModeBorderColor())
}

// setupSearchBehavior sets up the behavior for the search input
func (sh *SearchHandler) setupSearchBehavior() {
	// Set up real-time search using SetChangedFunc
	sh.searchInput.SetChangedFunc(func(text string) {
		sh.log.Info("SetChangedFunc called", "text", text, "isExitingInput", sh.isExitingInput)
		// Don't process empty text if we're exiting input mode to avoid clearing search
		if !sh.isExitingInput || text != "" {
			sh.processSearch(text)
		} else {
			sh.log.Info("Skipping processSearch for empty text during exit")
		}
	})

	// Handle special keys and trigger search after input processing
	sh.searchInput.SetDoneFunc(func(key tcell.Key) {
		// Only process search for non-Enter keys to avoid double processing
		if key != tcell.KeyEnter {
			text := sh.searchInput.GetText()
			sh.log.Info("SetDoneFunc called", "key", key, "text", text)
			sh.processSearch(text)
		} else {
			sh.log.Info("SetDoneFunc called with Enter - skipping processSearch to avoid double processing")
		}
		sh.HandleInput(key)
	})
}

// hideSearchInput makes the search input completely invisible
func (sh *SearchHandler) hideSearchInput() {
	if sh.searchInput == nil {
		return
	}
	sh.searchInput.SetText("")
	sh.searchInput.SetPlaceholder("")
	sh.searchInput.SetLabel("")
	sh.searchInput.SetTitle("")
}

// showSearchInput makes the search input visible
func (sh *SearchHandler) showSearchInput() {
	if sh.searchInput == nil {
		return
	}
	sh.searchInput.SetPlaceholder("Type to search... (ESC to clear)")
	sh.searchInput.SetLabel("/ ")
	sh.searchInput.SetTitle(" Search ")
}
