package elasticp_test

import (
	"context"
	"github.com/DiogoJunqueiraGeraldo/elasticp"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestElasticTest(t *testing.T) {
	ep := elasticp.New[int, int]()

	ep.Launch(1024)
	ep.Grow(1024)

	// it will take a while to shut down the goroutines
	// may spice cpu usage due to heavy comparisons (proportional to pool size)
	ep.Shrink(1024)

	gotSize := ep.Size()
	wantSize := 1024

	assert.Equal(t, wantSize, gotSize)
}

type IncrOperation struct{}

func (op *IncrOperation) Execute(_ context.Context, input *int, output *int) {
	*output = *input + 1
}

func TestTestSubmitWorkUnit(t *testing.T) {
	ep := elasticp.New[int, int]()
	ep.Launch(10)

	inp := 10
	out := 0

	wg := sync.WaitGroup{}
	op := IncrOperation{}
	wg.Add(1)
	wu := elasticp.NewWorkUnit[int, int](context.Background(), &op, &inp, &out, &wg)
	ep.Submit(wu)
	wg.Wait()

	assert.Equal(t, out, inp+1)
}
