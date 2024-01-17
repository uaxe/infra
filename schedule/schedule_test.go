package schedule_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/uaxe/infra/schedule"
)

func TestSchedule(t *testing.T) {
	start := time.Now()
	schedule.ScheduleAtFixRate(schedule.CalculateDelay(0, 2), 5*time.Second, func(now time.Time) error {
		fmt.Printf("A-->time:%v\n", now.UTC())
		if (time.Now().Unix() - start.Unix()) > (65) {
			//scheduleA <- now
		}
		time.Sleep(1 * time.Second)
		return nil
	})

	time.Sleep(12 * time.Second)
}
