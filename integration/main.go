package integration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ravi2015t/distributedQueue/server"
	"github.com/ravi2015t/distributedQueue/web"
)

func InitAndServe(dirname string, port uint) error {
	filename := filepath.Join(dirname, "write_test")
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("could not create test file %q: %s", filename, err)
	}
	fp.Close()
	os.Remove(fp.Name())

	backend, err := server.NewOnDisk(dirname)
	if err != nil {
		return fmt.Errorf("could not initialise on-disk backend: %v", err)
	}

	s := web.NewServer(backend, port)

	log.Printf("Listening connections")
	return s.Serve()

}
