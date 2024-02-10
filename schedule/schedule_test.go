package schedule_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/uaxe/infra/schedule"
)

func TestSchedule(t *testing.T) {
	schedule.ScheduleAtFixRate(schedule.CalculateDelay(0, 1), 1*time.Second, func(now time.Time) error {
		fmt.Printf("A-->time:%v\n", now.UTC())
		time.Sleep(1 * time.Second)
		return nil
	})

	time.Sleep(2 * time.Second)
}
