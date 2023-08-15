package schedule

import (
	"time"
)

//计算延迟的时间
func CalculateDelay(minute int, second int) time.Duration {
	now := time.Now()
	delay := 0
	if minute > 0 && now.Minute()%minute != 0 {
		delay = ((now.Minute()/minute+1)*minute - now.Minute()) * 60
	}
	if second < 0 {
		second = 0
	}
	delay += (60 - now.Second()) + second

	return time.Duration(int64(delay) * int64(time.Second))
}

//按照固定的时间执行
func ScheduleAtFixRate(delay time.Duration, period time.Duration, callback func(now time.Time) error) {

	go func() {
		time.Sleep(delay)
		initT := time.Now()
		for {
			now := time.Now()
			//执行命令
			callback(now)
			timeRange := int64(now.Sub(initT)) / int64(period)
			if int64(now.Sub(initT))%int64(period) != 0 {
				timeRange += 1
			}
			//下一次执行的时间
			next := initT.Add(time.Duration(timeRange) * period)

			now = time.Now()
			if next.Before(now) {
				time.Sleep(now.Sub(next))
			} else {
				time.Sleep(next.Sub(now))
			}
		}
	}()
}