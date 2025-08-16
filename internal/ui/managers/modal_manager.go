package managers

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/wikczerski/D5r/internal/ui/interfaces"
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
	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("modal", modal, true, true)
}

// ShowError displays an error modal
func (mm *ModalManager) ShowError(err error) {
	errorText := fmt.Sprintf("Error: %v", err)
	modal := mm.createModal(errorText, []string{"OK"})
	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("modal", modal, true, true)
}

// ShowConfirm displays a confirmation modal with Yes/No buttons
func (mm *ModalManager) ShowConfirm(text string, callback func(bool)) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("modal")
			callback(buttonIndex == 0)
		})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("modal", modal, true, true)
}

// createModal creates a standard modal with consistent styling
func (mm *ModalManager) createModal(text string, buttons []string) *tview.Modal {
	return tview.NewModal().
		SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("modal")
		})
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
