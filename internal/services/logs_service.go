package services

// LogsService implements logs operations
type logsService struct{}

// NewLogsService creates a new logs service
func NewLogsService() LogsService {
	return &logsService{}
}

// GetActions returns the available actions for logs as a map
func (s *logsService) GetActions() map[rune]string {
	return map[rune]string{
		'f': "Follow logs",
		't': "Tail logs",
		's': "Save logs",
		'c': "Clear logs",
		'w': "Wrap text",
	}
}

// GetActionsString returns the available actions for logs as a formatted string
func (s *logsService) GetActionsString() string {
	return "<f> Follow logs\n<t> Tail logs\n<s> Save logs\n<c> Clear logs\n<w> Wrap text"
}
