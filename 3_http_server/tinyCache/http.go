package tinyCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_tinycache/"

// HTTPPool implements PeerPicker for a pool of HTTP peers.
// 首先创建一个结构体 HTTPPool，作为承载节点间 HTTP 通信的核心数据结构（包括服务端和客户端，这里先只实现服务端）
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	// 那么 http://example.com/_tinycache/ 开头的请求，就用于节点间的访问。
	// 因为一个主机上还可能承载其他的服务，加一段 Path 是一个好习惯。比如，大部分网站的 API 接口，一般以 /api 作为前缀。
	self     string // 用来记录自己的地址，包括主机名/IP 和端口
	basePath string // 作为节点间通讯地址的前缀，默认是 /_tinycache/

}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP handle all http requests
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 首先判断访问路径的前缀是否是 basePath
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	// 约定访问路径格式为 /<basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	// 通过 groupname 得到 group 实例
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	// 再使用 group.Get(key) 获取缓存数据
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/otcet-stream")
	// 最终使用 w.Write() 将缓存值作为 httpResponse 的 body 返回
	w.Write(view.ByteSlice())
}
