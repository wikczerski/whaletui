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

	mm.setupHelpModalHandlers(modal)
	mm.addHelpModalToUI(modal)
}

// ShowError displays an error modal
func (mm *ModalManager) ShowError(err error) {
	errorText := fmt.Sprintf("Error: %v", err)
	modal := mm.createModal(errorText, []string{"OK"})

	mm.setupErrorModalHandlers(modal)
	mm.showErrorModal(modal)
}

// ShowInfo displays an info modal
func (mm *ModalManager) ShowInfo(message string) {
	modal := mm.createModal(message, []string{"OK"})

	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		mm.closeInfoModalAndRestoreFocus()
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeInfoModalAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})

	mm.showInfoModal(modal)
}

// ShowContextualHelp displays context-sensitive help modal
func (mm *ModalManager) ShowContextualHelp(context, operation string) {
	helpContent := mm.generateContextualHelp(context, operation)
	modal := mm.createModal(helpContent, []string{"OK"})

	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		mm.closeContextualHelpModalAndRestoreFocus()
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeContextualHelpModalAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})

	mm.addContextualHelpModalToPages(modal)

	// Set focus to the modal so it can receive keyboard input
	mm.setFocusToModal(modal)
}

// ShowConfirm displays a confirmation modal with Yes/No buttons
func (mm *ModalManager) ShowConfirm(text string, callback func(bool)) {
	modal := mm.createConfirmModal(text, callback)
	mm.setupConfirmModalHandlers(modal, callback)
	mm.showConfirmModal(modal)
}

// ShowServiceScaleModal displays a modal for scaling a service
func (mm *ModalManager) ShowServiceScaleModal(
	serviceName string,
	currentReplicas uint64,
	onConfirm func(int),
) {
	// Create input field for replicas
	inputField := mm.createReplicasInputField(currentReplicas)

	// Create form with input and buttons including help
	form := mm.createScaleForm(inputField, onConfirm)

	// Create modal container
	modal := mm.createScaleModal(serviceName, currentReplicas)

	// Create a flex container to hold both modal and form
	flex := mm.createScaleFlexContainer(modal, form)

	// Add the modal to the pages
	mm.addScaleModalToPages(flex)

	// Set focus to the form
	mm.setFocusToForm(form)
}

// ShowNodeAvailabilityModal displays a modal for updating node availability
func (mm *ModalManager) ShowNodeAvailabilityModal(
	nodeName, currentAvailability string,
	onConfirm func(string),
) {
	// Create the modal content
	content := mm.createNodeAvailabilityContent(nodeName, currentAvailability)

	// Create modal with help button
	modal := mm.createModal(content, []string{"Active", "Pause", "Drain", "Help", "Cancel"})

	// Add done function to handle button clicks
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		mm.handleNodeAvailabilityButtonClick(buttonLabel, onConfirm)
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeNodeAvailabilityModalAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})

	// Add the modal to the pages
	mm.addNodeAvailabilityModalToPages(modal)

	// Set focus to the modal so it can receive keyboard input
	mm.setFocusToModal(modal)
}

// ShowRetryDialog displays a retry dialog with automatic retry logic
func (mm *ModalManager) ShowRetryDialog(
	operation string,
	err error,
	retryFunc func() error,
	onSuccess func(),
) {
	// Create retry dialog content
	content := mm.createRetryDialogContent(operation, err)

	// Create modal with retry options
	modal := mm.createModal(content, []string{"Retry", "Retry (Auto)", "Cancel"})

	// Add done function to handle button clicks
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		mm.handleRetryDialogButtonClick(buttonLabel, operation, retryFunc, onSuccess)
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeRetryDialogAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})

	// Add the modal to the pages
	mm.addRetryDialogToPages(modal)

	// Set focus to the modal so it can receive keyboard input
	mm.setFocusToModal(modal)
}

// ShowFallbackDialog displays a fallback operations dialog
func (mm *ModalManager) ShowFallbackDialog(
	operation string,
	err error,
	fallbackOptions []string,
	onFallback func(string),
) {
	content := mm.createFallbackContent(operation, err)
	buttons := mm.createFallbackButtons(fallbackOptions)
	modal := mm.createModal(content, buttons)

	mm.setupFallbackModalHandlers(modal, onFallback)
	mm.showFallbackModal(modal)
}

// createFallbackContent creates the content for the fallback dialog
func (mm *ModalManager) createFallbackContent(operation string, err error) string {
	return fmt.Sprintf(
		"⚠️  Operation Failed: %s\n\nError: %v\n\nAlternative operations are available:",
		operation,
		err,
	)
}

// setupErrorModalHandlers sets up the handlers for the error modal
func (mm *ModalManager) setupErrorModalHandlers(modal *tview.Modal) {
	// Add done function to handle OK button click
	modal.SetDoneFunc(func(_ int, _ string) {
		mm.closeErrorModalAndRestoreFocus()
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeErrorModalAndRestoreFocus()
			return nil // Consume the event
		}
		return event
	})
}

// closeErrorModalAndRestoreFocus closes the error modal and restores focus
func (mm *ModalManager) closeErrorModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("error_modal")
	mm.restoreFocusToMainView()
}

// showErrorModal shows the error modal
func (mm *ModalManager) showErrorModal(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("error_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// createReplicasInputField creates the replicas input field
func (mm *ModalManager) createReplicasInputField(currentReplicas uint64) *tview.InputField {
	return tview.NewInputField().
		SetLabel("Replicas: ").
		SetText(fmt.Sprintf("%d", currentReplicas)).
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldInteger)
}

// createScaleForm creates the scale form with buttons
func (mm *ModalManager) createScaleForm(
	inputField *tview.InputField,
	onConfirm func(int),
) *tview.Form {
	return tview.NewForm().
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
			mm.closeScaleModalAndRestoreFocus()
			onConfirm(replicas)
		}).
		AddButton("Help", func() {
			// Show contextual help for service scaling
			mm.ShowContextualHelp("swarm_services", "scale")
		}).
		AddButton("Cancel", func() {
			// Close modal without action
			mm.closeScaleModalAndRestoreFocus()
		})
}

// createScaleModal creates the scale modal
func (mm *ModalManager) createScaleModal(serviceName string, currentReplicas uint64) *tview.Modal {
	return tview.NewModal().
		SetText(fmt.Sprintf("Scale Service: %s\nCurrent Replicas: %d", serviceName, currentReplicas)).
		SetDoneFunc(func(_ int, _ string) {
			mm.closeScaleModalAndRestoreFocus()
		})
}

// createScaleFlexContainer creates the flex container for the scale modal
func (mm *ModalManager) createScaleFlexContainer(modal *tview.Modal, form *tview.Form) *tview.Flex {
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(modal, 0, 1, false).
		AddItem(form, 0, 1, true)

	// Add keyboard handling for ESC key
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeScaleModalAndRestoreFocus()
			return nil
		}
		return event
	})

	return flex
}

// createNodeAvailabilityContent creates the content for the node availability modal
func (mm *ModalManager) createNodeAvailabilityContent(nodeName, currentAvailability string) string {
	return fmt.Sprintf(
		"Update Node Availability: %s\n\nCurrent Availability: %s\n\nSelect new availability:",
		nodeName,
		currentAvailability,
	)
}

// handleNodeAvailabilityButtonClick handles button clicks in the node availability modal
func (mm *ModalManager) handleNodeAvailabilityButtonClick(
	buttonLabel string,
	onConfirm func(string),
) {
	switch buttonLabel {
	case "Active":
		onConfirm("active")
		mm.closeNodeAvailabilityModalAndRestoreFocus()
	case "Pause":
		onConfirm("pause")
		mm.closeNodeAvailabilityModalAndRestoreFocus()
	case "Drain":
		onConfirm("drain")
		mm.closeNodeAvailabilityModalAndRestoreFocus()
	case "Help":
		// Show contextual help for node availability updates
		mm.ShowContextualHelp("swarm_nodes", "update_availability")
	case "Cancel":
		// Close the modal without action
		mm.closeNodeAvailabilityModalAndRestoreFocus()
	}
}

// closeNodeAvailabilityModalAndRestoreFocus closes the node availability modal and restores focus
func (mm *ModalManager) closeNodeAvailabilityModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("availability_modal")
	mm.restoreFocusToMainView()
}

// addNodeAvailabilityModalToPages adds the node availability modal to the pages
func (mm *ModalManager) addNodeAvailabilityModalToPages(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("availability_modal", modal, true, true)
}

// createRetryDialogContent creates the content for the retry dialog
func (mm *ModalManager) createRetryDialogContent(operation string, err error) string {
	return fmt.Sprintf(
		"🔄 Operation Failed: %s\n\nError: %v\n\nThis may be a temporary issue. Would you like to retry?",
		operation,
		err,
	)
}

// handleRetryDialogButtonClick handles button clicks in the retry dialog
func (mm *ModalManager) handleRetryDialogButtonClick(
	buttonLabel, operation string,
	retryFunc func() error,
	onSuccess func(),
) {
	switch buttonLabel {
	case "Retry":
		// Manual retry - close dialog and let user retry
		mm.closeRetryDialogAndRestoreFocus()
	case "Retry (Auto)":
		// Automatic retry with progress indication
		mm.performAutomaticRetry(operation, retryFunc, onSuccess)
	case "Cancel":
		// Close dialog without retry
		mm.closeRetryDialogAndRestoreFocus()
	}
}

// closeRetryDialogAndRestoreFocus closes the retry dialog and restores focus
func (mm *ModalManager) closeRetryDialogAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("retry_modal")
	mm.restoreFocusToMainView()
}

// addRetryDialogToPages adds the retry dialog to the pages
func (mm *ModalManager) addRetryDialogToPages(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("retry_modal", modal, true, true)
}

// createFallbackButtons creates the buttons for the fallback dialog
func (mm *ModalManager) createFallbackButtons(fallbackOptions []string) []string {
	buttons := make([]string, len(fallbackOptions)+1)
	copy(buttons, fallbackOptions)
	buttons[len(fallbackOptions)] = "Cancel"
	return buttons
}

// setupFallbackModalHandlers sets up the handlers for the fallback modal
func (mm *ModalManager) setupFallbackModalHandlers(modal *tview.Modal, onFallback func(string)) {
	modal.SetDoneFunc(func(_ int, buttonLabel string) {
		mm.handleFallbackButtonClick(buttonLabel, onFallback)
	})

	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.handleFallbackEscape()
			return nil
		}
		return event
	})
}

// handleFallbackButtonClick handles button clicks in the fallback modal
func (mm *ModalManager) handleFallbackButtonClick(buttonLabel string, onFallback func(string)) {
	if buttonLabel == "Cancel" {
		mm.closeFallbackModal()
		mm.restoreFocusToMainView()
	} else {
		onFallback(buttonLabel)
		mm.closeFallbackModal()
		mm.restoreFocusToMainView()
	}
}

// handleFallbackEscape handles ESC key press in the fallback modal
func (mm *ModalManager) handleFallbackEscape() {
	mm.closeFallbackModal()
	mm.restoreFocusToMainView()
}

// closeFallbackModal closes the fallback modal
func (mm *ModalManager) closeFallbackModal() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("fallback_modal")
}

// showFallbackModal displays the fallback modal
func (mm *ModalManager) showFallbackModal(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("fallback_modal", modal, true, true)

	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
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
	helpSections := mm.getHelpSections()
	return mm.formatHelpText(helpSections)
}

// getHelpSections returns the help sections configuration
func (mm *ModalManager) getHelpSections() []struct {
	title   string
	content []string
} {
	return []struct {
		title   string
		content []string
	}{
		mm.createGlobalHelpSection(),
		mm.createNavigationHelpSection(),
		mm.createTableNavigationHelpSection(),
		mm.createContainerActionsHelpSection(),
		mm.createImageActionsHelpSection(),
		mm.createVolumeActionsHelpSection(),
		mm.createNetworkActionsHelpSection(),
		mm.createSwarmServiceActionsHelpSection(),
		mm.createSwarmNodeActionsHelpSection(),
		mm.createConfigurationHelpSection(),
	}
}

// createGlobalHelpSection creates the global help section
func (mm *ModalManager) createGlobalHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Global",
		content: []string{
			"ESC       Close modal",
			"Ctrl+C    Exit application",
			"Q         Exit application",
			"F5        Refresh",
			"?         Show help",
		},
	}
}

// createNavigationHelpSection creates the navigation help section
func (mm *ModalManager) createNavigationHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Navigation",
		content: []string{
			"1, c      Containers view",
			"2, i      Images view",
			"3, v      Volumes view",
			"4, n      Networks view",
			"s         Swarm Services view",
			"w         Swarm Nodes view",
		},
	}
}

// createTableNavigationHelpSection creates the table navigation help section
func (mm *ModalManager) createTableNavigationHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Table Navigation",
		content: []string{
			"↑/↓       Navigate rows",
			"Enter     View details & actions",
			"ESC       Close details",
		},
	}
}

// createContainerActionsHelpSection creates the container actions help section
func (mm *ModalManager) createContainerActionsHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Container Actions",
		content: []string{
			"s         Start container",
			"S         Stop container",
			"r         Restart container",
			"d         Delete container",
			"l         View logs",
			"i         Inspect container",
		},
	}
}

// createImageActionsHelpSection creates the image actions help section
func (mm *ModalManager) createImageActionsHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Image Actions",
		content: []string{
			"d         Delete image",
			"i         Inspect image",
		},
	}
}

// createVolumeActionsHelpSection creates the volume actions help section
func (mm *ModalManager) createVolumeActionsHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Volume Actions",
		content: []string{
			"d         Delete volume",
			"i         Inspect volume",
		},
	}
}

// createNetworkActionsHelpSection creates the network actions help section
func (mm *ModalManager) createNetworkActionsHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Network Actions",
		content: []string{
			"d         Delete network",
			"i         Inspect network",
		},
	}
}

// createSwarmServiceActionsHelpSection creates the swarm service actions help section
func (mm *ModalManager) createSwarmServiceActionsHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Swarm Service Actions",
		content: []string{
			"i         Inspect service",
			"s         Scale service",
			"r         Remove service",
			"l         View logs",
		},
	}
}

// createSwarmNodeActionsHelpSection creates the swarm node actions help section
func (mm *ModalManager) createSwarmNodeActionsHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Swarm Node Actions",
		content: []string{
			"i         Inspect node",
			"a         Update availability",
			"r         Remove node",
		},
	}
}

// createConfigurationHelpSection creates the configuration help section
func (mm *ModalManager) createConfigurationHelpSection() struct {
	title   string
	content []string
} {
	return struct {
		title   string
		content []string
	}{
		title: "Configuration",
		content: []string{
			":         Command mode",
			"theme     Custom themes (YAML/JSON)",
			"refresh   Auto-refresh settings",
		},
	}
}

// formatHelpText formats the help sections into a readable string
func (mm *ModalManager) formatHelpText(helpSections []struct {
	title   string
	content []string
},
) string {
	helpText := "whaletui Keyboard Shortcuts\n\n"
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
		return mm.getServiceScalingHelp()
	case "remove":
		return mm.getServiceRemovalHelp()
	case "inspect":
		return mm.getServiceInspectionHelp()
	case "logs":
		return mm.getServiceLogsHelp()
	default:
		return mm.getServiceGeneralHelp()
	}
}

// getServiceScalingHelp returns help text for service scaling
func (mm *ModalManager) getServiceScalingHelp() string {
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
}

// getServiceRemovalHelp returns help text for service removal
func (mm *ModalManager) getServiceRemovalHelp() string {
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
}

// getServiceInspectionHelp returns help text for service inspection
func (mm *ModalManager) getServiceInspectionHelp() string {
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
}

// getServiceLogsHelp returns help text for service logs
func (mm *ModalManager) getServiceLogsHelp() string {
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
}

// getServiceGeneralHelp returns general help text for swarm services
func (mm *ModalManager) getServiceGeneralHelp() string {
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

// generateSwarmNodesHelp creates help content for swarm nodes context
func (mm *ModalManager) generateSwarmNodesHelp(operation string) string {
	switch operation {
	case "update_availability":
		return mm.getNodeAvailabilityHelp()
	case "remove":
		return mm.getNodeRemovalHelp()
	case "inspect":
		return mm.getNodeInspectionHelp()
	default:
		return mm.getNodeGeneralHelp()
	}
}

// getNodeAvailabilityHelp returns help text for node availability updates
func (mm *ModalManager) getNodeAvailabilityHelp() string {
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
}

// getNodeRemovalHelp returns help text for node removal
func (mm *ModalManager) getNodeRemovalHelp() string {
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
}

// getNodeInspectionHelp returns help text for node inspection
func (mm *ModalManager) getNodeInspectionHelp() string {
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
}

// getNodeGeneralHelp returns general help text for swarm nodes
func (mm *ModalManager) getNodeGeneralHelp() string {
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
func (mm *ModalManager) performAutomaticRetry(
	operation string,
	retryFunc func() error,
	onSuccess func(),
) {
	// Close the retry dialog
	mm.closeRetryDialog()

	// Show progress modal
	progressModal := mm.createProgressModal(operation)

	// Add the progress modal to the pages
	mm.addProgressModalToPages(progressModal)

	// Perform retry in a goroutine to avoid blocking UI
	go mm.executeRetryOperation(operation, retryFunc, onSuccess)
}

// closeRetryDialog closes the retry dialog
func (mm *ModalManager) closeRetryDialog() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("retry_modal")
}

// createProgressModal creates the progress modal for automatic retry
func (mm *ModalManager) createProgressModal(operation string) *tview.Modal {
	progressContent := fmt.Sprintf(
		"🔄 Retrying: %s\n\nPlease wait while we attempt to recover...",
		operation,
	)
	progressModal := mm.createModal(progressContent, []string{"Cancel"})

	// Add cancel functionality
	progressModal.SetDoneFunc(func(_ int, _ string) {
		mm.closeProgressModalAndRestoreFocus()
	})

	return progressModal
}

// closeProgressModalAndRestoreFocus closes the progress modal and restores focus
func (mm *ModalManager) closeProgressModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("retry_progress_modal")
	mm.restoreFocusToMainView()
}

// addProgressModalToPages adds the progress modal to the pages
func (mm *ModalManager) addProgressModalToPages(progressModal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("retry_progress_modal", progressModal, true, true)
}

// executeRetryOperation executes the retry operation in a goroutine
func (mm *ModalManager) executeRetryOperation(
	operation string,
	retryFunc func() error,
	onSuccess func(),
) {
	// Attempt retry
	err := retryFunc()

	// Close progress modal from main thread
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.QueueUpdateDraw(func() {
		mm.handleRetryResult(operation, err, onSuccess)
	})
}

// handleRetryResult handles the result of the retry operation
func (mm *ModalManager) handleRetryResult(operation string, err error, onSuccess func()) {
	mm.closeProgressModal()

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
	mm.restoreFocusToMainView()
}

// closeProgressModal closes the progress modal
func (mm *ModalManager) closeProgressModal() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("retry_progress_modal")
}

// createConfirmModal creates the confirmation modal
func (mm *ModalManager) createConfirmModal(text string, callback func(bool)) *tview.Modal {
	return tview.NewModal().
		SetText(text).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, _ string) {
			mm.handleConfirmButtonClick(buttonIndex, callback)
		})
}

// setupConfirmModalHandlers sets up the keyboard handlers for the confirm modal
func (mm *ModalManager) setupConfirmModalHandlers(modal *tview.Modal, callback func(bool)) {
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.handleConfirmEscape(callback)
			return nil
		}
		return event
	})
}

// handleConfirmButtonClick handles button clicks in the confirm modal
func (mm *ModalManager) handleConfirmButtonClick(buttonIndex int, callback func(bool)) {
	mm.closeConfirmModal()
	callback(buttonIndex == 0)
	mm.restoreFocusToMainView()
}

// handleConfirmEscape handles ESC key press in the confirm modal
func (mm *ModalManager) handleConfirmEscape(callback func(bool)) {
	mm.closeConfirmModal()
	callback(false)
	mm.restoreFocusToMainView()
}

// closeConfirmModal closes the confirmation modal
func (mm *ModalManager) closeConfirmModal() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("confirm_modal")
}

// showConfirmModal displays the confirmation modal
func (mm *ModalManager) showConfirmModal(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("confirm_modal", modal, true, true)

	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// setupHelpModalHandlers configures the event handlers for the help modal
func (mm *ModalManager) setupHelpModalHandlers(modal *tview.Modal) {
	// Add done function to handle Close button click
	modal.SetDoneFunc(func(_ int, _ string) {
		mm.closeHelpModal()
	})

	// Add keyboard handling for ESC key to close modal
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			mm.closeHelpModal()
			return nil // Consume the event
		}
		return event
	})
}

// closeHelpModal closes the help modal and restores focus
func (mm *ModalManager) closeHelpModal() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("help_modal")

	// Restore focus to the main view after closing modal
	if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
		if vc, ok := viewContainer.(*tview.Flex); ok {
			app, ok := mm.ui.GetApp().(*tview.Application)
			if !ok {
				return
			}
			app.SetFocus(vc)
		}
	}
}

// addHelpModalToUI adds the help modal to the UI and sets focus
func (mm *ModalManager) addHelpModalToUI(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("help_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// closeInfoModalAndRestoreFocus closes the info modal and restores focus
func (mm *ModalManager) closeInfoModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("info_modal")
	mm.restoreFocusToMainView()
}

// showInfoModal shows the info modal and sets focus
func (mm *ModalManager) showInfoModal(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("info_modal", modal, true, true)

	// Set focus to the modal so it can receive keyboard input
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}

// restoreFocusToMainView restores focus to the main view container
func (mm *ModalManager) restoreFocusToMainView() {
	if viewContainer := mm.ui.GetViewContainer(); viewContainer != nil {
		if vc, ok := viewContainer.(*tview.Flex); ok {
			app, ok := mm.ui.GetApp().(*tview.Application)
			if !ok {
				return
			}
			app.SetFocus(vc)
		}
	}
}

// closeContextualHelpModalAndRestoreFocus closes the contextual help modal and restores focus
func (mm *ModalManager) closeContextualHelpModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("contextual_help_modal")
	mm.restoreFocusToMainView()
}

// closeScaleModalAndRestoreFocus closes the scale modal and restores focus
func (mm *ModalManager) closeScaleModalAndRestoreFocus() {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.RemovePage("scale_modal")
	mm.restoreFocusToMainView()
}

// addScaleModalToPages adds the scale modal to the pages
func (mm *ModalManager) addScaleModalToPages(flex *tview.Flex) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("scale_modal", flex, true, true)
}

// setFocusToForm sets focus to the form
func (mm *ModalManager) setFocusToForm(form *tview.Form) {
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(form)
}

// addContextualHelpModalToPages adds the contextual help modal to the pages
func (mm *ModalManager) addContextualHelpModalToPages(modal *tview.Modal) {
	pages, ok := mm.ui.GetPages().(*tview.Pages)
	if !ok {
		return
	}
	pages.AddPage("contextual_help_modal", modal, true, true)
}

// setFocusToModal sets focus to the modal
func (mm *ModalManager) setFocusToModal(modal *tview.Modal) {
	app, ok := mm.ui.GetApp().(*tview.Application)
	if !ok {
		return
	}
	app.SetFocus(modal)
}
