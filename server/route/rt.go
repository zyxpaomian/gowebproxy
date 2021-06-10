package route
import (
	"net/url"
	"fmt"
)

// 一个host可能匹配多个path
// TODO: 添加tls认证
type RoutingTable struct {
	// CertByHost *tls.Certificate
	Backends map[string][]routingTableBackend
}

// 初始化一个新的路由表
func NewRoutingTable() *RoutingTable {
	rt := &RoutingTable{
		//certificatesByHost: make(map[string]map[string]*tls.Certificate),
		Backends: make(map[string][]routingTableBackend),
	}
	// 真实的服务器host + 端口
	rtb, _ := newroutingTableBackend("hello", "127.0.0.1", 12345)
	rt.Backends["www.zyx.com"] = append(rt.Backends["www.zyx.com"], rtb)
	return rt
}

// 根据访问的host 以及 path 获取真实的backend地址
func (rt *RoutingTable) GetBackend(host, path string) (*url.URL, error) {
	backends := rt.Backends[host]
	for _, backend := range backends {
		if backend.matches(path) {
			return backend.svcUrl, nil
		}
	}
	return nil, fmt.Errorf("no backend server found")
}