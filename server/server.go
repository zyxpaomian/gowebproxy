package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"gowebproxy/server/route"
	"golang.org/x/sync/errgroup"
	"k8s.io/klog"
)

// server 结构体
type Server struct {
	port int
	routingTables *route.RoutingTable
}

// New 创建一个新的服务器
func NewServer(port int) *Server {
	// 先等待路由表初始化	
	rtb := route.NewRoutingTable()
	s := &Server{
		port: port,
		routingTables:  rtb,
	}
	return s
}

func (s *Server) Run() error {
	var eg errgroup.Group
	// 启动http server 
	eg.Go(func() error {
		srv := http.Server{
			Addr:    fmt.Sprintf(":%d", s.port),
			Handler: s,
		}
		klog.Infof("[webproxy] starting http proxy server")
		err := srv.ListenAndServe()
		if err != nil {
			return fmt.Errorf("[webproxy] start http proxy server failed, error: %v", err)
		}
		return nil
	})
	return eg.Wait()
}

// 实际的代理服务功能
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取后端的真实服务地址
	backendURL, err := s.routingTables.GetBackend(r.Host, r.URL.Path)
	if err != nil {
		http.Error(w, "upstream server not found", http.StatusNotFound)
		return
	}

	klog.Infof("[webproxy] get proxy request from: %s%s to: %v", r.Host, r.URL.Path, backendURL)
	// 使用 NewSingleHostReverseProxy 进行代理请求
	p := httputil.NewSingleHostReverseProxy(backendURL)
	p.ServeHTTP(w, r)
}