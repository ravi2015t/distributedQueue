package main

import (
	"log"

	"github.com/ravi2015t/distributedQueue/server"
	"github.com/ravi2015t/distributedQueue/web"
)

func main() {
	s := web.NewServer(&server.InMemory{})

	log.Printf("Listening connections")
	s.Serve()
}
