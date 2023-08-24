package schedule

import (
	"fmt"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {

	//scheduleA := make(chan time.Time, 100)
	start := time.Now()
	ScheduleAtFixRate(CalculateDelay(0, 2), 5*time.Second, func(now time.Time) error {
		fmt.Printf("A-->time:%v\n", now.UTC())
		if (time.Now().Unix() - start.Unix()) > (65) {
			//scheduleA <- now
		}
		time.Sleep(1 * time.Second)
		return nil
	})

	time.Sleep(12 * time.Second)
	//scheduleB := make(chan time.Time, 100)
	//ScheduleAtFixRate(CalculateDelay(0, 2), 5*time.Second, func(now time.Time) error {
	//	fmt.Printf("B-->time:%v\n", now.UTC())
	//	scheduleB <- now
	//	return nil
	//})
	//
	//time.Sleep(5 * time.Minute)
}
