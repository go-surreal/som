//go:build embed

package internal

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrTxClosed          = errors.New("transaction is closed")
	ErrTxAlreadyActive   = errors.New("transaction is already active in this context")
	ErrTxLiveNotSupported = errors.New("live queries are not supported within transactions")
)

type transactable interface {
	Commit(ctx context.Context) error
	Cancel(ctx context.Context) error
}

type TxState struct {
	once     sync.Once
	mu       sync.Mutex
	tx       any
	done     bool
	started  bool
	beginErr error
}

func (s *TxState) EnsureTx(beginFn func() (any, error)) (any, error) {
	s.mu.Lock()
	if s.done {
		s.mu.Unlock()
		return nil, ErrTxClosed
	}
	s.mu.Unlock()

	s.once.Do(func() {
		tx, err := beginFn()
		s.mu.Lock()
		defer s.mu.Unlock()
		s.tx = tx
		s.beginErr = err
		if err == nil {
			s.started = true
		}
	})

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tx, s.beginErr
}

func (s *TxState) commit(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.done {
		return ErrTxClosed
	}
	s.done = true

	if !s.started {
		return nil
	}

	t, ok := s.tx.(transactable)
	if !ok {
		return errors.New("transaction does not implement commit")
	}
	return t.Commit(ctx)
}

func (s *TxState) cancel(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.done {
		return nil
	}
	s.done = true

	if !s.started {
		return nil
	}

	t, ok := s.tx.(transactable)
	if !ok {
		return errors.New("transaction does not implement cancel")
	}
	return t.Cancel(ctx)
}

type txKey struct{}

func TxStart(ctx context.Context) context.Context {
	if ctx.Value(txKey{}) != nil {
		panic(ErrTxAlreadyActive)
	}
	return context.WithValue(ctx, txKey{}, &TxState{})
}

func GetTxState(ctx context.Context) *TxState {
	if s, ok := ctx.Value(txKey{}).(*TxState); ok {
		return s
	}
	return nil
}

func TxActive(ctx context.Context) bool {
	return GetTxState(ctx) != nil
}

func TxCommit(ctx context.Context) error {
	s := GetTxState(ctx)
	if s == nil {
		return nil
	}
	return s.commit(ctx)
}

func TxCancel(ctx context.Context) error {
	s := GetTxState(ctx)
	if s == nil {
		return nil
	}
	return s.cancel(ctx)
}
