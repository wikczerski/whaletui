// Package interfaces provides manager interface definitions for WhaleTUI.
//
//nolint:revive // Package name "interfaces" is for interface-only package
package interfaces

import (
	"github.com/rivo/tview"
)

// HeaderManagerInterface defines the interface for header management
type HeaderManagerInterface interface {
	CreateHeaderSection() tview.Primitive
	UpdateDockerInfo()
	UpdateNavigation()
	UpdateActions()
}

// ModalManagerInterface defines the interface for modal management
type ModalManagerInterface interface {
	ShowError(error)
	ShowInfo(string)
	ShowConfirm(string, func(bool))
	ShowServiceScaleModal(string, uint64, func(int))
	ShowNodeAvailabilityModal(string, string, func(string))
	ShowRetryDialog(string, error, func() error, func())
	ShowFallbackDialog(string, error, []string, func(string))
}
