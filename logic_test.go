package goshared_test

import (
	"context"
	"testing"
	"time"

	"github.com/fulululu/goshared"
	"golang.org/x/sync/errgroup"
)

func TestSlicePaginate(t *testing.T) {
	type st struct {
		ID uint64
	}
	var before = []st{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
	}
	var after []st

	var offset, limit = uint64(1), uint64(2)
	after = goshared.SlicePaginate(before, &offset, &limit)
	if len(after) != 2 && after[0].ID != uint64(2) && after[1].ID != uint64(3) {
		t.Errorf("incorrect result: %v", after)
	}
	t.Logf("correct result: %v", after)

	offset, limit = uint64(3), uint64(5)
	after = goshared.SlicePaginate(before, &offset, &limit)
	if len(after) != 2 && after[0].ID != uint64(4) && after[1].ID != uint64(5) {
		t.Errorf("incorrect result: %v", after)
	}
	t.Logf("correct result: %v", after)

	offset, limit = uint64(3), uint64(0)
	after = goshared.SlicePaginate(before, &offset, &limit)
	if len(after) != 0 {
		t.Errorf("incorrect result: %v", after)
	}
	t.Logf("correct result: %v", after)

	offset, limit = uint64(6), uint64(3)
	after = goshared.SlicePaginate(before, &offset, &limit)
	if len(after) != 0 {
		t.Errorf("incorrect result: %v", after)
	}
	t.Logf("correct result: %v", after)
}

func TestSliceFilter(t *testing.T) {
	type st struct {
		ID uint64
	}
	var before = []st{
		{ID: 1},
		{ID: 2},
		{ID: 3},
		{ID: 4},
		{ID: 5},
	}
	var after []st

	condition1 := func(element st) bool {
		if element.ID > 3 {
			return true
		} else {
			return false
		}
	}
	after = goshared.SliceFilter(before, condition1)
	if len(after) != 2 && after[0].ID != uint64(4) && after[1].ID != uint64(5) {
		t.Errorf("incorrect result: %v", after)
	}
	t.Logf("correct result: %v", after)
}

func TestMakeFrequencyLimiter(t *testing.T) {
	op := func(ctx context.Context) error {
		time.Sleep(time.Second)
		return nil
	}

	freq1, freq2 := uint64(3), uint64(2)
	count1, count2 := uint64(0), uint64(0)
	limiter1Do := goshared.MakeFrequencyLimiter(5*time.Second, freq1)
	limiter2Do := goshared.MakeFrequencyLimiter(5*time.Second, freq2)

	timer := time.NewTimer(10*time.Second - 100*time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if (count1>>1) != freq1 || (count2>>1) != freq2 {
				t.Errorf("freq1: %d, count1: %d; freq2: %d, count2: %d", freq1, count1, freq2, count2)
			}
			return
		default:
			eg := errgroup.Group{}
			eg.Go(func() error {
				if done, err := limiter1Do(context.Background(), op); err != nil {
					t.Errorf("limiter1 do: %v", err)
				} else if done {
					count1++
					t.Logf("%v limiter1 done", time.Now().Format(time.StampMilli))
				}
				return nil
			})
			eg.Go(func() error {
				if done, err := limiter2Do(context.Background(), op); err != nil {
					t.Errorf("limiter2 do: %v", err)
				} else if done {
					count2++
					t.Logf("%v limiter2 done", time.Now().Format(time.StampMilli))
				}
				return nil
			})
			eg.Wait()
		}
	}
}
