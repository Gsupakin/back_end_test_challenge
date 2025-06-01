package domain

import (
	"time"
)

type RequestLog struct {
	Method    string    `bson:"method"`
	Path      string    `bson:"path"`
	Status    int       `bson:"status"`
	LatencyMS int64     `bson:"latency_ms"`
	IP        string    `bson:"ip"`
	UserAgent string    `bson:"user_agent"`
	Timestamp time.Time `bson:"timestamp"`
}
