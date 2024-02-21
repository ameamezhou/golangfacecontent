package main

import (
	"fmt"
	"time"
	"unsafe"
)

func main(){
	// ch 是长度为4的带缓冲的 channel
	// 初始 hchan结构体中的buf为空，sendx和recvx均为0
	ch := make(chan string, 4)
	fmt.Println(ch, unsafe.Sizeof(ch))
	go sendTask(ch)
	go receiveTask(ch)
	time.Sleep(1 * time.Second)
}

// G1 是发送者
// 当 G1 向ch里面发送数据的时候，首先会对buf加锁，然后将 task 存储的数据 copy 到 buf 中，然后 sendx++ ， 然后释放对 buf 的锁
func sendTask(ch chan string){
	taskList := []string{"I", "like", "jia jia", "and", "my", "id" , "is", "zhou ", "jia jia"}
	for _, task := range taskList {
		ch <- task
	}
}

// G2 是接收者
// 当 G2 消费 ch 的时候，会首先对 buf 加锁，然后将 buf 中的数据 copy 到 task 变量对应的内存里，然后 recvx++, 并释放锁
func receiveTask(ch chan string){
	for {
		task := <- ch
		fmt.Println("received: ", task)
	}
}
