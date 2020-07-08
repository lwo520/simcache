package cache

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DefaultMaxMem string = "128000KB"
)

type vItem struct {
	Value      interface{}
	Expiration int64
}

func (item vItem) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().Unix() > item.Expiration
}

type mycache struct {
	items     map[string]vItem
	maxMemory string
	lock      sync.Mutex
}

func NewMyCache() (*mycache, error) {
	mc := &mycache{
		items:     make(map[string]vItem),
		maxMemory: DefaultMaxMem,
	}
	return mc, nil
}

func (c *mycache) CheckMaxMemeory() bool {
	sInt, err := strconv.ParseUint(
		strings.TrimSuffix(c.maxMemory, "KB"), 10, 64)
	if err != nil {
		log.Fatal(err)
		return false
	}
	memStat := GetMemStatus()
	if memStat.Self < sInt {
		return true
	}
	return false
}

func (c *mycache) SetMaxMemory(size string) bool {
	c.lock.Lock()
	if size == "" {
		return false
	}
	c.maxMemory = DefaultMaxMem
	var (
		sInt uint64
		err  error
	)
	defer c.lock.Unlock()
	if strings.HasSuffix(size, "KB") {
		sInt, err = strconv.ParseUint(
			strings.TrimSuffix(size, "KB"), 10, 64)
	} else if strings.HasSuffix(size, "MB") {
		sInt, err = strconv.ParseUint(
			strings.TrimSuffix(size, "MB"), 10, 64)
		if err == nil {
			sInt = sInt * 1000
		}
	} else if strings.HasSuffix(size, "GB") {
		sInt, err = strconv.ParseUint(
			strings.TrimSuffix(size, "GB"), 10, 64)
		if err == nil {
			sInt = sInt * 1000 * 1000
		}
	}
	if err != nil {
		log.Fatal("您设置的最大内存值不正确！")
		c.maxMemory = DefaultMaxMem
		return false
	}
	memStat := GetMemStatus()
	if sInt > memStat.All {
		log.Fatal("您设置的最大内存超过了系统的最大内存范围！")
		c.maxMemory = DefaultMaxMem
		return false
	} else {
		size = strconv.FormatUint(sInt, 10)
		c.maxMemory = size + "KB"
	}
	return true
}

func (c *mycache) Set(key string, val interface{}, expired time.Duration) {
	if !c.CheckMaxMemeory() {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, existed := c.items[key]; existed {
		if !c.items[key].Expired() {
			log.Fatalf("key: %v alerdy existed.", key)
			return
		}
	}

	c.items[key] = vItem{
		Value:      val,
		Expiration: time.Now().Unix() + int64(expired.Seconds()),
	}
}

func (c *mycache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if val, existed := c.items[key]; existed {
		if c.items[key].Expired() {
			return "Expired.", false
		}
		return val, true
	}
	return "Doesn't existed.", false
}

func (c *mycache) Del(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, existed := c.items[key]; existed {
		delete(c.items, key)
		return true
	}
	return false
}

func (c *mycache) Exists(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, existed := c.items[key]; existed {
		return true
	}
	return false
}

func (c *mycache) Flush() bool {
	c.lock.Lock()
	for key, _ := range c.items {
		delete(c.items, key)
	}
	c.lock.Unlock()
	return true
}

func (c *mycache) Keys() int64 {
	c.lock.Lock()
	defer c.lock.Unlock()
	var keysCount int64 = 0
	for key, vit := range c.items {
		if !vit.Expired() {
			keysCount += 1
		} else {
			delete(c.items, key)
		}
	}
	return keysCount
}
