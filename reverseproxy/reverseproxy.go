package reverseproxy

import (
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Config struct {
	TargetAddress         string
	Timeout               time.Duration
	KeepAlive             time.Duration
	MaxIdleConns          int
	IdleConnTimeout       time.Duration
	ExpectContinueTimeout time.Duration
}

func New(config Config) (http.Handler, error) {
	parsedUrl, err := url.Parse(config.TargetAddress)
	if err != nil {
		return nil, err
	}

	if !parsedUrl.IsAbs() {
		return nil, errors.Errorf("invalid target address. scheme is empty. targetAddress=%s",
			config.TargetAddress)
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   config.Timeout,
			KeepAlive: config.KeepAlive,
		}).DialContext,
		MaxIdleConns:          config.MaxIdleConns,
		IdleConnTimeout:       config.IdleConnTimeout,
		MaxIdleConnsPerHost:   config.MaxIdleConns,
		ExpectContinueTimeout: config.ExpectContinueTimeout,
	}
	rProxy := httputil.NewSingleHostReverseProxy(parsedUrl)
	rProxy.Transport = transport

	return rProxy, nil
}
