package main

// 为什么这里是绝对路径，而在之前写的TinyGin中是相对路径呢
// main.go 和 tinyCache/ 在同级目录，但 go modules 不再支持 import <相对路径>，相对路径需要在 go.mod 中声明：
// require tinyCache v0.0.0
// replace tinyCache => ./tinyCache
import (
	"TinyCache/3_http_server/tinyCache"
	"fmt"
	"log"
	"net/http"
)

// 用一个 map 模拟数据源 db
var db = map[string]string{
	"Amadeus": "600",
	"Beta":    "500",
	"Cell":    "400",
}

func main() {
	// 创建一个名为 scores 的 Group，若缓存为空，回调函数会从 db 中获取数据并返回
	tinyCache.NewGroup("scores", 2<<10, tinyCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := tinyCache.NewHTTPPool(addr)
	log.Println("tinyCache is running at", addr)
	// 使用 http.ListenAndServe 在 9999 端口启动了 HTTP 服务
	log.Fatal(http.ListenAndServe(addr, peers))
}

/*
测试方法：
curl http://localhost:9999/_tinycache/scores/Amadeus
curl http://localhost:9999/_tinycache/scores/Beta
curl http://localhost:9999/_tinycache/scores/Cell
*/
