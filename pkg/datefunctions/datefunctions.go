package datefunctions

import (
	"time"
)

type DateTool struct {
}

func NewDateTool() *DateTool {
	return &DateTool{}
}

func (u *DateTool) Now() time.Time {
	return time.Now()
}
