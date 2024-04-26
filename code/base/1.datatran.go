package main

import (
	"context"
	"fmt"
)

type UserInfo struct {
	Name 	string
	Age 	int
}

func GetUser(ctx context.Context){
	fmt.Println(ctx.Value("info").(UserInfo).Name) // 可以使用断言转化类型的
}

func main1(){
	ctx := context.Background()
	ctx = context.WithValue(ctx, "info", UserInfo{Name: "xiaoqizhou", Age: 18})
	GetUser(ctx)
}
