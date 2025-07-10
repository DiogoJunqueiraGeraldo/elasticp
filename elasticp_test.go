package elasticp_test

import (
	"github.com/DiogoJunqueiraGeraldo/elasticp"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
)

func GetPool() *elasticp.ElasticPool {
	pool := elasticp.New(16, 8)

	pool.Start()

	return pool
}

func TestDX(t *testing.T) {
	pool := GetPool()

	var count atomic.Int64
	var wg sync.WaitGroup

	wantCount := 2_000_000

	for i := 0; i < wantCount; i++ {
		wg.Add(1)
		pool.Go(func() {
			defer wg.Done()
			count.Add(1)
		})
	}
	wg.Wait()

	assert.Equal(t, int64(wantCount), count.Load())
}
