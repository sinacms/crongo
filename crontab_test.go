package crongo

import (
	"fmt"
	"time"
	"strconv"
	"testing"
)

func init(){

}

func taskEnqueue(id, format string, extra map[string]interface{})error{
	//time.Sleep(2 * time.Second)
	fmt.Println("task " + id +" "+ format+" " + time.Now().String())
	return nil
}
func taskEnqueue2(id, format string, extra map[string]interface{})error{
	//time.Sleep(2 * time.Second)
	fmt.Println("task2 " + id +" "+ format+" " + time.Now().String())
	return nil
}
func taskEnqueue3(id, format string, extra map[string]interface{})error{
	//time.Sleep(2 * time.Second)
	fmt.Println("task3 " + id +" "+ format+" " + time.Now().String())
	return nil
}
func taskEnqueue4(id, format string, extra map[string]interface{})error{
	if  end,ok := extra["endDate"].(int64);  ok && end < time.Now().Unix()  {
		cron.Unregister(id)
		fmt.Println("task "+id+" Unregister ")
		return nil
	}
	fmt.Println("task4 " + id +" "+ format+" " + time.Now().String())

	return nil
}
func taskEnqueue5(id, format string, extra map[string]interface{})error{
	fmt.Println("task5 " + id +" "+ format+" " + time.Now().String())

	return nil
}
var cron *Crontab
func TestCrontab(t *testing.T)  {
	cron = New()
	cron.Setup(&MinuteAdapter{})
	//cron.Register("123", "* * * * *", taskEnqueue)
	//test *
	cron.Register("123", "* * * * *", taskEnqueue2, nil)
	//test current minute
	cron.Register("task3", strconv.Itoa(time.Now().Minute())+" * * * *", taskEnqueue3, nil)
	//test endDate
	cron.Register("task4", "* * * * *", taskEnqueue4, map[string]interface{}{"endDate": time.Now().Unix()+5})
	cron.Register("task5", "*/5 * * * *", taskEnqueue5,nil)

	var sign chan int

	go func(){
		time.Sleep(3*time.Second)
		cron.Unregister("123")
		fmt.Println("Unregister 123")
		time.Sleep(90*time.Second)
		cron.Reset()
		fmt.Println("Reset ")

	}()
	cron.Run()
	<-sign

}