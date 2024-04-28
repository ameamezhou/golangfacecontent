package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 很常见的一个案例，假设有一个获取ip的协程，但是这是一个非常耗时的操作每用户随时可能会取消
// 如果用户取消了，那么之前那个获取协程的函数就要停止了

var Wait = sync.WaitGroup{}

func main2()  {
	t := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	Wait.Add(1)
	go func() {
		// Wait.Done()
		ip, err := GetIp1(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(ip, err)
	}()
	go func() {
		time.Sleep(2*time.Second)
		// 取消协程
		cancel()
	}()
	Wait.Wait()
	fmt.Println("执行结束:", time.Since(t))
}

func GetIp1(ctx context.Context)(ip string, err error){
	go func() {
		select {
		case <- ctx.Done():
			fmt.Println("协程取消", ctx.Err())
			err = ctx.Err()
			Wait.Done()
			return
		}
	}()
	defer Wait.Done()
	time.Sleep(4*time.Second)

	ip = "192.16.8.0.1"

	return
}