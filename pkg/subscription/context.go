package subscription

import (
	"context"
	"sync"
	"net/http"
)

type InitialHttpRequestContext struct {
	context.Context
	Request *http.Request
}

func NewInitialHttpRequestContext(r *http.Request) *InitialHttpRequestContext {
	return &InitialHttpRequestContext{
		Context: r.Context(),
		Request: r,
	}
}

func newSubscriptionCancellations() subscriptionCancellations {
	return subscriptionCancellations{
		cancelFuncs: make(map[string]context.CancelFunc),
		mux:         &sync.Mutex{},
	}
}

type subscriptionCancellations struct {
	cancelFuncs map[string]context.CancelFunc
	mux *sync.Mutex
}

func (sc subscriptionCancellations) AddWithParent(id string, parent context.Context) context.Context {
	sc.mux.Lock()
	defer sc.mux.Unlock()

	ctx, cancelFunc := context.WithCancel(parent)
	sc.cancelFuncs[id] = cancelFunc
	return ctx
}

func (sc subscriptionCancellations) Cancel(id string) (ok bool) {
	sc.mux.Lock()
	defer sc.mux.Unlock()

	cancelFunc, ok := sc.cancelFuncs[id]
	if !ok {
		return false
	}

	cancelFunc()
	delete(sc.cancelFuncs, id)
	return true
}

func (sc subscriptionCancellations) CancelAll() {
	sc.mux.Lock()
	defer sc.mux.Unlock()

	for _, cancelFunc := range sc.cancelFuncs {
		cancelFunc()
	}
}

func (sc subscriptionCancellations) Count() int {
	sc.mux.Lock()
	defer sc.mux.Unlock()

	return len(sc.cancelFuncs)
}
