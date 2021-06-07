package resolve

import (
	"fmt"
	"sync"

	"github.com/buger/jsonparser"

	"github.com/jensneuse/graphql-go-tools/pkg/fastbuffer"
)

type DataLoader struct {
	totalNum           int
	inputs             [][]byte
	inputPosToBufPairs map[int]*BufPair
	fetch              *BatchFetch

	batch *batch

	mu *sync.Mutex
}

type batch struct {
	err  error
	done chan struct{}
}

func (d *DataLoader) Load(ctx *Context, input []byte, bufPair *BufPair) (err error) {
	d.mu.Lock()
	currentPosition := len(d.inputs)
	d.inputPosToBufPairs[currentPosition] = bufPair
	d.inputs = append(d.inputs, input)

	if d.batch == nil {
		d.batch = &batch{done: make(chan struct{})}
	}

	d.mu.Unlock()

	if d.isBatchReady(currentPosition) {
		go d.resolveFetch(ctx)
	}

	fmt.Println("before waiting", d.totalNum)
	select {
	case <-d.batch.done:
		err = d.batch.err
	case <- ctx.Context.Done():
		err = ctx.Context.Err()
	}
	fmt.Println("after waiting")

	return
}

func (d *DataLoader) resolveFetch(ctx *Context) {
	batchInput, err := d.fetch.PrepareBatch(d.inputs...)
	defer func() { close(d.batch.done) }()

	if err != nil {
		d.batch.err = err
		return
	}

	fmt.Println("batch request", string(batchInput.Input))

	fmt.Println()

	if ctx.beforeFetchHook != nil {
		ctx.beforeFetchHook.OnBeforeFetch(d.hookCtx(ctx), batchInput.Input)
	}

	batchBufferPair := &BufPair{
		Data:   fastbuffer.New(),
		Errors: fastbuffer.New(),
	}

	if err = d.fetch.Fetch.DataSource.Load(ctx.Context, batchInput.Input, batchBufferPair); err != nil {
		d.batch.err = err
		return
	}

	if ctx.afterFetchHook != nil {
		if batchBufferPair.HasData() {
			ctx.afterFetchHook.OnData(d.hookCtx(ctx), batchBufferPair.Data.Bytes(), false)
		}
		if batchBufferPair.HasErrors() {
			ctx.afterFetchHook.OnError(d.hookCtx(ctx), batchBufferPair.Errors.Bytes(), false)
		}
	}

	var outPosition int

	_, d.batch.err = jsonparser.ArrayEach(batchBufferPair.Data.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		inputPositions := batchInput.OutToInPositions[outPosition]

		for _, pos := range inputPositions {
			bufPair := d.inputPosToBufPairs[pos]
			bufPair.Errors.WriteBytes(batchBufferPair.Errors.Bytes())
			bufPair.Data.WriteBytes(value)
		}

		outPosition++
	})

	return
}

func (d *DataLoader) isBatchReady(currentPosition int) bool {
	return currentPosition == (d.totalNum - 1)
}

func (d *DataLoader) hookCtx(ctx *Context) HookContext {
	return HookContext{
		CurrentPath: ctx.path(),
	}
}

func NewDataLoader(fetch *BatchFetch, totalNum int) *DataLoader {
	return &DataLoader{
		totalNum:           totalNum,
		inputs:             make([][]byte, 0, totalNum),
		inputPosToBufPairs: make(map[int]*BufPair, totalNum),
		fetch:              fetch,
		mu:                 &sync.Mutex{},
	}
}
