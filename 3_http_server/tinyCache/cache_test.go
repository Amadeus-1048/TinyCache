package tinyCache

import (
	"fmt"
	"log"
	"testing"
)

// 用一个 map 模拟耗时的数据库
var db = map[string]string{
	"Amadeus": "600",
	"Beta":    "500",
	"Cell":    "400",
}

func Test_Get(t *testing.T) {
	// 创建 group 实例，并测试 Get 方法
	loadCounts := make(map[string]int, len(db))
	tinyCache := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	// 在这个测试用例中，我们主要测试了 2 种情况
	for k, v := range db {
		// 在缓存为空的情况下，能够通过回调函数获取到源数据
		if view, err := tinyCache.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Amadeus")
		} // load from callback function
		// 在缓存已经存在的情况下，是否直接从缓存中获取
		// 为了实现这一点，使用 loadCounts 统计某个键调用回调函数的次数，如果次数大于1，则表示调用了多次回调函数，没有缓存。
		if _, err := tinyCache.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := tinyCache.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be empty, but got %s", view)
	}
}
