package goroutinepool

import "runtime"

type GoroutinePool interface {
	Put(job Job)
	PutDone()
	Run()
}

type Job interface {
	Exec()
}

type pool struct {
	cpuNum    int
	workerNum int
	receiveCh chan Job
	handleCh  chan Job
}

func NewPool(workerNum int) GoroutinePool {
	return newPool(1, workerNum)
}

func NewPoolWithCPUNum(cpuNum, workerNum int) GoroutinePool {
	return newPool(cpuNum, workerNum)
}

func newPool(cpu, worker int) GoroutinePool {
	return &pool{
		cpuNum:    cpu,
		workerNum: worker,
		receiveCh: make(chan Job),
		handleCh:  make(chan Job),
	}
}

func (p *pool) Put(job Job) {
	p.receiveCh <- job
}

func (p *pool) PutDone() {
	close(p.receiveCh)
}

func (p *pool) Run() {
	runtime.GOMAXPROCS(p.cpuNum)
	for i := 0; i < p.workerNum; i++ {
		go p.worker()
	}
	for job := range p.receiveCh {
		p.handleCh <- job
	}
	close(p.handleCh)
}

func (p *pool) worker() {
	for job := range p.handleCh {
		job.Exec()
	}
}

type job struct {
	fn  func(arg interface{}) error
	arg interface{}
}

func NewJob(fn func(arg interface{}) error, arg interface{}) Job {
	return &job{
		fn:  fn,
		arg: arg,
	}
}

func (j *job) Exec() {
	err := j.fn(j.arg)
	if err != nil {
		println(err)
	}
}
