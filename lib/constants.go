package lib

import (
	"time"
)

const (
	EXPIRE_MIN time.Duration = 1 * time.Minute  // 1 minute
	EXPIRE_MAX time.Duration = 1440 * time.Hour // 60 days
)
