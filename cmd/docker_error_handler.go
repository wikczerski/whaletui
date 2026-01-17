package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/wikczerski/whaletui/internal/app"
	"github.com/wikczerski/whaletui/internal/config"
	"github.com/wikczerski/whaletui/internal/logger"
)

// AppRunner defines the function signature for running the application
type AppRunner func(application *app.App, log *slog.Logger) error

type dockerErrorHandler struct {
	err         error
	cfg         *config.Config
	interaction UserInteraction
	appRunner   AppRunner
}

func newDockerErrorHandler(err error, cfg *config.Config, appRunner AppRunner) *dockerErrorHandler {
	return &dockerErrorHandler{
		err:         err,
		cfg:         cfg,
		interaction: UserInteraction{},
		appRunner:   appRunner,
	}
}

func (h *dockerErrorHandler) handle() error {
	h.prepareTerminal()
	h.showErrorHeader()

	h.handleErrorDetails()
	h.handleConnectionType()

	h.showGeneralHelp()
	h.showLogOption()
	h.waitForExit()
	h.showRemoteOption()

	return fmt.Errorf("docker connection failed: %w", h.err)
}

// handleErrorDetails handles the error details display
func (h *dockerErrorHandler) handleErrorDetails() {
	if h.askForDetails() {
		h.showDetailedError()
	}
}

// handleConnectionType handles different connection types
func (h *dockerErrorHandler) handleConnectionType() {
	if h.isRemoteConnection() {
		h.showRemoteHelp()
	} else {
		h.showLocalHelp()
		if h.askForRetry() {
			if err := h.attemptRetry(); err != nil {
				fmt.Fprintf(os.Stderr, "Retry attempt failed: %v\n", err)
			}
		}
	}
}

func (h *dockerErrorHandler) prepareTerminal() {
	logger.SetTUIMode(false)
	if _, err := fmt.Fprint(os.Stdout, "\033[2J\033[H"); err != nil {
		// Log the error but continue since this is just terminal clearing
		fmt.Fprintf(os.Stderr, "Failed to clear terminal: %v\n", err)
	}
}

func (h *dockerErrorHandler) showErrorHeader() {
	fmt.Println("❌ Docker Connection Failed")
	fmt.Println("==========================")
	fmt.Println()
	fmt.Printf("Error details: %v\n", h.err)
	fmt.Println()
}

func (h *dockerErrorHandler) askForDetails() bool {
	return h.interaction.askYesNo("Show detailed error information?")
}

func (h *dockerErrorHandler) askForRetry() bool {
	return h.interaction.askYesNo("Would you like to retry?")
}

func (h *dockerErrorHandler) isRemoteConnection() bool {
	return h.cfg.RemoteHost != ""
}

func (h *dockerErrorHandler) showDetailedError() {
	fmt.Println()
	fmt.Println("Detailed error information:")
	fmt.Println("==========================")
	fmt.Printf("Error type: %T\n", h.err)
	fmt.Printf("Error message: %s\n", h.err.Error())

	h.showExitErrorDetails()
	h.showSpecificErrorGuidance()
	fmt.Println()
}

func (h *dockerErrorHandler) showExitErrorDetails() {
	if dockerErr, ok := h.err.(*exec.ExitError); ok {
		fmt.Printf("Exit code: %d\n", dockerErr.ExitCode())
		if len(dockerErr.Stderr) > 0 {
			fmt.Printf("Stderr: %s\n", string(dockerErr.Stderr))
		}
	}
}

func (h *dockerErrorHandler) attemptRetry() error {
	fmt.Println("Retrying connection...")
	time.Sleep(2 * time.Second)

	application, retryErr := app.New(h.cfg)
	if retryErr == nil {
		fmt.Println("✅ Connection successful! Starting whaletui...")
		time.Sleep(1 * time.Second)
		return h.appRunner(application, logger.GetLogger())
	}

	fmt.Printf("❌ Retry failed: %v\n", retryErr)
	h.interaction.waitForEnter()
	return retryErr
}

func (h *dockerErrorHandler) waitForExit() {
	h.interaction.waitForEnter()
}
