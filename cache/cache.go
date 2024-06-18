package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache *cache.Cache

func init() {
	// キャッシュの初期化
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}