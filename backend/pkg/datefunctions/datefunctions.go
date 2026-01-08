package datefunctions

import (
	"time"
)

// DateTool struct for working with time
type DateTool struct {
}

// NewDateTool constructor for DateTool
func NewDateTool() *DateTool {
	return &DateTool{}
}

// Now getting the current time
func (u *DateTool) Now() time.Time {
	return time.Now()
}
