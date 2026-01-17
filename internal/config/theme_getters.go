package config

import (
	"github.com/gdamore/tcell/v2"
)

// GetHeaderColor returns the header color
func (tm *ThemeManager) GetHeaderColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Header)
}

// GetBorderColor returns the border color
func (tm *ThemeManager) GetBorderColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Border)
}

// GetTextColor returns the text color
func (tm *ThemeManager) GetTextColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Text)
}

// GetBackgroundColor returns the background color
func (tm *ThemeManager) GetBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Background)
}

// GetSuccessColor returns the success color
func (tm *ThemeManager) GetSuccessColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Success)
}

// GetWarningColor returns the warning color
func (tm *ThemeManager) GetWarningColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Warning)
}

// GetErrorColor returns the error color
func (tm *ThemeManager) GetErrorColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Error)
}

// GetInfoColor returns the info color
func (tm *ThemeManager) GetInfoColor() tcell.Color {
	return tm.GetColor(tm.config.Colors.Info)
}

// GetShellBorderColor returns the shell border color
func (tm *ThemeManager) GetShellBorderColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Border)
}

// GetShellTitleColor returns the shell title color
func (tm *ThemeManager) GetShellTitleColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Title)
}

// GetShellTextColor returns the shell text color
func (tm *ThemeManager) GetShellTextColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Text)
}

// GetShellBackgroundColor returns the shell background color
func (tm *ThemeManager) GetShellBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Background)
}

// GetShellCmdLabelColor returns the shell command label color
func (tm *ThemeManager) GetShellCmdLabelColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Label)
}

// GetShellCmdBorderColor returns the shell command border color
func (tm *ThemeManager) GetShellCmdBorderColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Border)
}

// GetShellCmdTextColor returns the shell command text color
func (tm *ThemeManager) GetShellCmdTextColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Text)
}

// GetShellCmdBackgroundColor returns the shell command background color
func (tm *ThemeManager) GetShellCmdBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Background)
}

// GetShellCmdPlaceholderColor returns the shell command placeholder color
func (tm *ThemeManager) GetShellCmdPlaceholderColor() tcell.Color {
	return tm.GetColor(tm.config.Shell.Cmd.Placeholder)
}

// GetContainerExecLabelColor returns the container exec label color
func (tm *ThemeManager) GetContainerExecLabelColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Label)
}

// GetContainerExecBorderColor returns the container exec border color
func (tm *ThemeManager) GetContainerExecBorderColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Border)
}

// GetContainerExecTextColor returns the container exec text color
func (tm *ThemeManager) GetContainerExecTextColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Text)
}

// GetContainerExecBackgroundColor returns the container exec background color
func (tm *ThemeManager) GetContainerExecBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Background)
}

// GetContainerExecPlaceholderColor returns the container exec placeholder color
func (tm *ThemeManager) GetContainerExecPlaceholderColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Placeholder)
}

// GetContainerExecTitleColor returns the container exec title color
func (tm *ThemeManager) GetContainerExecTitleColor() tcell.Color {
	return tm.GetColor(tm.config.ContainerExec.Title)
}

// GetCommandModeLabelColor returns the command mode label color
func (tm *ThemeManager) GetCommandModeLabelColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Label)
}

// GetCommandModeBorderColor returns the command mode border color
func (tm *ThemeManager) GetCommandModeBorderColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Border)
}

// GetCommandModeTextColor returns the command mode text color
func (tm *ThemeManager) GetCommandModeTextColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Text)
}

// GetCommandModeBackgroundColor returns the command mode background color
func (tm *ThemeManager) GetCommandModeBackgroundColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Background)
}

// GetCommandModePlaceholderColor returns the command mode placeholder color
func (tm *ThemeManager) GetCommandModePlaceholderColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Placeholder)
}

// GetCommandModeTitleColor returns the command mode title color
func (tm *ThemeManager) GetCommandModeTitleColor() tcell.Color {
	return tm.GetColor(tm.config.CommandMode.Title)
}
