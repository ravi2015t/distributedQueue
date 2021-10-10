package integration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	s := web.NewServer(dirname, port)

	log.Printf("Listening connections")
	return s.Serve()

}
