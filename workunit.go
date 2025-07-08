package elasticp

import (
	"context"
	"sync"
)

type Operation[I, O any] interface {
	Execute(context.Context, *I, *O)
}

type WorkUnit[I, O any] struct {
	ctx       context.Context
	operation Operation[I, O]
	input     *I
	output    *O
	wg        *sync.WaitGroup
}

func NewWorkUnit[I, O any](ctx context.Context, op Operation[I, O], input *I, output *O, wg *sync.WaitGroup) WorkUnit[I, O] {
	return WorkUnit[I, O]{
		ctx:       ctx,
		operation: op,
		input:     input,
		output:    output,
		wg:        wg,
	}
}

func (wu *WorkUnit[I, O]) Execute() {
	defer wu.wg.Done()
	wu.operation.Execute(wu.ctx, wu.input, wu.output)
}
