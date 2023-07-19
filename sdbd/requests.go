package sdbd

import (
	"github.com/google/uuid"
	"sync"
)

type requests struct {
	store sync.Map
}

func (r *requests) prepare() (string, chan any) {
	key := uuid.New()
	ch := make(chan any)

	r.store.Store(key.String(), ch)

	return key.String(), ch
}

func (r *requests) get(key string) (chan<- any, bool) {
	val, ok := r.store.Load(key)
	if !ok {
		return nil, false
	}

	return val.(chan any), true
}

func (r *requests) cleanup(key string) {
	if ch, ok := r.store.LoadAndDelete(key); ok {
		close(ch.(chan any))
	}
}
