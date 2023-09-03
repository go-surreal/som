//go:build embed

package query

import (
	"context"
	"fmt"
)

type Database interface {
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
	Live(ctx context.Context, statement string, vars map[string]any) (<-chan []byte, error)
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

//
// -- LIVE
//

type liveResult[T any] struct {
	res T
	err error
}

func (r *liveResult[T]) Get() (T, error) {
	return r.res, r.err
}

func live[T any](ctx context.Context, in <-chan []byte, unmarshal func(buf []byte, val any) error) <-chan *liveResult[T] {
	out := make(chan *liveResult[T], 1)

	go func() {
		defer close(out)

		for {
			select {

			case <-ctx.Done():
				return

			case data, closed := <-in:
				if closed {
					return
				}

				var res T
				var outErr error

				if err := unmarshal(data, &res); err != nil {
					outErr = fmt.Errorf("could not unmarshal live result: %w", err)
				}

				out <- &liveResult[T]{
					res: res,
					err: outErr,
				}
			}
		}
	}()

	return out
}

// questions:
// - do the live query channels block the reading of new websocket messages?
// - should identical liver queries only be registered once?
