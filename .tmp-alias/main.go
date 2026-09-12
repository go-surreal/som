package main

import "fmt"

type builder[M any, R any] struct {
	name string
}

func (b builder[M, R]) Where(s string) builder[M, R] {
	b.name += s
	return b
}

func (b builder[M, R]) All() []R {
	return nil
}

func (b builder[M, R]) First() R {
	var zero R
	return zero
}

// Builder keeps the single-parameter form every node query uses.
type Builder[M any] = builder[M, *M]

type Drawing struct{ Name string }

type Shape interface{ Area() float64 }

func main() {
	var nodes Builder[Drawing]
	var _ []*Drawing = nodes.Where("x").All()

	var union builder[Shape, Shape]
	var _ []Shape = union.Where("x").All()
	var _ Shape = union.First()

	fmt.Println("ok")
}
