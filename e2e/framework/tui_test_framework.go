package framework

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"
)

// TUITestFramework provides TUI-specific testing utilities using tcell simulation.
type TUITestFramework struct {
	t            *testing.T
	screen       tcell.SimulationScreen
	app          *tview.Application
	cleanupFuncs []func()
}

// NewTUITestFramework creates a new TUI test framework with simulated screen.
func NewTUITestFramework(t *testing.T) *TUITestFramework {
	t.Helper()

	// Create simulation screen
	screen := tcell.NewSimulationScreen("UTF-8")
	err := screen.Init()
	require.NoError(t, err, "Failed to initialize simulation screen")

	// Set screen size (80x24 is standard terminal size)
	screen.SetSize(80, 24)

	fw := &TUITestFramework{
		t:            t,
		screen:       screen,
		cleanupFuncs: make([]func(), 0),
	}

	// Register cleanup
	t.Cleanup(func() {
		fw.TearDown()
	})

	return fw
}

// CreateApp creates a tview application with the simulated screen.
func (fw *TUITestFramework) CreateApp() *tview.Application {
	fw.t.Helper()

	app := tview.NewApplication()
	app.SetScreen(fw.screen)

	fw.app = app

	// Add cleanup
	fw.AddCleanup(func() {
		if fw.app != nil {
			fw.app.Stop()
		}
	})

	return app
}

// StartApp starts the application in a goroutine.
func (fw *TUITestFramework) StartApp(root tview.Primitive) {
	fw.t.Helper()

	if fw.app == nil {
		fw.CreateApp()
	}

	fw.app.SetRoot(root, true)

	// Start app in background
	go func() {
		if err := fw.app.Run(); err != nil {
			fw.t.Logf("App run error: %v", err)
		}
	}()

	// Wait for app to initialize and render
	time.Sleep(200 * time.Millisecond)
	fw.screen.Show()
}

// InjectKey simulates a key press event.
func (fw *TUITestFramework) InjectKey(key tcell.Key, ch rune, mod tcell.ModMask) {
	fw.t.Helper()

	fw.screen.InjectKey(key, ch, mod)

	// Allow time for event processing
	time.Sleep(50 * time.Millisecond)
	fw.screen.Show()
}

// InjectKeyRune simulates a rune key press.
func (fw *TUITestFramework) InjectKeyRune(ch rune) {
	fw.InjectKey(tcell.KeyRune, ch, tcell.ModNone)
}

// InjectKeyPress simulates a special key press.
func (fw *TUITestFramework) InjectKeyPress(key tcell.Key) {
	fw.InjectKey(key, 0, tcell.ModNone)
}

// InjectString simulates typing a string.
func (fw *TUITestFramework) InjectString(s string) {
	fw.t.Helper()

	for _, ch := range s {
		fw.InjectKeyRune(ch)
	}
}

// GetScreenContent returns the current screen content as a string.
func (fw *TUITestFramework) GetScreenContent() string {
	fw.t.Helper()

	width, height := fw.screen.Size()
	var content string

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			mainc, _, _, _ := fw.screen.GetContent(x, y)
			if mainc == 0 {
				mainc = ' '
			}
			content += string(mainc)
		}
		if y < height-1 {
			content += "\n"
		}
	}

	return content
}

// GetScreenLine returns a specific line from the screen.
func (fw *TUITestFramework) GetScreenLine(line int) string {
	fw.t.Helper()

	width, height := fw.screen.Size()
	if line < 0 || line >= height {
		fw.t.Fatalf("Line %d out of bounds (height: %d)", line, height)
	}

	var lineContent string
	for x := 0; x < width; x++ {
		mainc, _, _, _ := fw.screen.GetContent(x, line)
		if mainc == 0 {
			mainc = ' '
		}
		lineContent += string(mainc)
	}

	return lineContent
}

// GetCellContent returns the character at a specific position.
func (fw *TUITestFramework) GetCellContent(x, y int) rune {
	fw.t.Helper()

	mainc, _, _, _ := fw.screen.GetContent(x, y)
	return mainc
}

// GetCellStyle returns the style at a specific position.
func (fw *TUITestFramework) GetCellStyle(x, y int) tcell.Style {
	fw.t.Helper()

	_, _, style, _ := fw.screen.GetContent(x, y)
	return style
}

// VerifyTextAt verifies that specific text appears at a given position.
func (fw *TUITestFramework) VerifyTextAt(x, y int, expectedText string) bool {
	fw.t.Helper()

	for i, ch := range expectedText {
		actualChar := fw.GetCellContent(x+i, y)
		if actualChar != ch {
			fw.t.Logf("Text mismatch at (%d,%d): expected '%c', got '%c'", x+i, y, ch, actualChar)
			return false
		}
	}
	return true
}

// VerifyTextContains verifies that the screen contains specific text.
func (fw *TUITestFramework) VerifyTextContains(expectedText string) bool {
	fw.t.Helper()

	content := fw.GetScreenContent()
	return contains(content, expectedText)
}

// VerifyLineContains verifies that a specific line contains text.
func (fw *TUITestFramework) VerifyLineContains(line int, expectedText string) bool {
	fw.t.Helper()

	lineContent := fw.GetScreenLine(line)
	return contains(lineContent, expectedText)
}

// WaitForText waits for specific text to appear on screen.
func (fw *TUITestFramework) WaitForText(expectedText string, timeout time.Duration) bool {
	fw.t.Helper()

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if fw.VerifyTextContains(expectedText) {
			return true
		}

		if time.Now().After(deadline) {
			fw.t.Logf("Timeout waiting for text: %s", expectedText)
			fw.t.Logf("Screen content:\n%s", fw.GetScreenContent())
			return false
		}

		<-ticker.C
		fw.screen.Show()
	}
}

// WaitForTextAt waits for specific text to appear at a position.
func (fw *TUITestFramework) WaitForTextAt(x, y int, expectedText string, timeout time.Duration) bool {
	fw.t.Helper()

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if fw.VerifyTextAt(x, y, expectedText) {
			return true
		}

		if time.Now().After(deadline) {
			fw.t.Logf("Timeout waiting for text at (%d,%d): %s", x, y, expectedText)
			return false
		}

		<-ticker.C
		fw.screen.Show()
	}
}

// Sync forces a screen update.
func (fw *TUITestFramework) Sync() {
	if fw.app != nil {
		// Force the app to redraw
		fw.app.Draw()
	}
	fw.screen.Sync()
	fw.screen.Show()
	// Give time for the draw to complete
	time.Sleep(50 * time.Millisecond)
}

// Clear clears the screen.
func (fw *TUITestFramework) Clear() {
	fw.screen.Clear()
}

// GetSize returns the screen dimensions.
func (fw *TUITestFramework) GetSize() (width, height int) {
	return fw.screen.Size()
}

// SetSize sets the screen dimensions.
func (fw *TUITestFramework) SetSize(width, height int) {
	fw.screen.SetSize(width, height)
}

// AddCleanup adds a cleanup function.
func (fw *TUITestFramework) AddCleanup(cleanup func()) {
	fw.cleanupFuncs = append(fw.cleanupFuncs, cleanup)
}

// TearDown cleans up resources.
func (fw *TUITestFramework) TearDown() {
	fw.t.Helper()

	// Stop app first
	if fw.app != nil {
		fw.app.Stop()
		time.Sleep(100 * time.Millisecond) // Give app time to stop
	}

	// Run cleanup functions in reverse order
	for i := len(fw.cleanupFuncs) - 1; i >= 0; i-- {
		fw.cleanupFuncs[i]()
	}

	// Finalize screen (only if not already finalized)
	if fw.screen != nil {
		// The screen may already be finalized by the app stopping
		// We'll use a defer/recover to handle this gracefully
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Screen was already finalized, that's okay
					fw.t.Logf("Screen already finalized (expected): %v", r)
				}
			}()
			fw.screen.Fini()
		}()
	}
}

// GetScreen returns the simulation screen.
func (fw *TUITestFramework) GetScreen() tcell.SimulationScreen {
	return fw.screen
}

// GetApp returns the tview application.
func (fw *TUITestFramework) GetApp() *tview.Application {
	return fw.app
}

// DumpScreen dumps the current screen content for debugging.
func (fw *TUITestFramework) DumpScreen() {
	fw.t.Helper()

	content := fw.GetScreenContent()
	fw.t.Logf("Screen dump:\n%s", content)
}

// Helper function to check if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
