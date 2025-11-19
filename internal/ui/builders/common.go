package builders

import (
	"github.com/rivo/tview"
)

// ViewBuilder provides methods to create views
type ViewBuilder struct {
	builder *ComponentBuilder
}

// NewViewBuilder creates a new view builder
func NewViewBuilder() *ViewBuilder {
	return &ViewBuilder{
		builder: NewComponentBuilder(),
	}
}

// CreateView creates a new view with consistent setup
func (vb *ViewBuilder) CreateView() *tview.Flex {
	return vb.builder.CreateFlex(tview.FlexRow)
}
