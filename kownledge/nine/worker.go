package main

import (
	"fmt"
	"time"
)

type worker struct {
	id  int
	err error
}

func (wk *worker) work(workerChan chan<- *worker) (err error) {
	// 任何Goroutine只要异常退出或者正常退出 都会调用defer 函数，所以在defer中想WorkerManager的WorkChan发送通知
	defer func() {
		//捕获异常信息，防止panic直接退出
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				wk.err = err
			} else {
				wk.err = fmt.Errorf("Panic happened with [%v]", r)
			}
		} else {
			wk.err = err
		}

		//通知 主 Goroutine，当前子Goroutine已经死亡
		workerChan <- wk
	}()

	// do something
	fmt.Println("Start Worker...ID = ", wk.id)

	// 每个worker睡眠一定时间之后，panic退出或者 Goexit()退出
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
	}

	panic("worker panic..")
	//runtime.Goexit()

	return err
}
