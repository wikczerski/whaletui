package builders

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/constants"
	"github.com/wikczerski/whaletui/internal/ui/utils"
)

// TableBuilder provides methods to create and configure tables
type TableBuilder struct {
	builder *ComponentBuilder
}

// NewTableBuilder creates a new table builder
func NewTableBuilder() *TableBuilder {
	return &TableBuilder{
		builder: NewComponentBuilder(),
	}
}

// CreateTable creates a new table with consistent styling
func (tb *TableBuilder) CreateTable() *tview.Table {
	return tb.builder.CreateTable()
}

// SetupHeaders sets up table headers with consistent styling
func (tb *TableBuilder) SetupHeaders(table *tview.Table, headers []string) {
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(constants.HeaderColor).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetExpansion(1)
		table.SetCell(0, i, cell)
	}
}

// SetupHeadersWithConfig sets up table headers with column configuration
func (tb *TableBuilder) SetupHeadersWithConfig(
	table *tview.Table,
	headers []string,
	columnTypes []string,
	formatter *utils.TableFormatter,
) {
	tb.SetupHeadersWithConfigForView(table, headers, columnTypes, formatter, "")
}

// SetupHeadersWithConfigForView sets up table headers with view-specific column configuration
func (tb *TableBuilder) SetupHeadersWithConfigForView(
	table *tview.Table,
	headers []string,
	columnTypes []string,
	formatter *utils.TableFormatter,
	viewName string,
) {
	tb.SetupHeadersWithConfigForViewWithTerminalSize(
		table, headers, columnTypes, formatter, viewName, 0)
}

// SetupHeadersWithConfigForViewWithTerminalSize sets up table headers with view-specific column configuration
// with a given terminal width for percentage-based width calculations
func (tb *TableBuilder) SetupHeadersWithConfigForViewWithTerminalSize(
	table *tview.Table,
	headers []string,
	columnTypes []string,
	formatter *utils.TableFormatter,
	viewName string,
	terminalWidth int,
) {
	visibleColumnIndex := 0
	for i, header := range headers {
		// Skip hidden columns
		if formatter != nil && i < len(columnTypes) &&
			!formatter.IsColumnVisibleForView(columnTypes[i], viewName) {
			continue
		}

		// Use display name if available
		displayHeader := header
		if formatter != nil && i < len(columnTypes) {
			displayHeader = formatter.GetColumnDisplayNameForView(columnTypes[i], viewName)
		}

		cell := tview.NewTableCell(displayHeader).
			SetTextColor(constants.HeaderColor).
			SetSelectable(false).
			SetAlign(tview.AlignCenter).
			SetExpansion(1)

		// Set fixed width if configured
		if formatter != nil && i < len(columnTypes) {
			tb.applyHeaderWidth(displayHeader, cell, formatter,
				columnTypes[i], viewName, terminalWidth, table)
		}

		table.SetCell(0, visibleColumnIndex, cell)
		visibleColumnIndex++
	}
}

// SetupRow sets up a table row with consistent styling
func (tb *TableBuilder) SetupRow(
	table *tview.Table,
	row int,
	cells []string,
	textColor tcell.Color,
) {
	for i, cell := range cells {
		tableCell := tview.NewTableCell(cell).
			SetTextColor(textColor).
			SetAlign(tview.AlignLeft).
			SetExpansion(1)
		table.SetCell(row, i, tableCell)
	}
}

// SetupRowWithAlignment sets up a table row with alignment based on column types
func (tb *TableBuilder) SetupRowWithAlignment(
	table *tview.Table,
	row int,
	cells []string,
	columnTypes []string,
	textColor tcell.Color,
	formatter *utils.TableFormatter,
) {
	visibleColumnIndex := 0
	for i, cell := range cells {
		// Skip hidden columns
		if formatter != nil && i < len(columnTypes) && !formatter.IsColumnVisible(columnTypes[i]) {
			continue
		}

		alignment := tview.AlignLeft // Default alignment

		if formatter != nil && i < len(columnTypes) {
			alignment = formatter.GetAlignmentForColumn(columnTypes[i])
		}

		tableCell := tview.NewTableCell(cell).
			SetTextColor(textColor).
			SetAlign(alignment).
			SetExpansion(1)

		// Set fixed width if configured
		if formatter != nil && i < len(columnTypes) {
			width := formatter.GetColumnWidth(columnTypes[i])
			if width > 0 {
				tableCell.SetExpansion(0)
			}
		}

		table.SetCell(row, visibleColumnIndex, tableCell)
		visibleColumnIndex++
	}
}

// SetupRowWithLimits sets up a table row with character limits and alignment applied
func (tb *TableBuilder) SetupRowWithLimits(
	table *tview.Table,
	row int,
	cells []string,
	columnTypes []string,
	textColor tcell.Color,
	formatter *utils.TableFormatter,
) {
	tb.SetupRowWithLimitsForView(table, row, cells, columnTypes, textColor, formatter, "")
}

// SetupRowWithLimitsForView sets up a table row with view-specific character limits and alignment applied
func (tb *TableBuilder) SetupRowWithLimitsForView(
	table *tview.Table,
	row int,
	cells []string,
	columnTypes []string,
	textColor tcell.Color,
	formatter *utils.TableFormatter,
	viewName string,
) {
	tb.SetupRowWithLimitsForViewWithTerminalSize(table, row, cells, columnTypes,
		textColor, formatter, viewName, 0)
}

// SetupRowWithLimitsForViewWithTerminalSize sets up a table row with view-specific character limits
// and alignment applied
// with a given terminal width for percentage-based width calculations
func (tb *TableBuilder) SetupRowWithLimitsForViewWithTerminalSize(
	table *tview.Table,
	row int,
	cells []string,
	columnTypes []string,
	textColor tcell.Color,
	formatter *utils.TableFormatter,
	viewName string,
	terminalWidth int,
) {
	visibleColumnIndex := 0
	for i, cell := range cells {
		// Skip hidden columns
		if formatter != nil && i < len(columnTypes) &&
			!formatter.IsColumnVisibleForView(columnTypes[i], viewName) {
			continue
		}

		// Apply character limits if formatter and column type are provided
		formattedCell := cell
		alignment := tview.AlignLeft // Default alignment

		if formatter != nil && i < len(columnTypes) {
			formattedCell = formatter.FormatCellForView(cell, columnTypes[i], viewName)
			alignment = formatter.GetAlignmentForColumnForView(columnTypes[i], viewName)
		}

		tableCell := tview.NewTableCell(formattedCell).
			SetTextColor(textColor).
			SetAlign(alignment).
			SetExpansion(1)

		// Set fixed width if configured
		if formatter != nil && i < len(columnTypes) {
			tb.applyRowWidth(formattedCell, tableCell, formatter,
				columnTypes[i], viewName, terminalWidth, table, alignment)
		}

		table.SetCell(row, visibleColumnIndex, tableCell)
		visibleColumnIndex++
	}
}

// getTerminalWidth gets the terminal width
func (tb *TableBuilder) getTerminalWidth(table *tview.Table) int {
	_, _, width, _ := table.GetRect()
	if width > 0 {
		return width
	}
	return 120 // Default if not yet visible
}

// applyHeaderWidth applies width configuration to a header cell
func (tb *TableBuilder) applyHeaderWidth(displayHeader string, cell *tview.TableCell,
	formatter *utils.TableFormatter, columnType, viewName string,
	terminalWidth int, table *tview.Table,
) {
	// Get terminal width if not provided
	if terminalWidth <= 0 {
		terminalWidth = tb.getTerminalWidth(table)
	}
	width := formatter.GetColumnWidthForViewWithTerminalSize(
		columnType, viewName, terminalWidth)
	if width > 0 {
		cell.SetExpansion(0)
		// Pad the header content to the desired width
		if len(displayHeader) < width {
			// Center-pad the header with spaces to reach desired width
			padding := width - len(displayHeader)
			leftPad := padding / 2
			rightPad := padding - leftPad
			displayHeader = fmt.Sprintf("%*s%s%*s", leftPad, "", displayHeader,
				rightPad, "")
		}
		// Update the cell with padded content
		cell.SetText(displayHeader)
	}
}

// applyRowWidth applies width configuration to a row cell
func (tb *TableBuilder) applyRowWidth(formattedCell string, tableCell *tview.TableCell,
	formatter *utils.TableFormatter, columnType, viewName string,
	terminalWidth int, table *tview.Table, alignment int,
) {
	// Get terminal width if not provided
	if terminalWidth <= 0 {
		terminalWidth = tb.getTerminalWidth(table)
	}
	width := formatter.GetColumnWidthForViewWithTerminalSize(
		columnType, viewName, terminalWidth)
	if width > 0 {
		tableCell.SetExpansion(0)
		// Pad the content to the desired width based on alignment
		if len(formattedCell) < width {
			switch alignment {
			case tview.AlignRight:
				// Left-pad with spaces for right alignment
				formattedCell = fmt.Sprintf("%*s", width, formattedCell)
			case tview.AlignCenter:
				// Center-pad with spaces
				padding := width - len(formattedCell)
				leftPad := padding / 2
				rightPad := padding - leftPad
				formattedCell = fmt.Sprintf("%*s%s%*s", leftPad, "", formattedCell, rightPad, "")
			default: // tview.AlignLeft
				// Right-pad with spaces for left alignment
				formattedCell = fmt.Sprintf("%-*s", width, formattedCell)
			}
		}
		// Update the cell with padded content
		tableCell.SetText(formattedCell)
	}
}
