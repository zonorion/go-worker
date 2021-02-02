package worker

import (
	"fmt"
	"time"
)

type Job interface {
	Process()
}

var jobRunning = make(chan Job)

type Worker struct {
	Id         int
	Done       chan bool
	JobRunning chan Job
}

func newWorker(id int, jobChan chan Job) *Worker {
	return &Worker{
		Id:         id,
		Done:       make(chan bool),
		JobRunning: jobChan,
	}
}

func (w *Worker) run() {
	fmt.Println("Run worker id: ", w.Id)
	go func() {
		for {
			select {
			case job := <-w.JobRunning:
				fmt.Println("Job running: ", w.Id)
				job.Process()
			case <-w.Done:
				fmt.Println("Stop worker: ", w.Id)
				return
			}
		}
	}()
}

func (w *Worker) stop() {
	w.Done <- true
}

type JobQueue struct {
	Workers    []*Worker
	JobRunning chan Job
	Done       chan bool
}

func NewJobQueue(numberOfWorkers int) JobQueue {
	workers := make([]*Worker, numberOfWorkers, numberOfWorkers)
	jobRunning = make(chan Job)

	for i := 0; i < numberOfWorkers; i++ {
		workers[i] = newWorker(i+1, jobRunning)
	}

	return JobQueue{
		Workers:    workers,
		JobRunning: jobRunning,
		Done:       make(chan bool),
	}
}

func (jq *JobQueue) Stop()  {
	jq.Done <- true
}

func (jq *JobQueue) Push(job Job) {
	jq.JobRunning <- job
}

func (jq *JobQueue) Start() {
	go func() {
		for i := 0; i < len(jq.Workers); i++ {
			jq.Workers[i].run()
		}
	}()

	go func() {
		for {
			select {
			case <-jq.Done:
				for i := 0; i < len(jq.Workers); i++ {
					jq.Workers[i].stop()
				}
				close(jobRunning)
				return
			}
		}
	}()
}

//implement Job interface
type Sender struct {
	Email string
}

func(s Sender) Process()  {
	fmt.Println("========================" + s.Email)
	time.Sleep(20 * time.Second)
}

func IsWorkerRunning() bool {
	select {
	case <-jobRunning:
		return false
	default:
	}

	return true
}

func init()  {
	close(jobRunning)
}
