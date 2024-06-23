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
	// WithMaxBlockingTasks 设置最大等待队列长度为 2
	p, _ := ants.NewPool(4, ants.WithMaxBlockingTasks(2))
	defer p.Release()

	var wg sync.WaitGroup
	wg.Add(8)
	for i := 1; i <= 8; i++ {
		// 提交任务必须并行进行。如果是串行提交，第 5 个任务提交时由于池中没有空闲的 goroutine 处理该任务，
		// Submit()方法会被阻塞，后续任务就都不能提交了。也就达不到验证的目的了
		go func(i int) {
			err := p.Submit(wrapper(i, &wg))
			if err != nil {
				fmt.Printf("task:%d err:%v\n", i, err)
				// 由于任务可能提交失败，失败的任务不会实际执行，所以实际上wg.Done()次数会小于 8。
				// 因而在err != nil分支中我们需要调用一次wg.Done()。否则wg.Wait()会永远阻塞
				wg.Done()
			}
		}(i)
	}

	wg.Wait()
}
