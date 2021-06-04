package resolve

import (
	"sync"

	"github.com/buger/jsonparser"
)

type DataLoader struct {
	totalNum           int
	inputs             [][]byte
	inputPosToBufPairs map[int]*BufPair
	fetch              *BatchFetch
	resolver           *Resolver

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

	if currentPosition == (d.totalNum - 1) { // The last element of the batch
		go d.resolveFetch(ctx)
	}

	<-d.batch.done
	err = d.batch.err

	return
}

func (d *DataLoader) resolveFetch(ctx *Context) {
	batchInput, err := d.fetch.MergeInputs(d.inputs...)
	defer func() { close(d.batch.done) }()

	if err != nil {
		d.batch.err = err
		return
	}

	if ctx.beforeFetchHook != nil {
		ctx.beforeFetchHook.OnBeforeFetch(d.hookCtx(ctx), batchInput.Input)
	}

	batchPair := d.resolver.getBufPair()
	defer d.resolver.freeBufPair(batchPair)

	if err = d.fetch.fetch.DataSource.Load(ctx.Context, batchInput.Input, batchPair); err != nil {
		d.batch.err = err
		return
	}

	if ctx.afterFetchHook != nil {
		if batchPair.HasData() {
			ctx.afterFetchHook.OnData(d.hookCtx(ctx), batchPair.Data.Bytes(), false)
		}
		if batchPair.HasErrors() {
			ctx.afterFetchHook.OnError(d.hookCtx(ctx), batchPair.Errors.Bytes(), false)
		}
	}

	var outPosition int

	_, d.batch.err = jsonparser.ArrayEach(batchPair.Data.Bytes(), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		inputPositions := batchInput.OutToInPositions[outPosition]

		for i := range inputPositions {
			bufPair := d.inputPosToBufPairs[i]
			bufPair.Errors.WriteBytes(batchPair.Errors.Bytes())
			bufPair.Data.WriteBytes(value)
		}

		outPosition++
	})

	return
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
