package views

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContainersView(t *testing.T) {
	ui := NewMockUI()
	containersView := NewContainersView(ui)

	require.NotNil(t, containersView)
	// Note: ui field is internal, we can't test it directly
	assert.NotNil(t, containersView.view)
	assert.NotNil(t, containersView.table)
	assert.Empty(t, containersView.items)
	assert.IsType(t, &tview.Flex{}, containersView.view)
	assert.IsType(t, &tview.Table{}, containersView.table)
}

func TestContainersView_GetView(t *testing.T) {
	ui := NewMockUI()
	containersView := NewContainersView(ui)
	view := containersView.GetView()

	assert.NotNil(t, view)
	assert.Equal(t, containersView.view, view)
}

func TestContainersView_Refresh(t *testing.T) {
	ui := NewMockUI()
	containersView := NewContainersView(ui)

	containersView.Refresh()

	assert.Empty(t, containersView.items)
}

func TestContainersView_TableStructure(t *testing.T) {
	ui := NewMockUI()
	containersView := NewContainersView(ui)
	table := containersView.table

	assert.NotNil(t, table)
	// Note: We can't easily test the exact content without a full tview application
	assert.NotNil(t, table)
}

func TestContainersView_EmptyState(t *testing.T) {
	ui := NewMockUI()
	containersView := NewContainersView(ui)

	assert.Empty(t, containersView.items)
	assert.NotNil(t, containersView.table)
	assert.NotNil(t, containersView.view)
}
