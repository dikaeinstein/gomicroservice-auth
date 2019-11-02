package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/sirupsen/logrus"
)

// Health represents the health handler
type Health struct {
	logger *log.Logger
	statsd *statsd.Client
}

// NewHealth creates a new Health handler
func NewHealth(l *log.Logger, s *statsd.Client) *Health {
	return &Health{l, s}
}

// Get returns the health status of this service
func (h *Health) Get(w http.ResponseWriter, r *http.Request) {
	defer func(startTime time.Time) {
		h.statsd.Timing("health.timing", time.Since(startTime), nil, 1)
	}(time.Now())

	h.statsd.Incr("health.success", nil, 1)
	fmt.Fprintln(w, "OK")
}
