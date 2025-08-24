package managers

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/wikczerski/whaletui/internal/ui/interfaces"
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

	// Add done function to handle Close button click
	modal.SetDoneFunc(func(_ int, _ string) {
		pages := mm.ui.GetPages().(*tview.Pages)
		pages.RemovePage("help_modal")
		// Restore focus to the main view after closing modal
		if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
			if vc, ok := viewContainer.(*tview.Flex); ok {
				app := mm.ui.GetApp().(*tview.Application)
				app.SetFocus(vc)
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("help_modal")
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("help_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowError displays an error modal
func (mm *ModalManager) ShowError(err error) {
	errorText := fmt.Sprintf("Error: %v", err)
	modal := mm.createModal(errorText, []string{"OK"})

	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		pages := mm.ui.GetPages().(*tview.Pages)
		pages.RemovePage("error_modal")
		// Restore focus to the main view after closing modal
		if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
			if vc, ok := viewContainer.(*tview.Flex); ok {
				app := mm.ui.GetApp().(*tview.Application)
				app.SetFocus(vc)
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("error_modal")
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("error_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowInfo displays an info modal
func (mm *ModalManager) ShowInfo(message string) {
	modal := mm.createModal(message, []string{"OK"})

	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		pages := mm.ui.GetPages().(*tview.Pages)
		pages.RemovePage("info_modal")
		// Restore focus to the main view after closing modal
		if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
			if vc, ok := viewContainer.(*tview.Flex); ok {
				app := mm.ui.GetApp().(*tview.Application)
				app.SetFocus(vc)
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("info_modal")
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("info_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowContextualHelp displays context-sensitive help modal
func (mm *ModalManager) ShowContextualHelp(context, operation string) {
	helpContent := mm.generateContextualHelp(context, operation)
	modal := mm.createModal(helpContent, []string{"OK"})

	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		pages := mm.ui.GetPages().(*tview.Pages)
		pages.RemovePage("contextual_help_modal")
		// Restore focus to the main view after closing modal
		if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
			if vc, ok := viewContainer.(*tview.Flex); ok {
				app := mm.ui.GetApp().(*tview.Application)
				app.SetFocus(vc)
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("contextual_help_modal")
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("contextual_help_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowConfirm displays a confirmation modal with Yes/No buttons
func (mm *ModalManager) ShowConfirm(text string, callback func(bool)) {
	modal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, _ string) {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("confirm_modal")
			callback(buttonIndex == 0)
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		})

	// Add keyboard handling for ESC key to close modal (cancel action)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("confirm_modal")
			// Call callback with false (No) when ESC is pressed
			callback(false)
			// Restore focus to the main view after closing modal
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("confirm_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowServiceScaleModal displays a modal for scaling a service
func (mm *ModalManager) ShowServiceScaleModal(serviceName string, currentReplicas uint64, onConfirm func(int)) {
	// Create input field for replicas
	inputField := tview.NewInputField().
		SetLabel("Replicas: ").
		SetText(fmt.Sprintf("%d", currentReplicas)).
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)

	// Create form with input and buttons including help
	form := tview.NewForm().
		AddFormItem(inputField).
		AddButton("Scale", func() {
			// Parse replicas from input
			replicasStr := inputField.GetText()
			replicas, err := strconv.Atoi(replicasStr)
			if err != nil || replicas < 0 {
				mm.ShowError(fmt.Errorf("invalid replicas value: %s", replicasStr))
				return
			}

			// Close modal and call callback
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("scale_modal")
			onConfirm(replicas)

			// Restore focus to main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		}).
		AddButton("Help", func() {
			// Show contextual help for service scaling
			mm.ShowContextualHelp("swarm_services", "scale")
		}).
		AddButton("Cancel", func() {
			// Close modal without action
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("scale_modal")

			// Restore focus to main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		})

	// Create modal container
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Scale Service: %s\nCurrent Replicas: %d", serviceName, currentReplicas)).
		SetDoneFunc(func(_ int, _ string) {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("scale_modal")
		})

	// Create a flex container to hold both modal and form
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(modal, 0, 1, false).
		AddItem(form, 0, 1, true)

	// Add keyboard handling for ESC key
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("scale_modal")
			return nil
		}
		return event
	})

	// Add the modal to the pages
	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("scale_modal", flex, true, true)

	// Set focus to the form
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(form)
}

// ShowNodeAvailabilityModal displays a modal for updating node availability
func (mm *ModalManager) ShowNodeAvailabilityModal(nodeName, currentAvailability string, onConfirm func(string)) {
	// Create the modal content
	content := fmt.Sprintf("Update Node Availability: %s\n\nCurrent Availability: %s\n\nSelect new availability:", nodeName, currentAvailability)

	// Create modal with help button
	modal := mm.createModal(content, []string{"Active", "Pause", "Drain", "Help", "Cancel"})

	// Add done function to handle button clicks
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		switch buttonLabel {
		case "Active":
			onConfirm("active")
			// Close the modal
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("availability_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		case "Pause":
			onConfirm("pause")
			// Close the modal
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("availability_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		case "Drain":
			onConfirm("drain")
			// Close the modal
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("availability_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		case "Help":
			// Show contextual help for node availability updates
			mm.ShowContextualHelp("swarm_nodes", "update_availability")
		case "Cancel":
			// Close the modal without action
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("availability_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("availability_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	// Add the modal to the pages
	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("availability_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowRetryDialog displays a retry dialog with automatic retry logic
func (mm *ModalManager) ShowRetryDialog(operation string, err error, retryFunc func() error, onSuccess func()) {
	// Create retry dialog content
	content := fmt.Sprintf("🔄 Operation Failed: %s\n\nError: %v\n\nThis may be a temporary issue. Would you like to retry?", operation, err)

	// Create modal with retry options
	modal := mm.createModal(content, []string{"Retry", "Retry (Auto)", "Cancel"})

	// Add done function to handle button clicks
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		switch buttonLabel {
		case "Retry":
			// Manual retry - close dialog and let user retry
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("retry_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		case "Retry (Auto)":
			// Automatic retry with progress indication
			mm.performAutomaticRetry(operation, retryFunc, onSuccess)
		case "Cancel":
			// Close dialog without retry
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("retry_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("retry_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	// Add the modal to the pages
	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("retry_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// ShowFallbackDialog displays a fallback operations dialog
func (mm *ModalManager) ShowFallbackDialog(operation string, err error, fallbackOptions []string, onFallback func(string)) {
	// Create fallback dialog content
	content := fmt.Sprintf("⚠️  Operation Failed: %s\n\nError: %v\n\nAlternative operations are available:", operation, err)

	// Create buttons for fallback options
	buttons := make([]string, len(fallbackOptions)+1)
	copy(buttons, fallbackOptions)
	buttons[len(fallbackOptions)] = "Cancel"
	modal := mm.createModal(content, buttons)

	// Add done function to handle button clicks
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			// Close dialog without action
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("fallback_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		} else {
			// Execute fallback operation
			onFallback(buttonLabel)
			// Close dialog
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("fallback_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
		}
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("fallback_modal")
			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app := mm.ui.GetApp().(*tview.Application)
					app.SetFocus(vc)
				}
			}
			return nil // Consume the event
		}
		return event
	})

	// Add the modal to the pages
	pages := mm.ui.GetPages().(*tview.Pages)
	pages.AddPage("fallback_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app := mm.ui.GetApp().(*tview.Application)
	app.SetFocus(modal)
}

// createModal creates a standard modal with consistent styling
func (mm *ModalManager) createModal(text string, buttons []string) *tview.Modal {
	return tview.NewModal().
		SetText(text).
		AddButtons(buttons)
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
				"s         Swarm Services view",
				"w         Swarm Nodes view",
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
		{
			title: "Swarm Service Actions",
			content: []string{
				"i         Inspect service",
				"s         Scale service",
				"r         Remove service",
				"l         View logs",
			},
		},
		{
			title: "Swarm Node Actions",
			content: []string{
				"i         Inspect node",
				"a         Update availability",
				"r         Remove node",
			},
		},
		{
			title: "Configuration",
			content: []string{
				":         Command mode",
				"theme     Custom themes (YAML/JSON)",
				"refresh   Auto-refresh settings",
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

// generateContextualHelp creates context-sensitive help content
func (mm *ModalManager) generateContextualHelp(context, operation string) string {
	var helpContent string

	switch context {
	case "swarm_services":
		helpContent = mm.generateSwarmServicesHelp(operation)
	case "swarm_nodes":
		helpContent = mm.generateSwarmNodesHelp(operation)
	case "containers":
		helpContent = mm.generateContainersHelp(operation)
	case "images":
		helpContent = mm.generateImagesHelp(operation)
	case "networks":
		helpContent = mm.generateNetworksHelp(operation)
	case "volumes":
		helpContent = mm.generateVolumesHelp(operation)
	default:
		helpContent = mm.generateGeneralHelp(operation)
	}

	return helpContent
}

// generateSwarmServicesHelp creates help content for swarm services context
func (mm *ModalManager) generateSwarmServicesHelp(operation string) string {
	switch operation {
	case "scale":
		return `🔧 Service Scaling Help

Scaling a service changes the number of replicas running.

What happens when you scale:
• Docker Swarm will start/stop tasks to match the new replica count
• Service remains available during scaling (rolling update)
• Load balancer automatically distributes traffic

Best practices:
• Scale gradually for production services
• Monitor resource usage after scaling
• Consider using auto-scaling for variable workloads

Common issues:
• Insufficient resources on nodes
• Service constraints preventing placement
• Network connectivity issues

Need more help? Check Docker Swarm documentation.`

	case "remove":
		return `⚠️ Service Removal Help

Removing a service will permanently delete it.

What happens when you remove:
• All running tasks are stopped immediately
• Service definition is removed from swarm
• Load balancer stops routing traffic
• Cannot be undone

Before removing:
• Ensure no critical dependencies
• Backup service configuration if needed
• Consider stopping instead of removing

Alternatives to removal:
• Scale to 0 replicas (pause service)
• Update service configuration
• Use service update for changes

Need more help? Check Docker Swarm documentation.`

	case "inspect":
		return `🔍 Service Inspection Help

Inspecting a service shows detailed information.

What you can see:
• Service configuration and settings
• Current replica count and status
• Network and volume mounts
• Environment variables and labels
• Update and rollback history

Useful for:
• Troubleshooting service issues
• Understanding service configuration
• Planning updates or changes
• Debugging network problems

Common inspection fields:
• Spec: Service configuration
• Endpoint: Network endpoints
• UpdateStatus: Update progress
• PreviousSpec: Previous configuration

Need more help? Check Docker Swarm documentation.`

	case "logs":
		return `📋 Service Logs Help

Viewing service logs helps with troubleshooting.

What you can see:
• Application output and errors
• System messages and warnings
• Network connection logs
• Container startup/shutdown events

Log viewing tips:
• Logs may be truncated for performance
• Use Docker CLI for full log access
• Consider log aggregation for production
• Monitor logs for error patterns

Common log issues:
• High log volume affecting performance
• Missing logs due to rotation
• Permission issues accessing logs
• Network connectivity problems

Need more help? Check Docker Swarm documentation.`

	default:
		return `📚 Swarm Services Help

Available operations:
• Scale (s): Change number of replicas
• Remove (r): Delete service permanently
• Inspect (i): View detailed information
• Logs (l): View service logs

Navigation:
• Use arrow keys to select services
• Press 'h' for this help
• Press 'q' to return to main view

Need specific help? Select an operation first.`
	}
}

// generateSwarmNodesHelp creates help content for swarm nodes context
func (mm *ModalManager) generateSwarmNodesHelp(operation string) string {
	switch operation {
	case "update_availability":
		return `🔄 Node Availability Help

Changing node availability affects task placement.

Availability options:
• Active: Accepts new tasks (default)
• Pause: No new tasks, existing tasks continue
• Drain: No new tasks, existing tasks are rescheduled

What happens when draining:
• Running tasks are moved to other nodes
• Service remains available during transition
• Node becomes unavailable for new tasks
• Useful for maintenance or updates

Best practices:
• Drain nodes before maintenance
• Ensure sufficient capacity on other nodes
• Monitor task rescheduling progress
• Use pause for temporary unavailability

Common issues:
• Insufficient capacity on remaining nodes
• Tasks that cannot be rescheduled
• Network connectivity problems
• Resource constraints preventing placement

Need more help? Check Docker Swarm documentation.`

	case "remove":
		return `⚠️ Node Removal Help

Removing a node affects swarm stability.

What happens when you remove:
• Node is forcefully removed from swarm
• All tasks on the node are stopped
• Swarm rebalances remaining tasks
• Node must be re-added to rejoin

⚠️ Important warnings:
• Removing manager nodes affects swarm stability
• Ensure sufficient manager nodes remain
• Consider draining before removal
• Backup swarm state if possible

Before removing:
• Drain the node first (recommended)
• Ensure sufficient capacity remains
• Check manager node count
• Plan for service redistribution

Need more help? Check Docker Swarm documentation.`

	case "inspect":
		return `🔍 Node Inspection Help

Inspecting a node shows detailed information.

What you can see:
• Node status and availability
• Resource usage and capacity
• Engine version and plugins
• Network configuration
• Manager status (if applicable)

Useful for:
• Troubleshooting node issues
• Planning capacity and scaling
• Understanding node configuration
• Debugging network problems

Common inspection fields:
• Status: Node health and readiness
• Availability: Task placement preference
• EngineVersion: Docker engine version
• ManagerStatus: Manager role information

Need more help? Check Docker Swarm documentation.`

	default:
		return `📚 Swarm Nodes Help

Available operations:
• Update Availability (a): Change node availability
• Remove (r): Remove node from swarm
• Inspect (i): View detailed information

Navigation:
• Use arrow keys to select nodes
• Press 'h' for this help
• Press 'q' to return to main view

Need specific help? Select an operation first.`
	}
}

// generateContainersHelp creates help content for containers context
func (mm *ModalManager) generateContainersHelp(_ string) string {
	return `📚 Containers Help

Available operations:
• Start: Start a stopped container
• Stop: Stop a running container
• Remove: Delete a container
• Inspect: View detailed information
• Logs: View container logs

Navigation:
• Use arrow keys to select containers
• Press 'h' for this help
• Press 'q' to return to main view

Need specific help? Select an operation first.`
}

// generateImagesHelp creates help content for images context
func (mm *ModalManager) generateImagesHelp(_ string) string {
	return `📚 Images Help

Available operations:
• Remove: Delete an image
• Inspect: View detailed information
• History: View image layers

Navigation:
• Use arrow keys to select images
• Press 'h' for this help
• Press 'q' to return to main view

Need specific help? Select an operation first.`
}

// generateNetworksHelp creates help content for networks context
func (mm *ModalManager) generateNetworksHelp(_ string) string {
	return `📚 Networks Help

Available operations:
• Remove: Delete a network
• Inspect: View detailed information

Navigation:
• Use arrow keys to select networks
• Press 'h' for this help
• Press 'q' to return to main view

Need specific help? Select an operation first.`
}

// generateVolumesHelp creates help content for volumes context
func (mm *ModalManager) generateVolumesHelp(_ string) string {
	return `📚 Volumes Help

Available operations:
• Remove: Delete a volume
• Inspect: View detailed information

Navigation:
• Use arrow keys to select volumes
• Press 'h' for this help
• Press 'q' to return to main view

Need specific help? Select an operation first.`
}

// generateGeneralHelp creates general help content
func (mm *ModalManager) generateGeneralHelp(_ string) string {
	return `📚 General Help

Available operations:
• Navigate between views
• Manage Docker resources
• View system information

Navigation:
• Use arrow keys to navigate
• Press 'h' for context-specific help
• Press 'q' to return to previous view

Need specific help? Navigate to a specific view first.`
}

// performAutomaticRetry performs automatic retry with progress indication
func (mm *ModalManager) performAutomaticRetry(operation string, retryFunc func() error, onSuccess func()) {
	// Close the retry dialog
	pages := mm.ui.GetPages().(*tview.Pages)
	pages.RemovePage("retry_modal")

	// Show progress modal
	progressContent := fmt.Sprintf("🔄 Retrying: %s\n\nPlease wait while we attempt to recover...", operation)
	progressModal := mm.createModal(progressContent, []string{"Cancel"})

	// Add cancel functionality
	progressModal.SetDoneFunc(func(_ int, _ string) {
		pages := mm.ui.GetPages().(*tview.Pages)
		pages.RemovePage("retry_progress_modal")
		// Restore focus to the main view
		if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
			if vc, ok := viewContainer.(*tview.Flex); ok {
				app := mm.ui.GetApp().(*tview.Application)
				app.SetFocus(vc)
			}
		}
	})

	// Add the progress modal to the pages
	pages.AddPage("retry_progress_modal", progressModal, true, true)

	// Perform retry in a goroutine to avoid blocking UI
	go func() {
		// Attempt retry
		err := retryFunc()

		// Close progress modal from main thread
		app := mm.ui.GetApp().(*tview.Application)
		app.QueueUpdateDraw(func() {
			pages := mm.ui.GetPages().(*tview.Pages)
			pages.RemovePage("retry_progress_modal")

			if err != nil {
				// Retry failed - show error
				mm.ShowError(fmt.Errorf("automatic retry failed for %s: %v", operation, err))
			} else {
				// Retry succeeded - show success and execute callback
				mm.ShowInfo(fmt.Sprintf("✅ Operation '%s' recovered successfully!", operation))
				if onSuccess != nil {
					onSuccess()
				}
			}

			// Restore focus to the main view
			if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
				if vc, ok := viewContainer.(*tview.Flex); ok {
					app.SetFocus(vc)
				}
			}
		})
	}()
}
