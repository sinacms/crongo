package crongo

import (
	"fmt"
	"time"
	"testing"
)

func init(){

}

func taskEnqueue(id, format string, extra map[string]interface{})error{
	fmt.Println("task1 " + id +" "+ format+" " + time.Now().String())
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
func taskEnqueue6(id, format string, extra map[string]interface{})error{
	fmt.Println("task6 " + id +" "+ format+" " + time.Now().String())
	return nil
}
var cron *Crontab
func TestCrontab(t *testing.T)  {
	cron = New()
	cron.Register("task1", "7 * * * * *", taskEnqueue, nil)
	//test *
	cron.Register("task2", "0-5 * * * * *", taskEnqueue2, nil)
	//test current minute
	cron.Register("task3", "* * * * * *", taskEnqueue3, nil)
	//test endDate
	cron.Register("task4", "0 * * * * *", taskEnqueue4, map[string]interface{}{"endDate": time.Now().Unix()+65})
	cron.Register("task5", "*/3 * * * * *", taskEnqueue5,nil)
	cron.Register("task6", "8,9,10 * * * * *", taskEnqueue6,nil)

	var sign chan int

	go func(){
		time.Sleep(6*time.Second)
		cron.Unregister("task3")
		fmt.Println("Unregister task3 "+time.Now().String())
		time.Sleep(80*time.Second)
		cron.Reset()
		fmt.Println("Reset "+time.Now().String())

	}()
	cron.Run()
	<-sign

}