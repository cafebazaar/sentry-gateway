package throttle

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net/http"
)

type Config struct {
	RPS   int
	Burst int
}

type Throttle struct {
	next     http.Handler
	config   Config
	limiters map[string]*rate.Limiter
}

func New(next http.Handler, config Config) http.Handler {
	return &Throttle{next, config, map[string]*rate.Limiter{}}
}

func (t *Throttle) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	logrus.Debug("Throttle Middleware")

	projectId := mux.Vars(r)["id"]
	logrus.Printf("Got %s, project: %s", r.URL, projectId)
	limiter, ok := t.limiters[projectId]
	if !ok {
		logrus.Printf("New limiter for project %s with rps %s and burst %s",
			projectId, t.config.RPS, t.config.Burst)
		limiter = rate.NewLimiter(rate.Limit(t.config.RPS), t.config.Burst)
		t.limiters[projectId] = limiter
	}
	if !limiter.Allow() {
		logrus.Infof("Too many requests! Project: %s", projectId)
		rw.WriteHeader(http.StatusTooManyRequests)
		return
	}
	t.next.ServeHTTP(rw, r)
}
