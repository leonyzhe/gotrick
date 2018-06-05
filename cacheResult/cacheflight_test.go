package cacheflight

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	ast := assert.New(t)
	var lock sync.Mutex
	var keyCount = make(map[string]int)

	g := NewGroup(time.Second / 2)

	kv := map[string]string{
		"七牛": "七牛云",
		"上海": "上海市",
		"蓝":  "蓝色",
		"小黑": "是只猫",
		"":   "empty",
	}

	var wg sync.WaitGroup
	for key := range kv {
		for j := 0; j < rand.Intn(100)+1; j++ {
			k := key
			wg.Add(1)
			go func() {
				result, err := g.Do(k, func() (interface{}, error) {
					lock.Lock()
					keyCount[k]++
					lock.Unlock()
					return kv[k], nil
				})
				ast.Nil(err)
				ast.Equal(kv[k], result.(string))
				wg.Done()
			}()

		}
	}
	wg.Wait()

	for k := range kv {
		ast.Equal(1, keyCount[k], fmt.Sprint(k))
	}

	time.Sleep(time.Second)

	for key := range kv {
		for j := 0; j < rand.Intn(100)+1; j++ {
			k := key
			wg.Add(1)
			go func() {
				result, err := g.Do(k, func() (interface{}, error) {
					lock.Lock()
					keyCount[k]++
					lock.Unlock()
					return kv[k], nil
				})
				ast.Nil(err)
				ast.Equal(kv[k], result.(string))
				wg.Done()
			}()

		}
	}
	wg.Wait()

	for k := range kv {
		ast.Equal(2, keyCount[k], fmt.Sprint(k))
	}
}
