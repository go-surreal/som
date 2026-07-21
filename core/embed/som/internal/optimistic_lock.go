//go:build embed

package internal

type OptimisticLock struct {
	version int
}

func (o *OptimisticLock) Version() int {
	return o.version
}

func (o *OptimisticLock) SetVersion(v int) {
	o.version = v
}
