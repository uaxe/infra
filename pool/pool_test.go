package pool_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/uaxe/infra/pool"
)

func TestGPool_Queue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	gpool := pool.NewLimitPool(ctx, 100)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:

			}
			qsize, qcap := gpool.Monitor()
			fmt.Printf("Monitor:%d/%d\n", qsize, qcap)
			time.Sleep(1 * time.Second)
		}
	}()

	wu, err := gpool.Queue(ctx, func(ctx context.Context) (i any, e error) {
		time.Sleep(5 * time.Second)
		return "a", nil
	})

	if nil != err {
		fmt.Printf("Queue:%v\n", err)
		t.FailNow()
	}

	now := time.Now()
	resp, err := wu.Get()
	if nil != err {
		fmt.Printf("Get:%v\n", err)
		t.FailNow()
	}
	fmt.Printf("WaitResp:%v\t%v\n", err, resp)
	if resp != "a" {
		t.FailNow()
	}

	cost := time.Since(now)

	if cost < 5 {
		fmt.Printf("TooFast|WaitResp:%v\t%v\n", err, resp)
		t.FailNow()
	}

	cancel()
}

func TestNewBatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gpool := pool.NewLimitPool(ctx, 100)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:

			}
			qsize, qcap := gpool.Monitor()
			fmt.Printf("Monitor:%d/%d\n", qsize, qcap)
			time.Sleep(1 * time.Second)
		}
	}()

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	batch := gpool.NewBatch()
	wus, err := batch.Queue(func(ctx context.Context) (any, error) {
		time.Sleep(5 * time.Second)
		return "a", nil
	}).Queue(
		func(ctx context.Context) (any, error) {
			time.Sleep(5 * time.Second)
			return "b", nil
		}).Queue(func(ctx context.Context) (any, error) {
		time.Sleep(5 * time.Second)
		return "c", nil
	}).Wait(timeoutCtx)

	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	for _, wu := range wus {
		resp, err := wu.Get()
		fmt.Printf("%v|%v\n", err, resp)
		if err == nil {
			fmt.Printf("Should Timeout %v|%v\n", err, resp)
			t.FailNow()
		}
	}
	cancel()
	fmt.Println("FINISH...")

	timeoutCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	wus, err = batch.Queue(func(ctx context.Context) (any, error) {
		time.Sleep(5 * time.Second)
		return "a", nil
	}).Queue(
		func(ctx context.Context) (any, error) {
			time.Sleep(5 * time.Second)
			return "b", nil
		}).Queue(func(ctx context.Context) (any, error) {
		time.Sleep(5 * time.Second)
		return "c", nil
	}).Wait(timeoutCtx)

	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	for _, wu := range wus {
		resp, err := wu.Get()
		fmt.Printf("%v|%v\n", err, resp)
		if err != nil {
			fmt.Printf("Should Not Timeout %v|%v\n", err, resp)
			t.FailNow()
		}
	}

	time.Sleep(10 * time.Second)
	cancel()
}

func TestGPool_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gpool := pool.NewLimitPool(ctx, 100)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:

			}
			qsize, qcap := gpool.Monitor()
			fmt.Printf("Monitor:%d/%d\n", qsize, qcap)
			time.Sleep(1 * time.Second)
		}
	}()

	batch := gpool.NewBatch()
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	now := time.Now()
	wus, err := batch.Queue(func(ctx context.Context) (any, error) {
		return "a", nil
	}).Queue(
		func(ctx context.Context) (any, error) {
			time.Sleep(5 * time.Second)
			return "b", nil
		}).Queue(func(ctx context.Context) (any, error) {
		time.Sleep(5 * time.Second)
		return "c", nil
	}).Wait(timeoutCtx)

	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	cost := time.Since(now)
	fmt.Printf("Get Responses COST: %v\n", cost)
	if cost/time.Second > 2 {
		fmt.Printf("Get Responses Should Less Than 2s : %v\n", cost/time.Second)
		t.FailNow()
	}

	resps := make([]string, 0, 2)
	for _, wu := range wus {
		resp, err := wu.Get()
		fmt.Printf("%v|%v\n", err, resp)
		if err != nil && errors.Is(err, pool.ErrQueueContextDone) {
			fmt.Printf("Should Not Timeout %v|%v\n", err, resp)
			t.FailNow()
		}

		if nil != resp {
			resps = append(resps, resp.(string))
		}
	}

	sort.Strings(resps)
	if len(resps) != 1 {
		fmt.Printf("Responses Should Be Only 1 %v\n", resps)
		t.FailNow()
	}

	idx := sort.SearchStrings(resps, "a")
	if idx == len(resps) || resps[idx] != "a" {
		fmt.Printf("Responses Should Contains a %v\n", resps)
		t.FailNow()
	}

	cancel()
}
