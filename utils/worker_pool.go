package utils

import (
	"log"
	"sync"

	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/model"
)

type IworkerPool interface {
	Run()
	Close()
	AddWorker(w Iworker)
	GetWorkers() []Iworker
}

type Iworker interface {
	Execute(in <-chan []string, out chan<- *model.Pokemon)
	Id() int
}

type workerPool struct {
	waitGroup *sync.WaitGroup
	workers   []Iworker
	inChan    <-chan []string
	outChan   chan<- *model.Pokemon
}

func NewWorkerPool(wg *sync.WaitGroup, n int, in <-chan []string, out chan<- *model.Pokemon) IworkerPool {
	wp := &workerPool{
		waitGroup: wg,
		workers:   make([]Iworker, 0, n),
		inChan:    in,
		outChan:   out,
	}

	return wp
}

func (wp *workerPool) Run() {
	for _, worker := range wp.workers {
		log.Printf("WorkerPool: worker %v has been spawned.\n", worker.Id())
		wp.waitGroup.Add(1)
		go func(w Iworker) {
			defer wp.waitGroup.Done()
			log.Printf("WorkerPool: worker %v has started.\n", w.Id())
			w.Execute(wp.inChan, wp.outChan)
			log.Printf("WorkerPool: worker %v has completed.\n", w.Id())
		}(worker)
	}
}

func (wp *workerPool) AddWorker(w Iworker) {
	wp.workers = append(wp.workers, w)
}

func (wp *workerPool) Close() {
	close(wp.outChan)
}

func (wp *workerPool) GetWorkers() []Iworker {
	return wp.workers
}
