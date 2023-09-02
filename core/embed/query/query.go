//go:build embed

package query

import (
	"context"
)

type Database interface {
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
}

type idNode struct {
	ID string
}

type countResult struct {
	Count int
}

type queryResult[M any] struct {
	Result []M    `json:"result"`
	Status string `json:"status"`
	Time   string `json:"time"`
}

//
// -- ASYNC
//

type asyncResult[T any] struct {
	res <-chan T
	err <-chan error
}

func (r *asyncResult[T]) Val() <-chan T {
	return r.res
}

func (r *asyncResult[T]) Err() <-chan error {
	return r.err
}

func async[T any](ctx context.Context, fn func(ctx context.Context) (T, error)) *asyncResult[T] {
	resCh, errCh := make(chan T, 1), make(chan error, 1)

	go func() {
		defer close(resCh)
		defer close(errCh)

		res, err := fn(ctx)

		resCh <- res
		errCh <- err
	}()

	return &asyncResult[T]{
		res: resCh,
		err: errCh,
	}
}
