package route

import (
	"net/url"
	"regexp"
	"fmt"	
)

// 实际的real server和path的正则的匹配关系, 支持正则表达式进行path的匹配
type routingTableBackend struct {
	pathRe *regexp.Regexp
	svcUrl    *url.URL 
}

// 初始化一个新的rtb对象
func newroutingTableBackend(path string, serviceName string, servicePort int) (routingTableBackend, error) {
	rtb := routingTableBackend{
		svcUrl: &url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%s:%d", serviceName, servicePort),
		},
	}
	var err error
	if path != "" {
		rtb.pathRe, err = regexp.Compile(path)
	}
	return rtb, err
}

// 判断是否可以匹配正则条件
func (rtb routingTableBackend) matches(path string) bool {
	if rtb.pathRe == nil {
		return true
	}
	return rtb.pathRe.MatchString(path)
}