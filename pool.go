package main

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type ServerPool struct {
	backends []*Backend
	checker  []*Backend
	curr_idx uint64
}

func NewServerPool(backends ...*Backend) ServerPool {
	serverPool := ServerPool{backends: backends}
	serverPool.checker = serverPool.backends
	return serverPool
}
func (pool *ServerPool) makeHealthRequest(backend *Backend) error {
	req, err := http.NewRequest(http.MethodGet, backend.Url.Scheme+"://"+backend.Url.Host, nil)
	if err != nil {

		return err
	}

	c := &http.Client{Timeout: 5 * time.Second}
	_, err = c.Do(req)
	if err != nil {
		backend.PingTime = time.Second * 1000
		backend.Alive = false
		return err
	}
	return nil
}
func (pool *ServerPool) HealthCheck() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			wg := &sync.WaitGroup{}
			for _, backend := range pool.checker {
				wg.Add(1)
				go func(backend *Backend) {
					defer wg.Done()
					start := time.Now().UTC()
					err := pool.makeHealthRequest(backend)
					if err != nil {
						fmt.Println(err)
						backend.Alive = false
						return
					}
					duration := time.Since(start)
					backend.PingTime = duration
					backend.Alive = true
				}(backend)
			}
			wg.Wait()

			sort.Slice(pool.checker, func(i, j int) bool {
				return pool.checker[i].PingTime < pool.checker[j].PingTime
			})
			pool.backends = pool.checker

		}
	}

}
