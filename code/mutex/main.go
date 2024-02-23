package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// 可重入锁

type ReentrantLock struct {
	sync.Mutex
	recursion	int32 // goroutine 可重入的次数
	owner 		int64 // 当前持有锁的 goroutine id
}

// get returns the id of the current goroutine.
func GetGoroutineID() int64 {

	var buf [64]byte
	// 获取栈信息
	n := runtime.Stack(buf[:], false)
	// 抽取id
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine"))[0]
	// 转为64位整数
	gid, _ := strconv.Atoi(idField)
	return int64(gid)


	//var buf [64]byte
	//var s = buf[:runtime.Stack(buf[:], false)]
	//s = s[len("goroutine "):]
	//s = s[:bytes.IndexByte(s, )]
	//gid, _ := strconv.ParseInt(string(s), 10, 64)
	//return gid
}

func NewReentrantLock() sync.Locker {
	res := &ReentrantLock{
		Mutex: sync.Mutex{},
		recursion: 0,
		owner: 0,
	}
	return res
}

// ReentrantMutex 包装一个 Mutex 实现可重入
type ReentrantMutex struct {
	sync.Mutex
	recursion	int32 // goroutine 可重入的次数
	owner 		int64 // 当前持有锁的 goroutine id
}
func (m *ReentrantMutex) Lock(){
	gid := GetGoroutineID()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的 goroutine 第一次调用 记录下它的goroutine id 调用次数加1
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *ReentrantMutex) ULock(){
	gid := GetGoroutineID()
	// 非持有锁的goroutine尝试释放锁，错误使用
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("worng the owner(%d): %d", m.owner, gid))
	}
	// 调用次数减1
	m.recursion--
	if m.recursion != 0 {
		return
	}
	// 最后一次调用，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}


func main(){
	var mutex = &ReentrantMutex{}
	mutex.Lock()
	mutex.Lock()
	fmt.Println(111)
	mutex.Unlock()
	mutex.Unlock()
}
