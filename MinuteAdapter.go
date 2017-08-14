package crongo

import (
	"time"
	"strings"
	"strconv"
	"sync"
)

type MinuteAdapter struct {
	timeStep time.Duration
	formatSequence []uint
	ticker *time.Ticker
}
func (this *MinuteAdapter)Init()error{
	this.timeStep = 60 * time.Second
	this.formatSequence = []uint{MINUTE, HOUR, DAY, MONTH, WEEK}
	this.ticker = time.NewTicker(this.timeStep)
	return nil
}
func (this *MinuteAdapter)Run(tasks *sync.Map) {
	defer this.ticker.Stop()
	for {
		select {
		case <-(this.ticker.C):
			tasks.Range(func(id, task interface{})bool{
				_id,_ := id.(string)
				_task,_ := task.(*Task)
				if this.isHitNow(_task.format) {
					go func(){
						_task.callback(_id, _task.format, _task.extra)
					}()
				}
				return true
			})
		}
	}
}
func (this *MinuteAdapter) isHitNow(format string)bool{
	formats := strings.Split(strings.Trim(format, " "), " ")
	var expected,denominator int
	var err error
	for i,f := range formats {
		if f != "*" {
			if strings.Contains(f, "*/") {
				expected = 0
				denominator,err = strconv.Atoi(strings.Trim(f, "*/"))
			}else{
				denominator = 1<<8
				expected,err = strconv.Atoi(f)
			}
			if err != nil {
				return false
			}
			switch this.formatSequence[i] {
			case MINUTE:
				if time.Now().Minute()%denominator != expected {
					return false
				}
			case HOUR:
				if time.Now().Hour()%denominator != expected {
					return false
				}
			case DAY:
				if time.Now().Day()%denominator != expected {
					return false
				}
			case MONTH:
				if time.Now().Month()%time.Month(denominator) != time.Month(expected) {
					return false
				}
			case WEEK:
				if time.Now().Weekday()%time.Weekday(denominator) != time.Weekday(expected) {
					return false
				}
			}
		}
	}
	return true
}

