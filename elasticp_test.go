package elasticp_test

import (
	"runtime"
	"sync"
	"testing"

	"github.com/DiogoJunqueiraGeraldo/elasticp"
)

var pool *elasticp.Pool

func init() {
	pool = elasticp.New()
	pool.Grow(runtime.NumCPU())
}

func sumSliceRange(a, b []float64) {
	for i, v := range a {
		b[i] += v
	}
}

func BenchmarkSequential_100kTasks(b *testing.B) {
	const size = 1_000_000
	const tasks = 100_000
	chunk := size / tasks

	for i := 0; i < b.N; i++ {
		mem := make([]float64, size*2)
		input := mem[:size]
		output := mem[size:]

		for j := 0; j < size; j++ {
			input[j] = float64(j)
		}

		for t := 0; t < tasks; t++ {
			start := t * chunk
			end := start + chunk

			sumSliceRange(input[start:end], output[start:end])
		}
	}
}

func BenchmarkRawGoroutines_100kTasks(b *testing.B) {
	const size = 1_000_000
	const tasks = 100_000
	chunk := size / tasks

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(tasks)

		mem := make([]float64, size*2)
		input := mem[:size]
		output := mem[size:]

		for j := 0; j < size; j++ {
			input[j] = float64(j)
		}

		for t := 0; t < tasks; t++ {
			start := t * chunk
			end := start + chunk

			go func() {
				defer wg.Done()
				sumSliceRange(input[start:end], output[start:end])
			}()
		}

		wg.Wait()
	}
}

func BenchmarkElasticpPool100kTasks(b *testing.B) {
	const size = 1_000_000
	const tasks = 100_000
	chunk := size / tasks

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(tasks)

		mem := make([]float64, size*2)
		input := mem[:size]
		output := mem[size:]

		for j := 0; j < size; j++ {
			input[j] = float64(j)
		}

		for t := 0; t < tasks; t++ {
			start := t * chunk
			end := start + chunk

			inSlice := input[start:end]
			outSlice := output[start:end]

			pool.Submit(elasticp.Task{
				Input:  inSlice,
				Output: outSlice,
				Wg:     &wg,
			})
		}

		wg.Wait()
	}
}
