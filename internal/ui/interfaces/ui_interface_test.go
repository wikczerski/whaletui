package interfaces

import (
	"testing"

	uimocks "github.com/wikczerski/D5r/internal/ui/interfaces/mocks"
)

// Ensure the mockery-generated mock implements UIInterface at compile time
var _ UIInterface = (*uimocks.MockUIInterface)(nil)

func TestUIInterfaceImplementation(t *testing.T) {
	// Compile-time assertion above is sufficient; this test ensures the file is exercised
}

func TestUIInterfaceMethodCount(t *testing.T) {
	// This test intentionally has no runtime checks; it's a placeholder for interface evolution notes
}

func TestMockUIFunctionality(t *testing.T) {
	mock := uimocks.NewMockUIInterface(t)

	shutdown := make(chan struct{})
	mock.On("GetShutdownChan").Return(shutdown).Once()

	if ch := mock.GetShutdownChan(); ch == nil {
		t.Error("GetShutdownChan should return a channel")
	}
}

func TestUIInterfaceCompatibility(t *testing.T) {
	// Interface variable assignment should compile
	var ui UIInterface
	_ = ui

	// Function signature compatibility should compile
	testFunction := func(ui UIInterface) {
		_ = ui
	}
	_ = testFunction
}
