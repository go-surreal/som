package sdbd

import (
	"sync"
)

type liveQueries struct {
	store sync.Map
}

func (l *liveQueries) get(key string) chan []byte {
	val, ok := l.store.Load(key)

	if !ok {
		ch := make(chan []byte)
		l.store.Store(key, ch)
		return ch
	}

	return val.(chan []byte)
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
