package comF

import (
	"fmt"
	"sync/atomic"
	"time"
)

const (
	noStop   = 0
	stop     = 1
	stopTime = time.Second * 7
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
		taskChan: make(chan *TaskArgs, max*10),
		workChan: make(chan *workerTool, max),
		stop:     make(chan bool),
	}

	for i := 1; i <= init; i++ {
		worker := worker{
			id:       i,
			dispatch: &dispatch,
			workTool: &workerTool{
				workChan: make(chan *TaskArgs),
				stop:     noStop,
			},
		}
		atomic.AddInt32(&dispatch.cutCount, 1)
		go worker.run(nil, true)
	}

	go dispatch.dispatch()

	return &dispatch
}

func DestroyPool(pool Pool) {
	p := pool.(*dispatcher)
	close(p.stop)
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
	taskChan  chan *TaskArgs
	workChan  chan *workerTool
	stop      chan bool
}

func (d *dispatcher) dispatch() {
	for {
		select {
		case v := <- d.taskChan:
			select {
			case work := <- d.workChan:
				if atomic.LoadInt32(&work.stop) == noStop {
					work.workChan <- v
					continue
				}
			default:
			}
			if atomic.LoadInt32(&d.cutCount) < int32(d.max) {
				w := worker{
					id:       int(GetWorkId()),
					dispatch: d,
					workTool: &workerTool{
						workChan: make(chan *TaskArgs),
						stop:     noStop,
					},
				}
				atomic.AddInt32(&d.cutCount, 1)
				go w.run(v, false)
				continue
			}
			for work := range d.workChan {
				if atomic.LoadInt32(&work.stop) == noStop {
					work.workChan <- v
					break
				}
			}
		case <- d.stop:
			goto EXIT
		}
	}

EXIT:
	fmt.Println("退出dispatcher")
}

func (d *dispatcher) NewAsyncTask(ex AsyncF, args ...interface{}) (success bool) {
	task := &TaskArgs{
		Ex:   ex,
		Args: args,
	}

	d.taskChan <- task

	return true
}

type worker struct {
	id       int
	dispatch *dispatcher
	workTool *workerTool
}

type workerTool struct {
	workChan chan *TaskArgs
	stop     int32
}

func (w *worker) run(t *TaskArgs, isInit bool) {
	fmt.Println("创建worker: ", w.id)
	if t != nil {
		t.Ex(t.Args...)
	}
	if isInit {
		for {
			w.dispatch.workChan <- w.workTool
			select {
			case task := <- w.workTool.workChan:
				task.Ex(task.Args...)
			case <- w.dispatch.stop:
				goto EXIT
			}
		}
	} else {
		for {
			w.dispatch.workChan <- w.workTool
			select {
			case task := <- w.workTool.workChan:
				task.Ex(task.Args...)
			case <- w.dispatch.stop:
				goto EXIT
			case <- time.After(stopTime):
				atomic.AddInt32(&w.workTool.stop, stop)
				atomic.AddInt32(&w.dispatch.cutCount, -1)
				select {
				case task := <- w.workTool.workChan:
					task.Ex(task.Args...)
				default:
					goto EXIT
				}
			}
		}
	}
EXIT:
	fmt.Printf("退出worker ID:%d\n", w.id)
}
