package pool

import (
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/galaxy-future/BridgX/internal/model"
)

var expandWorkerPool gopool.Pool
var shrinkWorkerPool gopool.Pool
var ExpandTasksChan = make(chan *model.Task, 100)
var ShrinkTasksChan = make(chan *model.Task, 100)

func init() {
	expandWorkerPool = gopool.NewPool("expand-worker-pool", 100, gopool.NewConfig())
	shrinkWorkerPool = gopool.NewPool("shrink-worker-pool", 100, gopool.NewConfig())
	go daemon()
}

func daemon() {
	for {
		select {
		case et, ok := <-ExpandTasksChan:
			if ok {
				expandWorkerPool.Go(func() {
					doExpand(et)
				})
			}
		case st, ok := <-ShrinkTasksChan:
			if ok {
				shrinkWorkerPool.Go(func() {
					doShrink(st)
				})
			}
		}
	}
}
