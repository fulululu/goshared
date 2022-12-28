package goshared

import (
	"context"
	"math"
	"time"
)

// Ternary if b == true return t else return f
func Ternary[T comparable](b bool, t, f T) T {
	if b {
		return t
	}
	return f
}

// RepeatedlyDo do some operation at least once
// @Param op represent operation function which has 'func() error' signature
// @Param rt represent repeated times
func RepeatedlyDo(op func() error, rt uint) error {
	var count uint = 0
	var err error
	for err = op(); err != nil && count < rt; count++ {
		err = op()
		if err == nil {
			return nil
		}
	}
	return err
}

// SlicePaginate ...
func SlicePaginate[T comparable](in []T, offset *uint64, limit *uint64) (out []T) {
	sliceLength := uint64(len(in))
	defaultOffset := uint64(0)
	defaultLimit := uint64(math.MaxInt)

	if offset != nil {
		defaultOffset = *offset
	}
	if limit != nil {
		defaultLimit = *limit
	}

	if sliceLength == 0 || defaultOffset >= sliceLength || defaultLimit == 0 { // boundary condition
		return nil
	}

	if defaultOffset+defaultLimit > sliceLength {
		out = append(out, in[defaultOffset:]...)
	} else {
		out = append(out, in[defaultOffset:defaultOffset+defaultLimit]...)
	}

	return out
}

// SliceFilter ...
func SliceFilter[T comparable](in []T, condition func(element T) bool) (out []T) {
	for _, v := range in {
		if condition(v) {
			var tmp = v
			out = append(out, tmp)
		}
	}

	return out
}

// GetFirstElementFromSliceByCondition ...
func GetFirstElementFromSliceByCondition[T comparable](in []T, condition func(element T) bool) (out T, ok bool) {
	for _, v := range in {
		if condition(v) {
			var tmp = v
			return tmp, true
		}
	}

	return out, false
}

// SliceContains ...
func SliceContains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

type LimitedOperation func(context.Context) error
type FrequencyLimiter func(ctx context.Context, op LimitedOperation) (done bool, err error)

// MakeFrequencyLimiter ...
func MakeFrequencyLimiter(period time.Duration, frequency uint64) FrequencyLimiter {
	start := time.Now()
	count := uint64(0)
	return func(ctx context.Context, op LimitedOperation) (bool, error) {
		now := time.Now()
		if now.Add(-period).After(start) {
			start = now
			count = 0
		}

		if count < frequency {
			count++
			return true, op(ctx)
		} else {
			return false, nil
		}
	}
}
