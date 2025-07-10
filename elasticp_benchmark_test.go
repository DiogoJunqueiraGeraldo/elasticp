package elasticp_test

import (
	"github.com/DiogoJunqueiraGeraldo/elasticp"
	"runtime"
	"sync"
	"testing"
)

var (
	pool *elasticp.ElasticPool
	once sync.Once
)

func getPool() *elasticp.ElasticPool {
	once.Do(func() {
		pool = elasticp.New(runtime.NumCPU(), 512)
		pool.Start()
	})
	return pool
}

func benchmarkElasticPool(n int) {
	pool := getPool()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		pool.Go(func() {
			wg.Done()
			x := 0
			for i := 0; i < 10_000; i++ {
				x += i * i * i % 17
			}
		})
	}

	wg.Wait()
}

func benchmarkGoroutines(n int) {

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			wg.Done()

			x := 0
			for i := 0; i < 10_000; i++ {
				x += i * i * i % 17
			}
		}()
	}

	wg.Wait()
}

func BenchmarkElasticPool100k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkElasticPool(100_000)
	}
}

func BenchmarkElasticPool500k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkElasticPool(500_000)
	}
}

func BenchmarkElasticPool1kk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkElasticPool(1_000_000)
	}
}

func BenchmarkElasticPool5kk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkElasticPool(5_000_000)
	}
}

func BenchmarkGoroutines100k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkGoroutines(100_000)
	}
}

func BenchmarkGoroutines500k(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkGoroutines(500_000)
	}
}

func BenchmarkGoroutines1kk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkGoroutines(1_000_000)
	}
}

func BenchmarkGoroutines5kk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkGoroutines(5_000_000)
	}
}
