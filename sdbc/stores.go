package sdbc

import (
	"bytes"
	"github.com/google/uuid"
	"sync"
)

//
// -- BUFFERS
//

type bufPool struct {
	sync.Pool
}

// Get returns a buffer from the pool or
// creates a new one if the pool is empty.
func (p *bufPool) Get() *bytes.Buffer {
	b := p.Pool.Get()

	if b == nil {
		return &bytes.Buffer{}
	}

	return b.(*bytes.Buffer)
}

// Put returns a buffer into the pool.
func (p *bufPool) Put(b *bytes.Buffer) {
	b.Reset()

	p.Pool.Put(b)
}

//
// -- REQUESTS
//

type requests struct {
	store sync.Map
}

func (r *requests) prepare() (string, <-chan []byte) {
	key := uuid.New()
	ch := make(chan []byte)

	r.store.Store(key.String(), ch)

	return key.String(), ch
}

func (r *requests) get(key string) (chan<- []byte, bool) {
	val, ok := r.store.Load(key)
	if !ok {
		return nil, false
	}

	return val.(chan []byte), true
}

func (r *requests) cleanup(key string) {
	if ch, ok := r.store.LoadAndDelete(key); ok {
		close(ch.(chan []byte))
	}
}

func (r *requests) reset() {
	r.store.Range(func(key, ch any) bool {
		close(ch.(chan []byte))
		r.store.Delete(key)
		return true
	})
}

//
// -- LIVE QUERIES
//

type liveQueries struct {
	store sync.Map
}

func (l *liveQueries) get(key string, create bool) (chan []byte, bool) {
	val, ok := l.store.Load(key)

	if !ok && !create {
		return nil, false
	}

	if !ok {
		ch := make(chan []byte)
		l.store.Store(key, ch)
		return ch, true
	}

	return val.(chan []byte), true
}

func (l *liveQueries) del(key string) {
	if ch, ok := l.store.LoadAndDelete(key); ok {
		close(ch.(chan []byte))
	}
}

func (l *liveQueries) reset() {
	l.store.Range(func(key, ch any) bool {
		close(ch.(chan []byte))
		l.store.Delete(key)
		return true
	})
}
