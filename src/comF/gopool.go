package comF

import (
	"fmt"
	"sync/atomic"
)

type AsyncF func(args ...interface{})

type Pool interface {
	NewAsyncTask(ex AsyncF, args ...interface{}) (success bool)
}

func NewPool(init, max int) Pool {
	dispatch := dispatcher{
		id:       GetUUID(),
		init:     init,
		max:      max,
		taskChan: make(chan TaskArgs, max*10),
		workChan: make(chan chan TaskArgs, max),
		stop:     make(chan bool),
	}

	for i := 1; i <= init; i++ {
		worker := worker{
			id:       i,
			dispatch: &dispatch,
			workChan: make(chan TaskArgs),
			stop:     make(chan bool),
		}
		dispatch.workers = append(dispatch.workers, worker)
		dispatch.cutCount++
		go worker.run(nil)
	}

	go dispatch.dispatch()

	return &dispatch
}

func DestroyPool(pool Pool) {
	p := pool.(*dispatcher)
	for _, v := range p.workers {
		v.stop <- true
	}
	p.stop <- true
}

type TaskArgs struct {
	Ex   AsyncF
	Args []interface{}
}

type dispatcher struct {
	id        string
	init      int
	max       int
	cutCount  int32
	busyCount int32
	taskChan  chan TaskArgs
	workChan  chan chan TaskArgs
	workers   []worker
	stop      chan bool
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case v := <- d.taskChan:
			if atomic.LoadInt32(&d.busyCount) >= d.cutCount {
				if d.cutCount >= int32(d.max) {
					fmt.Println("max workers wait to run task")
				} else {
					w := worker{
						id:       int(d.cutCount + 1),
						dispatch: d,
						workChan: make(chan TaskArgs),
						stop:     make(chan bool),
					}
					d.workers = append(d.workers, w)
					d.cutCount++
					go w.run(&v)
					continue
				}
			}
			work := <- d.workChan
			work <- v
		case <- d.stop:
			goto EXIT
		}
	}

EXIT:
	fmt.Println("退出dispatcher")
}

func (d *dispatcher) NewAsyncTask(ex AsyncF, args ...interface{}) (success bool) {
	task := TaskArgs{
		Ex:   ex,
		Args: args,
	}

	d.taskChan <- task

	return true
}

type worker struct {
	id       int
	dispatch *dispatcher
	workChan chan TaskArgs
	stop     chan bool
}

func (w *worker) run(t *TaskArgs) {
	fmt.Println("创建worker", w.id)
	if t != nil {
		atomic.AddInt32(&w.dispatch.busyCount, 1)
		t.Ex(t.Args...)
		atomic.AddInt32(&w.dispatch.busyCount, -1)
	}
	for {
		w.dispatch.workChan <- w.workChan
		select {
		case task := <- w.workChan:
			atomic.AddInt32(&w.dispatch.busyCount, 1)
			task.Ex(task.Args...)
			atomic.AddInt32(&w.dispatch.busyCount, -1)
		case <- w.stop:
			goto EXIT
		}
	}

EXIT:
	fmt.Printf("退出worker ID:%d\n", w.id)
}
