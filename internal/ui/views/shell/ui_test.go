package shell

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestViewUI_Creation(t *testing.T) {
	view := createFullTestView()
	assert.NotNil(t, view)
}

func TestViewUI_Creation_Type(t *testing.T) {
	view := createFullTestView()
	assert.IsType(t, &View{}, view)
}

func TestViewUI_Structure(t *testing.T) {
	view := createFullTestView()
	assert.NotNil(t, view.view)
}

func TestViewUI_Structure_Components(t *testing.T) {
	view := createFullTestView()
	assert.NotNil(t, view.inputField)
	assert.NotNil(t, view.outputView)
}

func TestViewUI_Layout(t *testing.T) {
	view := createFullTestView()
	result := view.GetView()
	assert.IsType(t, &tview.Flex{}, result)
}

func TestViewUI_Layout_Direction(t *testing.T) {
	view := createFullTestView()
	result := view.GetView()
	flex, ok := result.(*tview.Flex)
	assert.True(t, ok)
	assert.Greater(t, flex.GetItemCount(), 0)
}

func TestViewUI_Layout_Items(t *testing.T) {
	view := createFullTestView()
	result := view.GetView()
	flex, ok := result.(*tview.Flex)
	assert.True(t, ok)
	assert.Greater(t, flex.GetItemCount(), 0)
}

func TestViewUI_Initialization(t *testing.T) {
	view := createTestView()
	assert.NotNil(t, view)
	assert.NotNil(t, view.inputField)
	assert.NotNil(t, view.outputView)
}

func TestViewUI_Initialization_State(t *testing.T) {
	view := createTestView()
	assert.Equal(t, "test-container", view.containerID)
	assert.Equal(t, "test-container", view.containerName)
	assert.Equal(t, 0, len(view.commandHistory))
	assert.Equal(t, 0, len(view.multiLineBuffer))
}

func TestViewUI_Functionality_Basic(t *testing.T) {
	view := createTestView()
	assert.NotNil(t, view.GetInputField())
	assert.NotNil(t, view.GetContainerID())
	assert.NotNil(t, view.GetContainerName())
}

func TestViewUI_EdgeCases_Empty(t *testing.T) {
	view := createTestView()
	view.containerID = ""
	view.containerName = ""
	assert.Equal(t, "", view.GetContainerID())
	assert.Equal(t, "", view.GetContainerName())
}
