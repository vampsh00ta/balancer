package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	Url          *url.URL
	Alive        bool
	mux          sync.RWMutex
	PingTime     time.Duration
	ReserveProxy *httputil.ReverseProxy
}

func NewBackend(host string) *Backend {
	backend := &Backend{Url: &url.URL{Host: host, Scheme: "http"}, Alive: true}
	backend.ReserveProxy = httputil.NewSingleHostReverseProxy(backend.Url)
	backend.ReserveProxy.Director = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = host
		req.Host = host
	}
	return backend
}
