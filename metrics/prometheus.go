package metrics

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

var responseDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "service_response_duration_seconds",
	Help: "Internal services duration",
}, []string{"status_code"})

var openConnections = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "circuit_breaker_state",
	Help: "If circuit breaker is open or notٔ",
})

var allCollector = []prometheus.Collector{
	responseDuration,
	openConnections,
}

type Config struct {
	Enabled       bool
	ListenAddress string
}

type Server struct {
	config   Config
	listener net.Listener
	handler  http.Handler
}

func New(config Config) (*Server, error) {
	ins := &Server{
		config: config,
	}
	logrus.Debugf("[prometheus] Going to listen on %s", ins.config.ListenAddress)
	newListener, err := net.Listen("tcp", ins.config.ListenAddress)
	if err != nil {
		return nil, errors.Wrap(err, "fail to listen")
	}

	if config.Enabled {
		for _, collector := range allCollector {
			if err := prometheus.Register(collector); err != nil {
				if err.Error() != (prometheus.AlreadyRegisteredError{}).Error() {
					return nil, errors.Wrapf(err, "fail to register metric %v", collector)
				}
			}
		}
	} else {
		logrus.Warn("prometheus is disabled by the config file")
	}

	ins.listener = newListener
	ins.handler = promhttp.Handler()

	return ins, nil
}

func (m *Server) Serve() {
	if !m.config.Enabled {
		return
	}

	err := http.Serve(m.listener, m.handler)
	if err != nil {
		logrus.WithError(err).Debugf("[prometheus] close previous http server that listen on: %s",
			m.config.ListenAddress)
	}

	return
}
