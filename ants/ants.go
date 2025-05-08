package ants

import "github.com/panjf2000/ants/v2"

func NewPool(size int) (*ants.Pool, error) {
	p, err := ants.NewPool(size)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func NewPoolWithFunc(size int, f func(any), o ...ants.Option) (*ants.PoolWithFunc, error) {
	p, err := ants.NewPoolWithFunc(size, f, o...)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func NewMultiPool(size, sizePerPool int, lbs ants.LoadBalancingStrategy) (*ants.MultiPool, error) {
	mp, err := ants.NewMultiPool(size, sizePerPool, lbs)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

func NewMultiPoolWithFunc(size, sizePerPool int, f func(any), lbs ants.LoadBalancingStrategy, o ...ants.Option) (*ants.MultiPoolWithFunc, error) {
	mp, err := ants.NewMultiPoolWithFunc(size, sizePerPool, f, lbs, o...)
	if err != nil {
		return nil, err
	}

	return mp, nil
}
