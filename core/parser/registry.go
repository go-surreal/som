package parser

import "sort"

type registry[H any] struct {
	handlers []H
	sorted   bool
}

func (r *registry[H]) register(h H) {
	r.handlers = append(r.handlers, h)
	r.sorted = false
}

func (r *registry[H]) all(priority func(H) int) []H {
	if !r.sorted {
		sort.Slice(r.handlers, func(i, j int) bool {
			return priority(r.handlers[i]) < priority(r.handlers[j])
		})
		r.sorted = true
	}
	return r.handlers
}
