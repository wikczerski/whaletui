package managers

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// ModalManager handles various modal dialogs
type ModalManager struct {
	ui interfaces.UIInterface
}

// NewModalManager creates a new modal manager
func NewModalManager(ui interfaces.UIInterface) *ModalManager {
	return &ModalManager{ui: ui}
}

// ShowHelp displays the help modal with keyboard shortcuts
func (mm *ModalManager) ShowHelp() {
	helpText := mm.buildHelpText()
	modal := mm.createModal(helpText, []string{"Close"})

	// Add done function to handle Close button click
	modal.SetDoneFunc(func(_ int, _ string) {
		pages := mm.ui.GetPages().(*tview.Pages)
		pages.RemovePage("help_modal")
		// Restore focus to the main view after closing modal
		if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
			if vc, ok := viewContainer.(*tview.Flex); ok {
				app := mm.ui.GetApp().(*tview.Application)
				app.SetFocus(vc)
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("help_modal")
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("help_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowError displays an error modal
func (mm *ModalManager) ShowError(err error) {
	errorText := fmt.Sprintf("Error: %v", err)
	modal := mm.createModal(errorText, []string{"OK"})

	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		pages := mm.ui.GetPages().(*tview.Pages)
		pages.RemovePage("error_modal")
		// Restore focus to the main view after closing modal
		if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
			if vc, ok := viewContainer.(*tview.Flex); ok {
				app := mm.ui.GetApp().(*tview.Application)
				app.SetFocus(vc)
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("error_modal")
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("error_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowConfirm displays a confirmation modal with Yes/No buttons
func (mm *ModalManager) ShowConfirm(text string, callback func(bool)) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, _ string) {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("confirm_modal")
			callback(buttonIndex == 0)
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		})

	// Add keyboard handling for ESC key to close modal (cancel action)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("confirm_modal")
			// Call callback with false (No) when ESC is pressed
			callback(false)
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("confirm_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// createModal creates a standard modal with consistent styling
func (mm *ModalManager) createModal(text string, buttons []string) *tview.Modal {
	return tview.NewModal().
		SetText(text).
		AddButtons(buttons)
}

// buildHelpText constructs the help text content
func (mm *ModalManager) buildHelpText() string {
	helpSections := []struct {
		title   string
		content []string
	}{
		{
			title: "Global",
			content: []string{
				"ESC       Close modal",
				"Ctrl+C    Exit application",
				"Q         Exit application",
				"F5        Refresh",
				"?         Show help",
			},
		},
		{
			title: "Navigation",
			content: []string{
				"1, c      Containers view",
				"2, i      Images view",
				"3, v      Volumes view",
				"4, n      Networks view",
			},
		},
		{
			title: "Table Navigation",
			content: []string{
				"↑/↓       Navigate rows",
				"Enter     View details & actions",
				"ESC       Close details",
			},
		},
		{
			title: "Container Actions",
			content: []string{
				"s         Start container",
				"S         Stop container",
				"r         Restart container",
				"d         Delete container",
				"l         View logs",
				"i         Inspect container",
			},
		},
		{
			title: "Image Actions",
			content: []string{
				"d         Delete image",
				"i         Inspect image",
			},
		},
		{
			title: "Volume Actions",
			content: []string{
				"d         Delete volume",
				"i         Inspect volume",
			},
		},
		{
			title: "Network Actions",
			content: []string{
				"d         Delete network",
				"i         Inspect network",
			},
		},
		{
			title: "Configuration",
			content: []string{
				":         Command mode",
				"theme     Custom themes (YAML/JSON)",
				"refresh   Auto-refresh settings",
			},
		},
	}

	helpText := "DockerK9s Keyboard Shortcuts\n\n"
	for _, section := range helpSections {
		helpText += section.title + ":\n"
		for _, item := range section.content {
			helpText += "  " + item + "\n"
		}
		helpText += "\n"
	}

	return helpText
}
