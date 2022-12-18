Group 是 TinyCache 最核心的数据结构，负责与用户的交互，并且控制缓存值存储和获取的流程。

在 `tinyCache.go` 中实现主体结构 Group

```go
// A Group is a cache namespace and associated data loaded spread over
// 一个 Group 可以认为是一个缓存的命名空间
type Group struct {
	// 每个 Group 拥有一个唯一的名称 name
	// 比如可以创建三个 Group，缓存学生的成绩命名为 scores，缓存学生信息的命名为 info，缓存学生课程的命名为 courses
	name string
	// 缓存未命中时获取源数据的回调(callback)
	getter Getter
	// 之前实现的并发缓存
	mainCache cache
}
```



```
                           是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |       
                |  否                        是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
```





TinyCache 的代码结构的雏形如下：

```
tinyCache/
    |--lru/
        |--lru.go  // lru 缓存淘汰策略
    |--byteview.go // 缓存值的抽象与封装
    |--cache.go    // 并发控制
    |--tinyCache.go // 负责与外部交互，控制缓存存储和获取的主流程
```