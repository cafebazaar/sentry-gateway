package main

import (
	"github.com/cafebazaar/sentry-gateway/metrics"
	"github.com/cafebazaar/sentry-gateway/reverseproxy"
	"github.com/cafebazaar/sentry-gateway/throttle"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	config, err := LoadConfig()
	panicIfErr(err)

	prometheusServer, err := metrics.New(config.MetricsConfig)
	panicIfErr(err)
	go prometheusServer.Serve()

	proxy, err := reverseproxy.New(config.ReverseProxyConfig)
	panicIfErr(err)

	router := mux.NewRouter()
	chain := alice.New(
		metrics.NewMetricsMiddleware)

	proxyWithInstrumentationHandler := chain.Then(proxy)

	proxyWithInstrumentationAndThrottleHandler := chain.Append(func(next http.Handler) http.Handler {
		return throttle.New(next, config.ThrottleConfig)
	}).Then(proxy)

	router.Handle("/api/{id:[0-9]+}/store/", proxyWithInstrumentationAndThrottleHandler)
	router.PathPrefix("/").Handler(proxyWithInstrumentationHandler)

	logrus.Infof("Start listening on %s", config.ListenAddress)
	err = http.ListenAndServe(config.ListenAddress, router)
	if err != nil {
		logrus.Fatal("Web server (HTTP): ", err)
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
