package lib

import (
	"time"
)

const (
	MinimumExpiry time.Duration = 1 * time.Minute  // 1 minute
	MaximumExpiry time.Duration = 1440 * time.Hour // 60 days
)
