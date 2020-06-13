package metrics

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type metricsMiddleware struct {
	next http.Handler
}

// NewEntryPointMiddleware creates a new metrics middleware for an Entrypoint.
func NewMetricsMiddleware(next http.Handler) http.Handler {
	return &metricsMiddleware{
		next: next,
	}
}

func (m *metricsMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	openConnections.Add(1)
	defer openConnections.Add(-1)

	recorder := newResponseRecorder(rw)
	start := time.Now()

	m.next.ServeHTTP(recorder, req)

	logrus.Debugf("[%s] %s%s %d %f",req.Method, req.Host, req.URL.String(), recorder.getCode(), time.Since(start).Seconds())
	responseDuration.With(map[string]string{
		"status_code": strconv.Itoa(recorder.getCode()),
	}).Observe(time.Since(start).Seconds())
}
