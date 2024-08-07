//go:build embed

package query

import (
	"context"
	"fmt"
	"github.com/fxamacker/cbor/v2"
	"strings"
)

type Database interface {
	Query(ctx context.Context, statement string, vars map[string]any) ([]byte, error)
	Live(ctx context.Context, statement string, vars map[string]any) (<-chan []byte, error)
	Unmarshal(buf []byte, val any) error
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

type liveResponse struct {
	Action string          `json:"action"`
	Result cbor.RawMessage `json:"result"`
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

func live[C, M any](
	ctx context.Context,
	in <-chan []byte,
	unmarshal func(buf []byte, val any) error,
	convert func(C) M,
) <-chan LiveResult[M] {
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

				out <- toLiveResult[C, M](data, unmarshal, convert)
			}
		}
	}()

	return out
}

func toLiveResult[C, M any](
	in []byte,
	unmarshal func(buf []byte, val any) error,
	convert func(C) M,
) LiveResult[M] {
	var response liveResponse

	if err := unmarshal(in, &response); err != nil {
		return &liveResult[M]{
			err: fmt.Errorf("could not unmarshal live response: %w", err),
		}
	}

	switch strings.ToLower(response.Action) {

	case "create":
		var out liveResult[M]

		var result C

		if err := unmarshal(response.Result, &result); err != nil {
			out.err = fmt.Errorf("could not unmarshal live create result: %w", err)
		}

		if out.err == nil {
			out.res = convert(result)
		}

		return &liveCreate[M]{
			liveResult: out,
		}

	case "update":
		var out liveResult[M]

		var result C

		if err := unmarshal(response.Result, &result); err != nil {
			out.err = fmt.Errorf("could not unmarshal live update result: %w", err)
		}

		if out.err == nil {
			out.res = convert(result)
		}

		return &liveUpdate[M]{
			liveResult: out,
		}

	case "delete":
		var out liveResult[M]

		var result C

		if err := unmarshal(response.Result, &result); err != nil {
			out.err = fmt.Errorf("could not unmarshal live delete result: %w", err)
		}

		if out.err == nil {
			out.res = convert(result)
		}

		return &liveDelete[M]{
			liveResult: out,
		}

	default:
		return &liveResult[M]{
			err: fmt.Errorf("unknown action type %s", response.Action),
		}
	}
}

type LiveResult[M any] interface {
	live()
}

type LiveCreate[M any] interface {
	Get() (M, error)
	create()
}

type LiveUpdate[M any] interface {
	Get() (M, error)
	update()
}

type LiveDelete[M any] interface {
	Get() (M, error)
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
