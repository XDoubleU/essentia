package wsinternal_test

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/internal/wsinternal"
)

type TestSubscriber struct {
	id       string
	output   string
	outputMu *sync.RWMutex
}

func NewTestSubscriber() *TestSubscriber {
	return &TestSubscriber{
		id:       uuid.NewString(),
		output:   "",
		outputMu: &sync.RWMutex{},
	}
}

func (sub TestSubscriber) ID() string {
	return sub.id
}

func (sub *TestSubscriber) OnSubscribeCallback() any {
	return nil
}

func (sub *TestSubscriber) OnEventCallback(event any) {
	sub.outputMu.Lock()
	defer sub.outputMu.Unlock()

	if v, ok := event.(string); ok {
		sub.output = v
	}
}

func (sub *TestSubscriber) Output() string {
	sub.outputMu.RLock()
	defer sub.outputMu.RUnlock()

	return sub.output
}

const sleep = 100 * time.Millisecond

func TestBasic(t *testing.T) {
	wp := wsinternal.NewWorkerPool(1, 10)

	tSub := NewTestSubscriber()
	wp.AddSubscriber(tSub)
	wp.Start()
	time.Sleep(sleep)
	require.True(t, wp.Active())

	event := "Hello, World!"
	wp.EnqueueEvent(event)
	time.Sleep(sleep)
	assert.Equal(t, event, tSub.Output())

	wp.RemoveSubscriber(tSub)
	time.Sleep(sleep)
	assert.Equal(t, false, wp.Active())
}

func TestMoreWorkersThanSubs(t *testing.T) {
	wp := wsinternal.NewWorkerPool(2, 10)

	tSub := NewTestSubscriber()
	wp.AddSubscriber(tSub)
	wp.Start()
	time.Sleep(sleep)
	require.True(t, wp.Active())

	event := "Hello, World!"
	wp.EnqueueEvent(event)
	time.Sleep(sleep)
	assert.Equal(t, event, tSub.Output())

	wp.RemoveSubscriber(tSub)
	time.Sleep(sleep)
	assert.Equal(t, false, wp.Active())
}

func TestAddRemoveSubscriberWhileWorkersActive(t *testing.T) {
	wp := wsinternal.NewWorkerPool(2, 10)

	tSub := NewTestSubscriber()
	wp.AddSubscriber(tSub)
	wp.Start()
	time.Sleep(sleep)
	require.True(t, wp.Active())

	tSub2 := NewTestSubscriber()
	wp.AddSubscriber(tSub2)

	event := "Hello, World!"
	wp.EnqueueEvent(event)
	time.Sleep(sleep)
	assert.Equal(t, event, tSub.Output())
	assert.Equal(t, event, tSub2.Output())

	wp.RemoveSubscriber(tSub2)
	time.Sleep(sleep)
	assert.Equal(t, true, wp.Active())

	wp.RemoveSubscriber(tSub)
	time.Sleep(sleep)
	assert.Equal(t, false, wp.Active())
}

func work(t *testing.T, wp *wsinternal.WorkerPool, nr int) {
	t.Logf("Run %d", nr)

	tSub := NewTestSubscriber()
	wp.AddSubscriber(tSub)
	wp.Start()
	time.Sleep(sleep)
	require.True(t, wp.Active())

	event := "Hello, World!"
	wp.EnqueueEvent(event)
	time.Sleep(sleep)
	assert.Equal(t, event, tSub.Output())

	wp.RemoveSubscriber(tSub)
	time.Sleep(sleep)
	assert.Equal(t, false, wp.Active())
}

func TestToggleWork(t *testing.T) {
	wp := wsinternal.NewWorkerPool(1, 10)

	work(t, wp, 1)
	work(t, wp, 2)
}
