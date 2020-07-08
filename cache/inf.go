package cache

import "time"

type Cache interface {
	//size 是⼀一个字符串串。⽀支持以下参数: 1KB， 100KB， 1MB， 2MB， 1GB 等
	SetMaxMemory(size string) bool
	// 设置⼀一个缓存项，并且在expire时间之后过期
	Set(key string, val interface{}, expire time.Duration)
	// 获取⼀一个值
	Get(key string) (interface{}, bool)
	// 删除⼀一个值
	Del(key string) bool
	// 检测⼀一个值 是否存在
	Exists(key string) bool
	// 情况所有值
	Flush() bool
	// 返回所有的key 多少
	Keys() int64
}
