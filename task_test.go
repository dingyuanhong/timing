package timer

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	fmt.Println("TestTime")
	startTime := "2019-08-11 21:30:01"
	endTime := "2019-08-11 21:45:00"
	cal := "2019-08-12 21:30:07"

	t1, err := time.ParseInLocation("2006-01-02 15:04:05", cal, time.Local)
	s, err1 := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	e, err2 := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)

	sTime := time.Now().Format("2006-01-02") + " " + s.Format("15:04:05")
	eTime := time.Now().Format("2006-01-02") + " " + e.Format("15:04:05")
	cal = time.Now().Format("2006-01-02") + " " + t1.Format("15:04:05")

	fmt.Println(t1, err, err1, err2, sTime, eTime)
	if cal >= sTime && cal <= eTime {
		fmt.Println(cal)
	}

}

//test add Func
func Test_AddFunc(t *testing.T) {
	fmt.Println("Test_AddFunc")
	cron := NewScheduler()

	go cron.Start()

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*1), func() {
		fmt.Println("one second after")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*1), func() {
		fmt.Println("one second after, task second")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("ten second after")
	})

	timer := time.NewTimer(11 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("over")
		}
		break
	}
	cron.Stop()
}

//test add Task and timing add Task
func Test_AddTask(t *testing.T) {
	fmt.Println("Test_AddTask")
	cron := NewScheduler()
	go cron.Start()

	cron.AddTask(&Task{
		Job: GetJob(func() {
			fmt.Println("hello cron")
		}),
		RunTime: time.Now().UnixNano() + int64(time.Second*2),
	})

	cron.AddTask(&Task{
		Job: GetJob(func() {
			fmt.Println("hello cron1")
		}),
		RunTime: time.Now().UnixNano() + int64(time.Second*3),
	})

	cron.AddTask(&Task{
		Job: GetJob(func() {
			fmt.Println("hello cron2")
		}),
		RunTime: time.Now().UnixNano() + +int64(time.Second*4),
	})

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("over")
		}
		break
	}
	cron.Stop()
}

//test add space task func
func Test_AddFuncSpace(t *testing.T) {
	fmt.Println("Test_AddFuncSpace")
	cron := NewScheduler()
	log.SetOutput(os.Stdout)
	go cron.Start()

	cron.AddFuncSpace(int64(time.Second*1), time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("one second after")
	})

	cron.AddFuncSpace(int64(time.Second*1), time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("one second after, task second")
	})

	cron.AddFunc(time.Now().UnixNano()+int64(time.Second*10), func() {
		fmt.Println("ten second after")
	})

	cron.AddFuncSpaceNumber(int64(time.Second*1), 10, func() {
		fmt.Println("number 10")
	})
	timer := time.NewTimer(11 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("over")
		}
		break
	}
	cron.Stop()
}

func Test_JobStartEvent(t *testing.T) {
	fmt.Println("Test_JobStartEvent")
	cron := NewScheduler()
	cron.Start()
	f := func() {
		fmt.Println("hello")
	}
	t1 := &Task{
		Job:     GetJob(f),
		RunTime: time.Now().UnixNano() + int64(time.Second)*1,
		Spacing: int64(3 * time.Second),
		EndTime: time.Now().UnixNano() + int64(time.Second*20),
		Uuid:    "123",
	}
	f1 := func(j *TaskJob, reply Reply) {
		fmt.Println(reply)
		fmt.Println("It's reply")
	}
	t1.GetJob().OnStart(f1)
	cron.AddTask(t1)

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("over")
		}
		break
	}
	cron.Stop()
}

func Test_DoubleJob(t *testing.T) {
	fmt.Println("Test_DoubleJob")

	cron := NewScheduler()
	cron.Start()
	f1 := func() {
		fmt.Println("hello 1")
	}
	f2 := func() {
		fmt.Println("hello 2")
	}

	cron.AddTask(&Task{
		Job:     GetJob(f1),
		RunTime: time.Now().UnixNano(),
		Number:  1,
		Uuid:    "222223",
	})
	cron.AddTask(&Task{
		Job:     GetJob(f2),
		RunTime: time.Now().Add(time.Duration(5) * time.Second).UnixNano(),
		Number:  1,
		Uuid:    "222222",
	})

	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			fmt.Println("over")
		}
		break
	}
	cron.Stop()
}
