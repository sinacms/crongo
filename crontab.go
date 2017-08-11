package crongo

import "sync"

func init(){

}

const (
	SECOND = iota
	MINUTE
	HOUR
	DAY
	MONTH
	WEEK
)


type Callback func(id , format string, extra map[string]interface{})error
type Task  struct{
	format string
	callback Callback
	extra map[string]interface{}
}
type Crontab struct {
	adapter CrontabAdapterInterface
	//tasks *map[string]*Task 不安全
	tasks *sync.Map
}
func New()*Crontab{
	return &Crontab{}
}
func (this *Crontab)Reset() {
	this.tasks.Range(func(id, task interface{})bool{
		_id,_ := id.(string)
		this.tasks.Delete(_id)
		return true
	})
}
func (this *Crontab)Unregister(id string) {
	//if _, ok :=  (*this.tasks).Load(id); !ok {
		//delete(*this.tasks, id)
		this.tasks.Delete(id)
	//}
}
func (this *Crontab)Register(id, format string, task Callback, extra map[string]interface{})  {
	//if _, ok :=  (*this.tasks).Load(id); !ok {
	(*this.tasks).Store(id, &Task{
		format: format,
		callback: task,
		extra: extra,
	})
	//}
}

func (this *Crontab)Run() {
	go this.adapter.Run(this.tasks)
}
func (this *Crontab)Setup(adapter CrontabAdapterInterface){
	adapter.Init()
	this.adapter = adapter
	this.tasks = &sync.Map{}
}


type CrontabAdapterInterface interface {
	Init()error
	Run(tasks *sync.Map)
}
