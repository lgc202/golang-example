package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants"
)

func wrapper(i int, wg *sync.WaitGroup) func() {
	return func() {
		fmt.Printf("hello from task:%d\n", i)
		// 为了避免任务执行过快，空出了 goroutine，观察不到现象，
		// 每个任务中我使用time.Sleep(2 * time.Second)休眠 2s
		time.Sleep(2 * time.Second)
		wg.Done()
	}
}

func main() {
	// 协程池的大小只有2,但后面提交了3个任务,采用非阻塞模式会用一个提交不成功返回错误
	p, _ := ants.NewPool(2, ants.WithNonblocking(true))
	defer p.Release()

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 1; i <= 3; i++ {
		err := p.Submit(wrapper(i, &wg))
		if err != nil {
			fmt.Printf("task:%d err:%v\n", i, err)
			wg.Done()
		}
	}

	wg.Wait()
}
