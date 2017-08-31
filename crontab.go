package crongo

import "sync"

func init() {

}

// six columns mean：
//       second：0-59
//       minute：0-59
//       hour：1-23
//       day：1-31
//       month：1-12
//       week：0-6（0 means Sunday）

// SetCron some signals：
//       *： any time
//       ,：　 separate signal
//　　    －：duration
//       /n : do as n times of time duration
/////////////////////////////////////////////////////////
//	0/30 * * * * *                        every 30s
//	0 43 21 * * *                         21:43
//	0 15 05 * * * 　　                     05:15
//	0 0 17 * * *                          17:00
//	0 0 17 * * 1                           17:00 in every Monday
//	0 0,10 17 * * 0,2,3                   17:00 and 17:10 in every Sunday, Tuesday and Wednesday
//	0 0-10 17 1 * *                       17:00 to 17:10 in 1 min duration each time on the first day of month
//	0 0 0 1,15 * 1                        0:00 on the 1st day and 15th day of month
//	0 42 4 1 * * 　 　                     4:42 on the 1st day of month
//	0 0 21 * * 1-6　　                     21:00 from Monday to Saturday
//	0 0,10,20,30,40,50 * * * *　           every 10 min duration
//	0 */10 * * * * 　　　　　　              every 10 min duration
//	0 * 1 * * *　　　　　　　　               1:00 to 1:59 in 1 min duration each time
//	0 0 1 * * *　　　　　　　　               1:00
//	0 0 */1 * * *　　　　　　　               0 min of hour in 1 hour duration
//	0 0 * * * *　　　　　　　　               0 min of hour in 1 hour duration

const (
	SECOND = iota
	MINUTE
	HOUR
	DAY
	MONTH
	WEEK
)

type Callback func(id, format string, extra map[string]interface{}) error
type Task struct {
	format   string
	callback Callback
	extra    map[string]interface{}
}
type Crontab struct {
	adapter CrontabAdapterInterface
	tasks   *sync.Map
}

func New() *Crontab {
	cron := &Crontab{}
	cron.Setup(&Adapter{})
	return cron
}
// todo
func (this *Crontab) Pause() {

}
func (this *Crontab) Reset() {
	this.tasks.Range(func(id, task interface{}) bool {
		_id, _ := id.(string)
		this.tasks.Delete(_id)
		return true
	})
}
func (this *Crontab) List() *sync.Map {
	return this.tasks
}
func (this *Crontab) Unregister(id string) {
	this.tasks.Delete(id)
}
func (this *Crontab) Register(id, format string, task Callback, extra map[string]interface{}) {
	(*this.tasks).Store(id, &Task{
		format:   format,
		callback: task,
		extra:    extra,
	})
}

func (this *Crontab) Run() {
	go this.adapter.Run(this.tasks)
}
func (this *Crontab) Setup(adapter CrontabAdapterInterface) {
	adapter.Init()
	this.adapter = adapter
	this.tasks = &sync.Map{}
}

type CrontabAdapterInterface interface {
	Init() error
	Run(tasks *sync.Map)error
}
