package main

import (
	"flag"
	"log"

	"github.com/ravi2015t/distributedQueue/integration"
)

var (
	dirname      = flag.String("dirname", "", "The directory name where to put all the data")
	instanceName = flag.String("instance-name", "", "The unique instance name where to put all the data")
	listenAddr   = flag.String("listen", "127.0.0.1:8080", "Network addr to listen on")
	etcdAddr     = flag.String("etcd", "127.0.0.1:2739", "The network address of etcd server(s)")
)

func main() {
	flag.Parse()

	if *dirname == "" {
		log.Fatalf("The flag `--dirname` must be provided")
	}

	if *instanceName == "" {
		log.Fatalf("The flag `--instance-name` must be provided")
	}

	if *etcdAddr == "" {
		log.Fatalf("The flag `--etcd` must be provided")
	}

	if err := integration.InitAndServe(*etcdAddr, *instanceName, *dirname, *listenAddr); err != nil {
		log.Fatalf("Init and Serve failed %v", err)
	}
}
