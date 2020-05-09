package util

// Agent represent workers agent.
type Agent struct {
	Workers  []Worker
	WorkChan chan *Work
}

// Do does work by registered func.
func (a *Agent) Do(in interface{}) chan interface{} {
	work := &Work{
		Wait:         in,
		DoneWorkChan: make(chan interface{}),
	}
	a.WorkChan <- work
	return work.DoneWorkChan
}

// NewAgent news an agent.
func NewAgent(workerNum, workBuf int, handler Handler) *Agent {
	workChan := make(chan *Work, workBuf)

	workers := make([]Worker, 0, workerNum)
	for i := 0; i < workerNum; i++ {
		worker := Worker{
			WorkChan: workChan,
			Handler:  handler,
		}
		worker.Start()
		workers = append(workers, worker)
	}

	return &Agent{
		WorkChan: workChan,
		Workers:  workers,
	}

}

// Handler represent how to do the work.
type Handler func(interface{}) interface{}

// Worker represent a worker.
type Worker struct {
	WorkChan chan *Work
	Handler  Handler
}

// Work represent a work.
type Work struct {
	Wait         interface{}
	DoneWorkChan chan interface{}
}

// Start let the worker ready to deal with coming work.
func (w *Worker) Start() {
	go func() {
		for {
			work := <-w.WorkChan
			Done := w.Handler(work.Wait)
			work.DoneWorkChan <- Done
		}
	}()
}
