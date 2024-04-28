package main

import (
	"context"
	"fmt"
	"time"
)

func main(){
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	go GetIp3(ctx)
	// 手动结束进程
	time.Sleep(5*time.Second)
	// 模拟线程阻塞
	time.Sleep(1*time.Second)
}

func GetIp3(ctx context.Context){
	fmt.Println("获取IP")
	select {
		case <- ctx.Done():
			fmt.Println("协程取消", ctx.Err())
	}
}