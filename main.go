package main

import (
	"fmt"
	"log"
	. "simcache/cache"
	"time"
)

func main() {
	var simCache Cache
	simCache, err := NewMyCache()
	if err != nil {
		log.Fatal("New cache error", err)
		return
	}
	simCache.Set("abcd", "Just for test!", 3*time.Second)
	time.Sleep(5*time.Second)
	fmt.Println(simCache.Keys())
}
