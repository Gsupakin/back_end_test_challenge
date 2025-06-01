package domain

import (
	"time"
)

type RequestLog struct {
	Method    string    `bson:"method"`
	Path      string    `bson:"path"`
	Status    int       `bson:"status"`
	LatencyMS int64     `bson:"latency_ms"`
	Timestamp time.Time `bson:"timestamp"`
}
