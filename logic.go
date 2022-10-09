package goshared

import (
	"context"
	"math"
	"reflect"
	"time"
)

// Ternary if b == true return t else return f
func Ternary(b bool, t, f interface{}) interface{} {
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

// Paginate ...
func Paginate(slice interface{}, offset *uint64, limit *uint64) (result []interface{}) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("parameter 'slice' given a non-slice type")
	}

	result = make([]interface{}, 0)
	sl := s.Len()
	of := 0
	li := math.MaxInt

	if offset != nil {
		of = int(*offset)
	}
	if limit != nil {
		li = int(*limit)
	}

	if sl == 0 || of >= sl || li == 0 { // boundary condition
		return nil
	}

	num := 0
	for i := 0; i < sl; i++ {
		if i >= of {
			result = append(result, s.Index(i).Interface())
			num++
		}
		if num >= li {
			break
		}
	}

	return result
}

// FilterSlice ...
func FilterSlice(source interface{}, condition func(element interface{}) bool) (result []interface{}) {
	s := reflect.ValueOf(source)
	if s.Kind() != reflect.Slice {
		panic("parameter 'source' is slice type, but is given a non-slice type")
	}

	result = make([]interface{}, 0)
	sl := s.Len()

	for i := 0; i < sl; i++ {
		if condition(s.Index(i).Interface()) {
			result = append(result, s.Index(i).Interface())
		}
	}

	return
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
