package formatters

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeFormatter(t *testing.T) {
	// Test that FormatTime function works correctly
	now := time.Now()
	result := FormatTime(now)
	assert.NotEmpty(t, result)
}

func TestTimeFormatter_FormatDuration(t *testing.T) {
	duration := 2*time.Hour + 30*time.Minute + 45*time.Second

	result := formatDuration(duration)
	assert.Equal(t, "2h 30m", result)
}

func TestTimeFormatter_FormatDuration_Zero(t *testing.T) {
	duration := time.Duration(0)

	result := formatDuration(duration)
	assert.Equal(t, "0s", result)
}

func TestTimeFormatter_FormatDuration_Short(t *testing.T) {
	duration := 45 * time.Second

	result := formatDuration(duration)
	assert.Equal(t, "45s", result)
}

func TestTimeFormatter_FormatDuration_Long(t *testing.T) {
	duration := 25*time.Hour + 30*time.Minute

	result := formatDuration(duration)
	assert.Equal(t, "1d 1h", result)
}
