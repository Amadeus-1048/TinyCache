package main

import (
	"TinyCache/3_http_server/tinyCache"
	"fmt"
	"log"
	"net/http"
)

// 用一个 map 模拟耗时的数据库
var db = map[string]string{
	"Amadeus": "600",
	"Beta":    "500",
	"Cell":    "400",
}

func main() {
	// 创建 group 实例，并测试 Get 方法
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
	log.Fatal(http.ListenAndServe(addr, peers))
}
