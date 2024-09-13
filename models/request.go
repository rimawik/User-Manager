package models

import "time"

// Struct to hold request info
type RequestLog struct {
	Method     string
	Path       string
	StartTime  time.Time
	FinishTime time.Time
	Duration   time.Duration
}
