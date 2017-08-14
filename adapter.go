package crongo

import (
	"time"
	"strings"
	"strconv"
	"sync"
	"errors"
	//"fmt"
)

type Adapter struct {
	timeStep       time.Duration
	formatSequence []uint
	ticker         *time.Ticker
}

func (this *Adapter) Init() error {
	this.timeStep = time.Second
	this.formatSequence = []uint{SECOND, MINUTE, HOUR, DAY, MONTH, WEEK}
	this.ticker = time.NewTicker(this.timeStep)
	return nil
}
func (this *Adapter) Run(tasks *sync.Map) error {
	defer this.ticker.Stop()
	var err error
	defer func() {
		if p := recover(); p != nil {
			switch x := p.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	for {
		select {
		case <-(this.ticker.C):
			tasks.Range(func(id, task interface{}) bool {
				_id, _ := id.(string)
				_task, _ := task.(*Task)
				if this.isHitNow(_task.format) {
					go func() {
						_task.callback(_id, _task.format, _task.extra)
					}()
				}
				return true
			})
		}
	}
}
func (this *Adapter) matchExpected(field string, timePart int) (bool, error) {
	var err error
	if strings.Contains(field, "*/") {
		var denominator int
		denominator, err = strconv.Atoi(strings.Trim(field, "*/"))
		//fmt.Printf("*/ %#v %#v %#v \n", field, denominator, timePart)
		if err != nil || denominator <= 0 {
			return false, errors.New("error occurs at parsing format */ in " + field)
		}
		return 0 == timePart%denominator, nil
	} else if strings.Contains(strings.Trim(field, " "), ",") {
		targets := strings.Split(field, ",")
		real := strconv.Itoa(timePart)
		//fmt.Printf(", %#v %#v %#v \n", field, targets, real)
		for _, v := range targets {
			if v == real {
				return true, nil
			}
		}
	} else if strings.Contains(field, "-") {
		ranges := strings.Split(field, "-")
		low, err1 := strconv.Atoi(ranges[0])
		high, err2 := strconv.Atoi(ranges[1])
		//fmt.Printf("- %#v %#v %#v \n", field, low, high)

		if err1 != nil || err2 != nil || low > high {
			return false, errors.New("error occurs at parsing format - in " + field)
		}
		for i := low; i <= high; i++ {
			if i == timePart {
				return true, nil
			}
		}
	} else {
		//fmt.Printf("number %#v %#v  \n", field, timePart)
		return field == strconv.Itoa(timePart), nil
	}
	return false, nil
}
func (this *Adapter) timePart(i int) int {
	var part int
	switch this.formatSequence[i] {
	case SECOND:
		part = time.Now().Second()
	case MINUTE:
		part = time.Now().Minute()
	case HOUR:
		part = time.Now().Hour()
	case DAY:
		part = time.Now().Day()
	case MONTH:
		part = int(time.Now().Month())
	case WEEK:
		part = int(time.Now().Weekday())
	}
	return part
}

func (this *Adapter) isHitNow(format string) bool {
	formats := strings.Split(strings.Trim(format, " "), " ")
	for i, f := range formats {
		f = strings.Trim(f, " ")
		if f != "*" {
			m, e := this.matchExpected(f, this.timePart(i))
			//fmt.Printf("matchExpected %#v %#v \n", m, e)

			if !m {
				if e != nil {
					panic(e)
				}
				return false
			}
		}
	}
	return true
}
