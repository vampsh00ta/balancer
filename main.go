package main

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct {
	Url          *url.URL
	Alive        bool
	mux          sync.Mutex
	ReserveProxy *httputil.ReverseProxy
}

type ServerPool struct {
	backends []*Backend
	curr_idx uint64
}

func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.curr_idx, uint64(1)) % uint64(len(s.backends)))
}
