package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func lb(serverPool ServerPool, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backend := serverPool.backends[0]
		if backend.Alive == false {
			http.Error(w, "all servers are down", 500)
		}
		if backend == nil {
			http.Error(w, "xyu", 500)
		}
		fmt.Println(backend.Url.Port())

		backend.ReserveProxy.ServeHTTP(w, r)

	})

}

func main() {

	var logger *log.Logger
	fmt.Println(os.Getenv("servers"))
	servers := strings.Split(os.Getenv("servers"), " ")
	fmt.Println(servers)
	var backends []*Backend
	for _, str := range servers {
		backends = append(backends, NewBackend(str))
	}

	serverPool := NewServerPool(backends...)
	go serverPool.HealthCheck()

	server := http.Server{Addr: ":8000", Handler: lb(serverPool, logger)}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}

}
