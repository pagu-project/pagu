package job

import (
	"sync"
)

type IScheduler interface {
	Submit(job Job)
	Run()
	Shutdown()
}

type Scheduler struct {
	jobs []Job
	mu   *sync.Mutex
	wg   *sync.WaitGroup
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: make([]Job, 0),
		mu:   &sync.Mutex{},
		wg:   &sync.WaitGroup{},
	}
}

func (s *Scheduler) Submit(job Job) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Run() {
	for _, j := range s.jobs {
		s.wg.Add(1)
		go j.Start()
	}
}

func (s *Scheduler) Shutdown() {
	for _, j := range s.jobs {
		j.Stop()
		s.wg.Done()
	}
}
