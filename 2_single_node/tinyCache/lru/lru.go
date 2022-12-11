package lru

import "container/list"

// Value use Len to count how many bytes it takes
// 值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// 键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射
type entry struct {
	key   string
	Value Value
}

// Cache is an LRU cache. It is not safe for concurrent access.
// 创建一个包含字典和双向链表的结构体类型 Cache，方便实现后续的增删查改操作
type Cache struct {
	maxBytes  int64                    // 允许使用的最大内存
	usedBytes int64                    // 当前已使用的内存
	ll        *list.List               // 使用 Go 语言标准库实现的双向链表list.List
	cache     map[string]*list.Element // 字典的定义是 map[string]*list.Element，键是字符串，值是双向链表中对应节点的指针
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可以为 nil
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
// 查找主要有 2 个步骤，第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 如果键对应的链表节点存在，则将对应节点移动到队尾，并返回查找到的值
	if ele, ok := c.cache[key]; ok {
		// 将链表中的节点 ele 移动到队尾（双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾）
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.Value, true
	}
	return
}

// RemoveOldest removes the oldest item
// 这里的删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // 取到队首节点
	if ele != nil {
		c.ll.Remove(ele)                                          // 从队列 ll 中删去队首节点
		kv := ele.Value.(*entry)                                  // 获取节点的值（*entry）
		delete(c.cache, kv.key)                                   // 从字典中 c.cache 删除该节点的映射关系
		c.usedBytes -= int64(len(kv.key)) + int64(kv.Value.Len()) // 更新当前所用的内存
		if c.OnEvicted != nil {                                   // 如果回调函数 OnEvicted 不为 nil，则调用回调函数
			c.OnEvicted(kv.key, kv.Value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { // 如果键存在，则更新对应节点的值，并将该节点移到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 这里是由于对key对应的value进行了修改，只有value的长度变化了
		c.usedBytes += int64(value.Len()) - int64(kv.Value.Len())
		kv.Value = value
	} else { // 不存在则是新增场景
		ele := c.ll.PushFront(&entry{key, value})           // 首先队尾添加新节点 &entry{key, value}
		c.cache[key] = ele                                  // 并字典中添加 key 和节点的映射关系
		c.usedBytes += int64(len(key)) + int64(value.Len()) // 更新 c.usedBytes
	}
	// 如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
// 获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
