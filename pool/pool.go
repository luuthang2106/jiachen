package pool

import "github.com/panjf2000/ants/v2"

var pool *ants.Pool

func NewPool(size int) {
	if pool != nil {
		panic("pool existed")
	}
	pool, _ = ants.NewPool(size)
}

func Submit(task func()) error {
	return pool.Submit(task)
}
