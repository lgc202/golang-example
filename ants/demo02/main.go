package main

import (
	"fmt"
	"sync"

	"github.com/panjf2000/ants"
	"golang.org/x/exp/rand"
)

type Task struct {
	index int
	nums  []int
	sum   int
}

const (
	// 10000个整数
	DataSize = 10000
	// 将整数分成 100 份, 每份 100 个
	DataPerTask = 100
)

type taskFunc func()

func taskFuncWrapper(nums []int, i int, sum *int, wg *sync.WaitGroup) taskFunc {
	return func() {
		for _, num := range nums[i*DataPerTask : (i+1)*DataPerTask] {
			*sum += num
		}

		fmt.Printf("task:%d sum:%d\n", i+1, *sum)
		wg.Done()
	}
}

func main() {
	p, _ := ants.NewPool(10)
	defer p.Release()

	nums := make([]int, DataSize, DataSize)
	for i := range nums {
		nums[i] = rand.Intn(1000)
	}

	var wg sync.WaitGroup
	wg.Add(DataSize / DataPerTask)
	partSums := make([]int, DataSize/DataPerTask, DataSize/DataPerTask)
	for i := 0; i < DataSize/DataPerTask; i++ {
		// Submit 的参数是一个不接受任何参数的函数
		// 所以 taskFuncWrapper 需要用闭包处理输入的数据
		p.Submit(taskFuncWrapper(nums, i, &partSums[i], &wg))
	}
	wg.Wait()

	// 将结果汇总
	var sum int
	for _, partSum := range partSums {
		sum += partSum
	}

	// 验证结果
	var expect int
	for _, num := range nums {
		expect += num
	}

	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks, result is %d expect:%d\n", sum, expect)
}
