//go:build embed

package query

import (
	"context"
	"fmt"
	"strings"
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

type liveResponse[M any] struct {
	Action string `json:"action"`
	Result M      `json:"result"`
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

func live[M any](ctx context.Context, in <-chan []byte, unmarshal func(buf []byte, val any) error) <-chan LiveResult[M] {
	out := make(chan LiveResult[M], 1)

	go func() {
		defer close(out)

		for {
			select {

			case <-ctx.Done():
				return

			case data, open := <-in:
				if !open {
					return
				}

				out <- toLiveResult[M](data, unmarshal)
			}
		}
	}()

	return out
}

func toLiveResult[M any](in []byte, unmarshal func(buf []byte, val any) error) LiveResult[M] {
	var out liveResult[M]

	var response liveResponse[M]

	if err := unmarshal(in, &response); err != nil {
		out.err = fmt.Errorf("could not unmarshal live result: %w", err)
	}

	out.res = response.Result

	switch strings.ToLower(response.Action) {

	case "create":
		return &liveCreate[M]{
			liveResult: out,
		}

	case "update":
		return &liveUpdate[M]{
			liveResult: out,
		}

	case "delete":
		return &liveDelete[M]{
			liveResult: out,
		}

	default:
		out.err = fmt.Errorf("unknown action type %s", response.Action)

		return &out
	}
}

type LiveResult[M any] interface {
	live()
}

type LiveAny[M any] interface {
	LiveResult[M]
	Get() (M, error)
}

type LiveCreate[M any] interface {
	LiveAny[M]
	create()
}

type LiveUpdate[M any] interface {
	LiveAny[M]
	update()
}

type LiveDelete[M any] interface {
	LiveAny[M]
	delete()
}

type liveCreate[M any] struct {
	liveResult[M]
}

func (*liveCreate[M]) create() {}

type liveUpdate[M any] struct {
	liveResult[M]
}

func (*liveUpdate[M]) update() {}

type liveDelete[M any] struct {
	liveResult[M]
}

func (*liveDelete[M]) delete() {}

type liveResult[M any] struct {
	res M
	err error
}

func (*liveResult[M]) live() {}

func (r *liveResult[M]) Get() (M, error) {
	return r.res, r.err
}

// questions:
// - do the live query channels block the reading of new websocket messages?
// - should identical liver queries only be registered once?
