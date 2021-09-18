package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ravi2015t/distributedQueue/server"
	"github.com/ravi2015t/distributedQueue/web"
)

var (
	dirname = flag.String("dirname", "", "The directory name where to put all the data")
	port    = flag.Uint("port", 8080, "Network port to listen on")
)

func main() {
	flag.Parse()

	var backend *server.OnDisk

	if *dirname == "" {
		log.Fatalf("The flag `--dirname` must be provided")
	}

	filename := filepath.Join(*dirname, "write_test")
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Could not create test file %q: %s", filename, err)
	}
	fp.Close()
	os.Remove(fp.Name())

	backend = server.NewOnDisk(*dirname)

	s := web.NewServer(backend, *port)

	log.Printf("Listening connections")
	s.Serve()
}
