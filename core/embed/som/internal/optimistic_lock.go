//go:build embed

package internal

type OptimisticLock struct {
	version int
}

func NewOptimisticLock(version int) OptimisticLock {
	return OptimisticLock{
		version: version,
	}
}

func (o *OptimisticLock) Version() int {
	return o.version
}

func (o *OptimisticLock) SetVersion(v int) {
	o.version = v
}
