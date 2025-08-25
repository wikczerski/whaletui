package handlers

import (
	"context"
	"time"

	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
)

// OperationExecutor handles common async operations with error handling and refresh
type OperationExecutor struct {
	ui           interfaces.UIInterface
	refreshDelay time.Duration
}

// NewOperationExecutor creates a new operation executor
func NewOperationExecutor(ui interfaces.UIInterface) *OperationExecutor {
	return &OperationExecutor{
		ui:           ui,
		refreshDelay: 500 * time.Millisecond,
	}
}

// Execute runs an operation asynchronously with error handling and refresh
func (oe *OperationExecutor) Execute(operation func() error, onRefresh func()) {
	go func() {
		if err := operation(); err != nil {
			oe.handleOperationError(err)
			return
		}

		time.Sleep(oe.refreshDelay)
		oe.handleOperationSuccess(onRefresh)
	}()
}

// ExecuteWithConfirmation shows a confirmation dialog before executing an operation
func (oe *OperationExecutor) ExecuteWithConfirmation(
	message string,
	operation func() error,
	onRefresh func(),
) {
	oe.ui.ShowConfirm(message, func(confirmed bool) {
		if !confirmed {
			return
		}
		oe.Execute(operation, onRefresh)
	})
}

// DeleteOperation handles resource deletion with confirmation
func (oe *OperationExecutor) DeleteOperation(
	resourceType, resourceID, resourceName string,
	deleteFunc func(context.Context, string, bool) error,
	onRefresh func(),
) {
	message := "Delete " + resourceType + " " + resourceName + "?"
	operation := func() error {
		return deleteFunc(context.Background(), resourceID, true)
	}
	oe.ExecuteWithConfirmation(message, operation, onRefresh)
}

// StartOperation handles resource startup
func (oe *OperationExecutor) StartOperation(
	_, resourceID string,
	startFunc func(context.Context, string) error,
	onRefresh func(),
) {
	operation := func() error {
		return startFunc(context.Background(), resourceID)
	}
	oe.Execute(operation, onRefresh)
}

// StopOperation handles resource shutdown
func (oe *OperationExecutor) StopOperation(
	_, resourceID string,
	stopFunc func(context.Context, string, *time.Duration) error,
	onRefresh func(),
) {
	operation := func() error {
		timeout := 10 * time.Second
		return stopFunc(context.Background(), resourceID, &timeout)
	}
	oe.Execute(operation, onRefresh)
}

// RestartOperation handles resource restart
func (oe *OperationExecutor) RestartOperation(
	_, resourceID string,
	restartFunc func(context.Context, string, *time.Duration) error,
	onRefresh func(),
) {
	operation := func() error {
		timeout := 10 * time.Second
		return restartFunc(context.Background(), resourceID, &timeout)
	}
	oe.Execute(operation, onRefresh)
}

// handleOperationError handles operation errors by showing them in the UI
func (oe *OperationExecutor) handleOperationError(err error) {
	app, ok := oe.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.QueueUpdateDraw(func() {
		oe.ui.ShowError(err)
	})
}

// handleOperationSuccess handles successful operations by refreshing the UI
func (oe *OperationExecutor) handleOperationSuccess(onRefresh func()) {
	app, ok := oe.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.QueueUpdateDraw(func() {
		if onRefresh != nil {
			onRefresh()
		}
	})
}
