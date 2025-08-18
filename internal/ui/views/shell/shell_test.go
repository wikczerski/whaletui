package shell

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestViewCreation(t *testing.T) {
	view := createTestView()
	assert.NotNil(t, view)
}

func TestViewCreation_Type(t *testing.T) {
	view := createTestView()
	assert.IsType(t, &View{}, view)
}

func TestViewCreation_Components(t *testing.T) {
	view := createFullTestView()
	assert.NotNil(t, view.view)
}

func TestViewCreation_InputField(t *testing.T) {
	view := createTestView()
	assert.NotNil(t, view.inputField)
}

func TestViewCreation_OutputView(t *testing.T) {
	view := createTestView()
	assert.NotNil(t, view.outputView)
}

func TestViewCreation_CommandHistory(t *testing.T) {
	view := createTestView()
	assert.NotNil(t, view.commandHistory)
}

func TestViewCreation_MultiLineBuffer(t *testing.T) {
	view := createTestView()
	assert.NotNil(t, view.multiLineBuffer)
}

func TestViewCreation_ContainerInfo(t *testing.T) {
	view := createTestView()
	assert.Equal(t, "test-container", view.containerID)
	assert.Equal(t, "test-container", view.containerName)
}

func TestGetView(t *testing.T) {
	view := createFullTestView()
	result := view.GetView()

	assert.NotNil(t, result)
}

func TestGetView_Type(t *testing.T) {
	view := createFullTestView()
	result := view.GetView()

	assert.IsType(t, &tview.Flex{}, result)
}

func TestGetContainerID(t *testing.T) {
	view := createTestView()
	result := view.GetContainerID()

	assert.Equal(t, "test-container", result)
}

func TestGetContainerID_Empty(t *testing.T) {
	view := createTestView()
	view.containerID = ""
	result := view.GetContainerID()

	assert.Equal(t, "", result)
}

func TestGetContainerName(t *testing.T) {
	view := createTestView()
	result := view.GetContainerName()

	assert.Equal(t, "test-container", result)
}

func TestGetContainerName_Empty(t *testing.T) {
	view := createTestView()
	view.containerName = ""
	result := view.GetContainerName()

	assert.Equal(t, "", result)
}

func TestGetInputField(t *testing.T) {
	view := createTestView()
	result := view.GetInputField()

	assert.NotNil(t, result)
}

func TestGetInputField_Type(t *testing.T) {
	view := createTestView()
	result := view.GetInputField()

	assert.IsType(t, &tview.InputField{}, result)
}

// Note: Private methods like exitShell, addOutput, isMultiLineCommand, etc.
// are not directly testable from outside the package.
// They are tested indirectly through the public interface.
