package integration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ravi2015t/distributedQueue/web"
	"go.etcd.io/etcd/client"
)

func InitAndServe(etcdAddr string, dirname string, port uint) error {
	cfg := client.Config{
		Endpoints:               strings.Split(etcdAddr, ","),
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("creating etcd client: %v", err)
	}
	kapi := client.NewKeysAPI(c)

	filename := filepath.Join(dirname, "write_test")
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("could not create test file %q: %s", filename, err)
	}
	fp.Close()
	os.Remove(fp.Name())

	s := web.NewServer(kapi, dirname, port)

	log.Printf("Listening connections")
	return s.Serve()

}
