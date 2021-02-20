package timer

import (
	"time"

	"github.com/pkg/errors"
)

//default Job
type TaskJob struct {
	Fn      func()
	err     chan error
	done    chan bool
	stop    chan bool
	replies map[string]func(job *TaskJob, reply Reply)
	Task    TaskInterface
	errorno error
}

func (j *TaskJob) GetTask() TaskInterface {
	return j.Task
}

func (j *TaskJob) SetTask(Task TaskInterface) {
	j.Task = Task
}

func (j *TaskJob) OnStart(f func(job *TaskJob, reply Reply)) {
	j.replies["start"] = f
}

func (j *TaskJob) OnStop(f func(job *TaskJob, reply Reply)) {
	j.replies["stop"] = f
}

func (j *TaskJob) OnFinish(f func(job *TaskJob, reply Reply)) {
	j.replies["finish"] = f
}

func (j *TaskJob) OnError(f func(job *TaskJob, reply Reply)) {
	j.replies["error"] = f
}

func (j *TaskJob) Finish() {
	if f, ok := j.replies["finish"]; ok {
		f(j, Reply{})
	}
}

func (j *TaskJob) run() {
	isPanic := false
	defer func() {
		if x := recover(); x != nil {
			err := errors.Errorf("job error with panic:%v", x)
			j.err <- err
			isPanic = true
			return
		}
	}()

	defer func() {
		if !isPanic {
			j.done <- true
		}
	}()
	j.Fn()
}

func (j *TaskJob) Stop() {
	j.stop <- true
}

//run job and catch error
func (j *TaskJob) Run() {
	if f, ok := j.replies["start"]; ok {
		f(j, Reply{})
	}

	go j.run()
	for {
		select {
		case e := <-j.err:
			if f, ok := j.replies["error"]; ok {
				reply := Reply{
					Code: 500,
					Msg:  e.Error(),
					Err:  e,
				}
				f(j, reply)
			}
			j.errorno = e
			return
		case <-j.done:
			task := j.GetTask()
			loop := false
			now := time.Now().UnixNano()
			if task.GetRunNumber() > 1 {
				task.SetRunNumber(task.GetRunNumber() - 1)
				loop = true
			} else if task.GetEndTime() > now {
				loop = true
			}

			if loop {
				spacing := task.GetSpacing()
				if spacing > 0 {
					//must use now time
					task.SetRuntime(now + spacing)
					task.SetStatus(0)
				}
			} else {
				j.close(false)
				j.Finish()
			}
			return
		case <-j.stop:
			j.close(true)
			if f, ok := j.replies["stop"]; ok {
				f(j, Reply{})
			}
			return
		}
	}
}

func (j *TaskJob) close(exit bool) {
	close(j.done)
	close(j.err)
	close(j.stop)
}
