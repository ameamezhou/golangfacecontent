package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main0(){
	var wg = sync.WaitGroup{}

	ctx1, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	wg.Add(1)
	go GetIp(ctx1, &wg)
	wg.Wait()
}

func GetIp(ctx context.Context, wg *sync.WaitGroup)(ip string, err error){
	go func() {
		select {
		case <- ctx.Done():
			fmt.Println("协程取消", ctx.Err())
			err = ctx.Err()
			wg.Done()
			return
		}
	}()
	defer wg.Done()
	time.Sleep(7*time.Second)

	ip = "192.16.8.0.1"

	return
}