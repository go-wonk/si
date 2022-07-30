package sicore

import (
	"sync"
)

type WorkerJob interface {
	Execute() (any, error)
}

type WorkerPool struct {
	numWorkers int
	queueSize  int
	jobsWg     *sync.WaitGroup
	jobs       chan WorkerJob
	resultsWg  *sync.WaitGroup
	results    chan any
	errorzWg   *sync.WaitGroup
	errorz     chan error
}

func NewPool(numWorkers, queueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		queueSize:  queueSize,
		jobsWg:     &sync.WaitGroup{},
		jobs:       make(chan WorkerJob, queueSize),
		resultsWg:  &sync.WaitGroup{},
		errorzWg:   &sync.WaitGroup{},
	}
}

func NewPoolWithResultsAndErrors(numWorkers, queueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		queueSize:  queueSize,

		jobsWg: &sync.WaitGroup{},
		jobs:   make(chan WorkerJob, queueSize),

		resultsWg: &sync.WaitGroup{},
		results:   make(chan any, queueSize),

		errorzWg: &sync.WaitGroup{},
		errorz:   make(chan error, queueSize),
	}
}

func (p *WorkerPool) startWorker() {
	defer p.jobsWg.Done()

	for j := range p.jobs {
		if j == nil {
			continue
		}

		res, err := j.Execute()
		// if err != nil {
		// 	fmt.Println(err)
		// }
		if p.errorz != nil {
			p.errorz <- err
		}
		if p.results != nil {
			p.results <- res
		}
	}
}

func (p *WorkerPool) Start() {
	p.jobsWg.Add(p.numWorkers)
	for i := 0; i < p.numWorkers; i++ {
		go p.startWorker()
	}
}

func (p *WorkerPool) Wait() {
	p.jobsWg.Wait()
	p.resultsWg.Wait()
	p.errorzWg.Wait()
}

func (p *WorkerPool) StartAndWait() {
	p.Start()
	p.Wait()
}

func (p *WorkerPool) Finish() {
	close(p.jobs)
	p.jobsWg.Wait()

	if p.errorz != nil {
		close(p.errorz)
	}
	p.errorzWg.Wait()
	if p.results != nil {
		close(p.results)
	}
	p.resultsWg.Wait()
}

func (p *WorkerPool) Queue() chan<- WorkerJob {
	return p.jobs
}

func (p *WorkerPool) QueueIgnore(j WorkerJob) {

	select {
	case p.jobs <- j:
	default:
		// log.Println("ignored")
		return
	}

	// this probably causes a race condition.
	// if len(p.jobs) == cap(p.jobs) {
	// 	log.Println("ignored")
	// 	return
	// }
	// p.jobs <- j
}

func (p *WorkerPool) Results() <-chan any {
	return p.results
}

func (p *WorkerPool) ResultsReady() {
	p.resultsWg.Add(1)
}

func (p *WorkerPool) ResultsDone() {
	p.resultsWg.Done()
}

func (p *WorkerPool) Errorz() <-chan error {
	return p.errorz
}

func (p *WorkerPool) ErrorzReady() {
	p.errorzWg.Add(1)
}

func (p *WorkerPool) ErrorzDone() {
	p.errorzWg.Done()
}
