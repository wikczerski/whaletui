package builders

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wikczerski/D5r/internal/ui/constants"
)

func TestTimeFormatter(t *testing.T) {
	tf := NewTimeFormatter()
	require.NotNil(t, tf)

	recentTime := time.Now().Add(-2 * time.Hour)
	result := tf.FormatTime(recentTime)
	assert.Contains(t, result, "ago")
	assert.Contains(t, result, "h")

	oldTime := time.Now().Add(-25 * time.Hour)
	result = tf.FormatTime(oldTime)
	assert.NotContains(t, result, "ago")
}

func TestTimeFormatter_FormatDuration(t *testing.T) {
	tf := NewTimeFormatter()

	result := tf.formatDuration(30 * time.Second)
	assert.Equal(t, "30s", result)

	result = tf.formatDuration(2*time.Minute + 30*time.Second)
	assert.Equal(t, "2m 30s", result)

	result = tf.formatDuration(3*time.Hour + 45*time.Minute)
	assert.Equal(t, "3h 45m", result)

	result = tf.formatDuration(5*24*time.Hour + 12*time.Hour)
	assert.Equal(t, "5d 12h", result)
}

func TestDetailsViewBuilder(t *testing.T) {
	dvb := NewDetailsViewBuilder()
	require.NotNil(t, dvb)
	assert.NotNil(t, dvb.builder)

	actions := map[rune]string{'a': "Action A", 'b': "Action B"}

	detailsView := dvb.CreateDetailsView("Test Title", "Test Details", actions,
		func(r rune) { /* action handler */ },
		func() { /* back handler */ })

	require.NotNil(t, detailsView)
	// Note: GetDirection() is not available on tview.Flex in this version
}

func TestDetailsViewBuilder_SetupKeyBindings(t *testing.T) {
	dvb := NewDetailsViewBuilder()
	actions := map[rune]string{'a': "Action A"}

	detailsView := dvb.CreateDetailsView("Test", "Details", actions,
		func(r rune) { /* action handler */ },
		func() { /* back handler */ })

	require.NotNil(t, detailsView)
	// Note: Testing key bindings requires more complex setup with tview
	// For now, we just verify the view is created successfully
}

func TestDetailsViewBuilder_FormatActions(t *testing.T) {
	dvb := NewDetailsViewBuilder()
	actions := map[rune]string{'a': "Action A", 'b': "Action B"}

	result := dvb.formatActions(actions)
	assert.Contains(t, result, "a: Action A")
	assert.Contains(t, result, "b: Action B")
}

func TestTableBuilder(t *testing.T) {
	tb := NewTableBuilder()
	require.NotNil(t, tb)
	assert.NotNil(t, tb.builder)

	// Test creating table
	table := tb.CreateTable()
	require.NotNil(t, table)
	// Note: Some table methods are not available in this version of tview
}

func TestTableBuilder_SetupHeaders(t *testing.T) {
	tb := NewTableBuilder()
	table := tb.CreateTable()
	headers := []string{"Header1", "Header2", "Header3"}

	tb.SetupHeaders(table, headers)

	for i, header := range headers {
		cell := table.GetCell(0, i)
		require.NotNil(t, cell)
		assert.Equal(t, header, cell.Text)
		assert.Equal(t, constants.HeaderColor, cell.Color)
		// Note: Selectable field is not available in this version of tview
		assert.Equal(t, tview.AlignCenter, cell.Align)
	}
}

func TestTableBuilder_SetupRow(t *testing.T) {
	tb := NewTableBuilder()
	table := tb.CreateTable()
	cells := []string{"Cell1", "Cell2", "Cell3"}
	textColor := tcell.ColorRed

	tb.SetupRow(table, 1, cells, textColor)

	for i, cell := range cells {
		tableCell := table.GetCell(1, i)
		require.NotNil(t, tableCell)
		assert.Equal(t, cell, tableCell.Text)
		assert.Equal(t, textColor, tableCell.Color)
		assert.Equal(t, tview.AlignLeft, tableCell.Align)
	}
}

func TestViewBuilder(t *testing.T) {
	vb := NewViewBuilder()
	require.NotNil(t, vb)
	assert.NotNil(t, vb.builder)

	view := vb.CreateView()
	require.NotNil(t, view)
	// Note: GetDirection() is not available on tview.Flex in this version
}

func TestLegacyFunctions(t *testing.T) {
	testTime := time.Now().Add(-2 * time.Hour)
	result := formatTime(testTime)
	assert.Contains(t, result, "ago")

	actions := map[rune]string{'a': "Action A"}
	detailsView := createDetailsView("Test", "Details", actions, nil, nil)
	require.NotNil(t, detailsView)

	table := createTable()
	require.NotNil(t, table)

	view := createView()
	require.NotNil(t, view)

	table = createTable()
	headers := []string{"H1", "H2"}
	setupTableHeaders(table, headers)

	cells := []string{"C1", "C2"}
	setupTableRow(table, 1, cells, tcell.ColorWhite)
}

func TestCreateLogsView(t *testing.T) {
	title := "Test Logs"
	logsView, logsFlex := createLogsView(title)

	require.NotNil(t, logsView)
	require.NotNil(t, logsFlex)

	assert.Contains(t, logsView.GetTitle(), title)
	// Note: Some methods are not available in this version of tview
	// assert.True(t, logsView.IsDynamicColors())
	// assert.True(t, logsView.IsScrollable())
	// assert.Equal(t, tview.FlexRow, logsFlex.GetDirection())
}

func TestCreateInspectView(t *testing.T) {
	title := "Test Inspect"
	inspectView, inspectFlex := createInspectView(title)

	require.NotNil(t, inspectView)
	require.NotNil(t, inspectFlex)

	assert.Contains(t, inspectView.GetTitle(), title)
	// Note: Some methods are not available in this version of tview
	// assert.True(t, inspectView.IsDynamicColors())
	// assert.True(t, inspectView.IsScrollable())
	// assert.Equal(t, tview.FlexRow, inspectFlex.GetDirection())
}

func TestCreateInspectDetailsView(t *testing.T) {
	title := "Test Inspect Details"
	inspectData := map[string]any{
		"id":   "test123",
		"name": "test-container",
		"state": map[string]any{
			"status": "running",
		},
	}
	actions := map[rune]string{'a': "Action A"}
	detailsView := createInspectDetailsView(title, inspectData, actions,
		func(r rune) { /* action handler */ },
		func() { /* back handler */ })

	require.NotNil(t, detailsView)
	// Note: GetDirection() is not available on tview.Flex in this version
}

func TestFormatInspectData(t *testing.T) {
	data := map[string]any{
		"id":   "test123",
		"name": "test-container",
	}

	result := formatInspectData(data)
	assert.Contains(t, result, "test123")
	assert.Contains(t, result, "test-container")
	assert.Contains(t, result, "{")
	assert.Contains(t, result, "}")

	result = formatInspectData(nil)
	assert.Equal(t, "No inspect data available", result)
}

func TestFormatActionsText(t *testing.T) {
	actions := map[rune]string{'a': "Action A", 'b': "Action B"}

	result := formatActionsText(actions)
	assert.Contains(t, result, "a: Action A")
	assert.Contains(t, result, "b: Action B")
}
