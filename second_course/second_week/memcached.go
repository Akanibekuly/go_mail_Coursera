package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	MemcachedAddresses := []string{"127.0.0.1:11211"}
	memcacheClient := memcache.New(MemcachedAddresses...)

	mKey := "coursera"

	memcacheClient.Set(&memcache.Item{
		Key:        mKey,
		Value:      []byte("1"),
		Expiration: 3,
	})

	memcacheClient.Increment("habrTag", 1)

	item, err := memcacheClient.Get(mKey)
	if err != nil && err != memcache.ErrCacheMiss {
		fmt.Println("MC error:", err)
	}

	fmt.Printf("mc value %#v\n", item)

	memcacheClient.Delete(mKey)

	item, err = memcacheClient.Get(mKey)
	if err != nil && err == memcache.ErrCacheMiss {
		fmt.Println("record not found ")
	}
}
