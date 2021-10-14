package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ravi2015t/distributedQueue/server/replication"
	"github.com/ravi2015t/distributedQueue/web"
	"go.etcd.io/etcd/clientv3"
)

func InitAndServe(etcdAddr string, instanceName string, dirname string, listenAddr string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("creating new etcd client: %v", err)
	}
	defer etcdClient.Close()

	_, err = etcdClient.Put(ctx, "test", "test")
	if err != nil {
		return fmt.Errorf("couldn't set test key to etcd: %v", err)
	}

	_, err = etcdClient.Put(ctx, "peers/"+instanceName, listenAddr)
	if err != nil {
		return fmt.Errorf("couldn't register peers name in etcd: %v", err)
	}

	filename := filepath.Join(dirname, "write_test")
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("could not create test file %q: %s", filename, err)
	}
	fp.Close()
	os.Remove(fp.Name())

	s := web.NewServer(etcdClient, dirname, instanceName, listenAddr, replication.NewStorage(etcdClient, instanceName))

	log.Printf("Listening connections")
	return s.Serve()

}
