# crongo

    This is a crontab service that supports hot plug and high performance. In addition, it supports second level (first parameters), multi time point and time section function.


    这是一个支持热插拨、高性能的crontab服务，另外，它还支持秒级别（第1个参数）、多时间点、时间段功能。

# Usage

``` go
	cron := New()
	cron.Register("task1", "7 * * * * *", taskEnqueue, nil)
	cron.Register("task2", "0-5 * * * * *", taskEnqueue2, nil)
	cron.Register("task3", "* * * * * *", taskEnqueue3, nil)
	cron.Register("task4", "0 * * * * *", taskEnqueue4, map[string]interface{}{"endDate": time.Now().Unix()+65})
	cron.Register("task5", "*/3 * * * * *", taskEnqueue5,nil)
	cron.Register("task6", "8,9,10 * * * * *", taskEnqueue6,nil)

	go func(){
		time.Sleep(6*time.Second)
		cron.Unregister("task3")
		fmt.Println("Unregister task3 "+time.Now().String())
		time.Sleep(80*time.Second)
		cron.Reset()
		fmt.Println("Reset "+time.Now().String())
	}()
	cron.Run()


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
```
## format 格式
```
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
```
