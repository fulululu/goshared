package goshared_test

import (
	"context"
	"testing"
	"time"

	"github.com/fulululu/goshared"
	"golang.org/x/sync/errgroup"
)

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
