package schedule_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/uaxe/infra/schedule"
)

var tw = schedule.NewTimerWheel(200*time.Millisecond, 1000)

func TestHeap(t *testing.T) {
	tw.RepeatedTimer(time.Second, func(now time.Time) {
		fmt.Printf("T1:%d\n", now.Unix())
	}, nil)

	tid2 := tw.RepeatedTimer(time.Second, func(now time.Time) {
		fmt.Printf("T2:%d\n", now.Unix())
	}, func(t time.Time) {
		fmt.Printf("T2:Cancelled %d\n", t.Unix())
	})

	time.Sleep(2 * time.Second)

	tw.CancelTimer(tid2)

	time.Sleep(3 * time.Second)
}

func TestTimer(t *testing.T) {

	tida, _ := tw.AddTimer(time.Second, func(t time.Time) {
		fmt.Printf("T1:%d\n", t.Unix())
	}, nil)

	time.Sleep(100 * time.Millisecond)
	tw.UpdateTimer(tida, time.Now().Add(2*time.Second))

	tidb, _ := tw.AddTimer(time.Second, func(t time.Time) {
		fmt.Printf("T2:%d\n", t.Unix())
	}, nil)

	time.Sleep(100 * time.Millisecond)

	tw.CancelTimer(tidb)

	time.Sleep(3 * time.Second)
}

func TestTimeWheel(t *testing.T) {

	id, ch := tw.After(1 * time.Second)
	start := time.Now().Unix()
	select {
	case <-ch:
		t.Logf("2 Seconds Timeout !")
	case <-time.After(2 * time.Second):
		tw.CancelTimer(id)
		t.FailNow()
		t.Logf("testTimeWheel 2 Seconds Not Timeout ")
	}

	t.Logf("Wait : %d s", time.Now().Unix()-start)

}
