package util

type Agent struct {
	Workers  []Worker
	WorkChan chan *Work
}

func (a *Agent) Do(in interface{}) chan interface{} {
	work := &Work{
		Wait:         in,
		DoneWorkChan: make(chan interface{}),
	}
	a.WorkChan <- work
	return work.DoneWorkChan
}

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

type Handler func(interface{}) interface{}

type Worker struct {
	WorkChan chan *Work
	Handler  Handler
}

type Work struct {
	Wait         interface{}
	DoneWorkChan chan interface{}
}

func (w *Worker) Start() {
	go func() {
		for {
			work := <-w.WorkChan
			Done := w.Handler(work.Wait)
			work.DoneWorkChan <- Done
		}
	}()
}
