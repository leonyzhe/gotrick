package cacheflight

import (
	"sync"
	"time"

	"github.com/golang/groupcache/singleflight"
)

type cacheResult struct {
	data  interface{}
	ctime time.Time
	// 保证在并发访问时 key 对应的 fn 只执行一次
	doing chan bool
}

type Group struct {
	cacheExpiration time.Duration
	sfg             *singleflight.Group
	cache           map[string]cacheResult
	lock            sync.RWMutex
}

func NewGroup(cacheExpiration time.Duration) (group *Group) {

	// do something
	group = &Group{
		sfg:             &singleflight.Group{},
		cache:           make(map[string]cacheResult),
		cacheExpiration: cacheExpiration,
	}
	return
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (ret interface{}, err error) {

	// do something
	return g.do(key, fn)
}

func (g *Group) do(key string, fn func() (interface{}, error)) (ret interface{}, err error) {

	return g.sfg.Do(key, func() (interface{}, error) {

		g.lock.Lock()
		result, ok := g.cache[key]
		if result.doing != nil {
			g.lock.Unlock()
			// 等待 fn 执行完成
			<-result.doing
			return g.do(key, fn)
		} else if !ok || result.ctime.Add(g.cacheExpiration).Before(time.Now()) {
			var doing = make(chan bool)
			result.doing = doing
			g.cache[key] = result
			// fn 执行完成后关闭 channel，使其他等待的 goroutine 能及时响应
			defer func() {
				close(doing)
			}()
		} else {
			g.lock.Unlock()
			return result.data, nil
		}
		g.lock.Unlock()

		// 不同 key 的 fn 能做到并发执行
		ret, err = fn()
		if err != nil {
			return ret, nil
		}

		result = cacheResult{
			data:  ret,
			ctime: time.Now(),
		}
		g.lock.Lock()
		g.cache[key] = result
		g.lock.Unlock()

		return ret, nil
	})
}
