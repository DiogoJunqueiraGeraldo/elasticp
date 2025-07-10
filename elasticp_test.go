package elasticp_test

import (
	"fmt"
	"github.com/DiogoJunqueiraGeraldo/elasticp"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestDX(t *testing.T) {
	pool := elasticp.New(elasticp.NewConfig(
		elasticp.WithDebug(false),
	))
	pool.Start()

	count := 0
	m := sync.Mutex{}
	wantCount := 1_000_000

	fmt.Println("Before Goroutines Count:", runtime.NumGoroutine())
	for i := 0; i < wantCount; i++ {
		pool.Go(func() {
			m.Lock()
			count++
			m.Unlock()
		})
	}

	fmt.Println("Immediately After Goroutines Count:", runtime.NumGoroutine())
	time.Sleep(2 * time.Second)
	fmt.Println("Two Sec After Goroutines Count:", runtime.NumGoroutine())

	assert.Equal(t, wantCount, count)
}
