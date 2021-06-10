package main

import (
	"math/rand"
	"k8s.io/klog"	
	"time"
	"flag"
	"golang.org/x/sync/errgroup"
	"runtime"
	"gowebproxy/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	// 启动参数
	var port int
	var tlsPort int
	flag.IntVar(&port, "port", 80, "http server port.")
	flag.IntVar(&tlsPort, "tls-port", 443, "https server port")
	flag.Parse()

	// http proxy server
	s := server.NewServer(80)
	// 多协程启动
	var eg errgroup.Group
	eg.Go(func() error {
		return s.Run()
	})
	if err := eg.Wait(); err != nil {
		klog.Fatalf("[webproxy] something is wrong: %v", err.Error())
	}
}