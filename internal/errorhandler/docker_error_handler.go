package errorhandler

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

// UserInteractor defines the interface for user interaction
type UserInteractor interface {
	AskYesNo(question string) bool
	WaitForEnter()
}

type DockerErrorHandler struct {
	err         error
	cfg         *config.Config
	interaction UserInteractor
	appRunner   AppRunner
}

func NewDockerErrorHandler(
	err error,
	cfg *config.Config,
	appRunner AppRunner,
	interaction UserInteractor,
) *DockerErrorHandler {
	return &DockerErrorHandler{
		err:         err,
		cfg:         cfg,
		interaction: interaction,
		appRunner:   appRunner,
	}
}

func (h *DockerErrorHandler) Handle() error {
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
func (h *DockerErrorHandler) handleErrorDetails() {
	if h.askForDetails() {
		h.showDetailedError()
	}
}

// handleConnectionType handles different connection types
func (h *DockerErrorHandler) handleConnectionType() {
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

func (h *DockerErrorHandler) prepareTerminal() {
	logger.SetTUIMode(false)
	if _, err := fmt.Fprint(os.Stdout, "\033[2J\033[H"); err != nil {
		// Log the error but continue since this is just terminal clearing
		fmt.Fprintf(os.Stderr, "Failed to clear terminal: %v\n", err)
	}
}

func (h *DockerErrorHandler) showErrorHeader() {
	fmt.Println("❌ Docker Connection Failed")
	fmt.Println("==========================")
	fmt.Println()
	fmt.Printf("Error details: %v\n", h.err)
	fmt.Println()
}

func (h *DockerErrorHandler) askForDetails() bool {
	return h.interaction.AskYesNo("Show detailed error information?")
}

func (h *DockerErrorHandler) askForRetry() bool {
	return h.interaction.AskYesNo("Would you like to retry?")
}

func (h *DockerErrorHandler) isRemoteConnection() bool {
	return h.cfg.RemoteHost != ""
}

func (h *DockerErrorHandler) showDetailedError() {
	fmt.Println()
	fmt.Println("Detailed error information:")
	fmt.Println("==========================")
	fmt.Printf("Error type: %T\n", h.err)
	fmt.Printf("Error message: %s\n", h.err.Error())

	h.showExitErrorDetails()
	h.showSpecificErrorGuidance()
	fmt.Println()
}

func (h *DockerErrorHandler) showExitErrorDetails() {
	if dockerErr, ok := h.err.(*exec.ExitError); ok {
		fmt.Printf("Exit code: %d\n", dockerErr.ExitCode())
		if len(dockerErr.Stderr) > 0 {
			fmt.Printf("Stderr: %s\n", string(dockerErr.Stderr))
		}
	}
}

func (h *DockerErrorHandler) attemptRetry() error {
	fmt.Println("Retrying connection...")
	time.Sleep(2 * time.Second)

	application, retryErr := app.New(h.cfg)
	if retryErr == nil {
		fmt.Println("✅ Connection successful! Starting whaletui...")
		time.Sleep(1 * time.Second)
		return h.appRunner(application, logger.GetLogger())
	}

	fmt.Printf("❌ Retry failed: %v\n", retryErr)
	h.interaction.WaitForEnter()
	return retryErr
}

func (h *DockerErrorHandler) waitForExit() {
	h.interaction.WaitForEnter()
}
