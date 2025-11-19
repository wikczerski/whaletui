package builders

import "github.com/rivo/tview"

// ModalBuilder handles the construction of tview Modals
type ModalBuilder struct {
	modal *tview.Modal
}

// NewModalBuilder creates a new ModalBuilder
func NewModalBuilder() *ModalBuilder {
	return &ModalBuilder{
		modal: tview.NewModal(),
	}
}

// SetText sets the modal text
func (b *ModalBuilder) SetText(text string) *ModalBuilder {
	b.modal.SetText(text)
	return b
}

// AddButtons adds buttons to the modal
func (b *ModalBuilder) AddButtons(buttons []string) *ModalBuilder {
	b.modal.AddButtons(buttons)
	return b
}

// SetDoneFunc sets the handler for button clicks
func (b *ModalBuilder) SetDoneFunc(handler func(buttonIndex int, buttonLabel string)) *ModalBuilder {
	b.modal.SetDoneFunc(handler)
	return b
}

// Build returns the constructed modal
func (b *ModalBuilder) Build() *tview.Modal {
	return b.modal
}
