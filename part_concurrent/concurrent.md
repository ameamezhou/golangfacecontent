## 并发
### Go 常用的并发模型
并发模型说的是系统中的线程如何协作完成并发任务，不同的并发模型，线程以不同的方式进行通信和协作

#### 线程间通信方式
线程间通信的方式有两种：共享内存和消息传递，无论是那种通信模型，线程或者协程最终都会从内存中获取数据，所以梗准确的说法是直接共享内存、发送消息的方式来同步信息

**共享内存**
- **抽象层级**：抽象层级低，当我们遇到资源进行更细粒度的控制或者对性能有极高要求的场景才会考虑抽象层级更低的方法
- **耦合**：高，线程需要在读取或者写入数据时优先获取保护该资源的互斥锁
- **线程竞争**：需要加锁，才能避免线程竞争和数据冲突

**发送消息**
- **抽象层级**：抽象层级高，提供了更良好的封装和生产领域更相关和契合的设计，比如 Go中的channel就提供了goroutine之间用于信息传递的方式，它在内部实现的时候就广泛运用刀了共享内存和锁，通过对二者进行组合提供了更高级的同步
- **耦合**：低，生产者消费者模型
- **线程竞争**：保证同一时间只有要给活跃的线程能够访问数据，channel会维护所有被channel阻塞的协程，保证有资源的时候只唤醒要给协程，避免竞争

Go语言中实现了两种并发模型，一种是共享内存并发模型，另一种是CSP模型

共享内存：

![img.png](img.png)

CSP并发模型

通过发送消息的方式来同步信息，Go语言推荐使用 通信顺序进程 communicating sequential processes 并发模型，通过 goroutine 和 channel 来实现

- goroutine 是Go语言中并发的执行单位，可以理解为"线程"
- channel 是Go语言中各个并发结构体(goroutine)之前的通信机制，通俗的来说，就是各个 goroutine 之间通信的"管道"，类似于 Linux 当中的管道

![img_1.png](img_1.png)

### Go 有哪些并发同步原语
mutex RWMutex 等并发源于的底层是通过 atomic 包中一些原子操作来实现的，原子操作是最基础的并发原语

#### atomic
1. 应用场景
    - 在有些场景中，我们不需要这些基本并发原语里面的复杂逻辑，而只需要其中的部分简单的原子操作，就可以通过atomic包中的方法去实现
    - 可以使用 atomic 实现自己定义的基本并发原语
    - atomic 原子操作还是实现 lock-free 数据结构的基石
2. 提供的方法（https://github.com/ameamezhou/golangfacecontent/blob/master/README.md#go-%E7%9A%84%E5%8E%9F%E5%AD%90%E6%93%8D%E4%BD%9C%E6%9C%89%E5%93%AA%E4%BA%9B）
    - Add：给第一个参数地址中的值增加一个delta值
    - CAS：比较当前addr地址里的值是不是old，如果不等于old就返回false;如果当前值等于old，就替换新值，然后返回true
    - Swap：如果不需要比较旧值，只是比较粗暴的替换的话，就可以使用Swap方法
    - Load：取出addr地址中的值，即使在多处理器、多核、有CPU cache 的情况下，这个操作也能保证 Load 是一个原子操作
    - store: 把一个值存入到指定的 addr 地址中
    - value 类型：可以把原子地存取对象类型，但也只能存取，不能CAS和Swap，常常用在配置变更等场景中
3. 第三方库的扩展 
    - uber-go/atomic：定义和封装了几种常见类型相对应的原子操作类型
    - bool类型：提供了 CAS、store、Swap、Toggle等原子方法，还提供 string、MarshalJson、UnmarshalJson等辅助方法

#### channel

channel 管道，高级同步原语，goroutine之间通信的桥梁

使用场景：消息队列、数据传递、信号通知、任务编排、锁

#### 基本并发原语
Go语言在sync包中提供了用于同步一些基本原语，这些基本原语提供了较为基础的同步功能，但是它们是一种相对原始的同步机制，在多数情况下，我们都应该使用抽象层更高的channel实现同步

常见的并发原语如下：sync.Mutex sync.RWMutex syncWaitGroup sync.Cond sync.Once sync.Pool(内存复用 减少GC压力) sync.Context(上下文信息传递，超时和取消机制、控制子 goroutine的执行) sync.Map

#### 扩展并发原语
ErrGroup

可以在一组Goroutine中提供同步、错误传播以及上下文取消的功能

使用场景：只要一个 goroutine 出错我们就不再等其他的 goroutine了，减少资源浪费，并且返回错误

### Go waitGroup 实现原理
Go标准库提供的原语，可以用它来等待一批Goroutine结束

```go
type WaitGroup struct {
	noCopy noCopy

	// 64-bit value: high 32 bits are counter, low 32 bits are waiter count.
	// 高32位是一个计数器，低32位是等待计数器
	// 64-bit atomic operations require 64-bit alignment, but 32-bit
	// compilers only guarantee that 64-bit fields are 32-bit aligned.
	// 64位原子操作需要64位对齐，但32位编译器只能保证64位字段是32位对齐的。
	// For this reason on 32 bit architectures we need to check in state()
	// if state1 is aligned or not, and dynamically "swap" the field order if
	// needed.
	// 因为这个原因所以我们32位的结构需要用state()方法进行检查，如果state1是否对齐，如果有需要的话会动态”交换“字段顺序
	state1 uint64
	state2 uint32
}
```
其中 noCopy 就是golang源码中检测禁止拷贝的技术。如果程序中有 WaitGroup 的赋值行为，使用 go vet 检查程序的时候就会发现报错。但需要注意的是，noCopy不会影响程序正常的编译和进行

**使用方法**

在WaitGroup里面主要有3个方法：

- WaitGroup.Add(): 可以添加或减少请求的 goroutine 数量，Add(n) 将会导致 counter += n
- WaitGroup.Done(): 相当于 Add(-1),Done()导致counter -= 1，当请求计数器 counter 为0的时候通过信号量调用 runtime_Semrelease 唤醒witer 线程
- WaitGroup.Wait(): 会将 waiter++, 同时通过信号量调用 runtime_Semacquire(semap) 阻塞当前goroutine

### Go Cond 实现原理
Go 标准库提供了 Cond 原语，可以让 goroutine 在满足特定条件的时候被阻塞和唤醒

**底层数据结构**
```go
// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *Mutex or *RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
//
// A Cond must not be copied after first use.
type Cond struct {
	noCopy noCopy

	// L is held while observing or changing the condition
	L Locker

	notify  notifyList
	checker copyChecker
}

type notifyList struct {
    wait   uint32
    notify uint32
    lock   uintptr // key field of the mutex
    head   unsafe.Pointer
    tail   unsafe.Pointer
}
```
- noCopy 同上
- checker：用于禁止运行期间发生拷贝，双重检查
- L：可以传入一个读写锁或互斥锁，当修改条件或者调用 Wait 方法的时候需要加锁
- notify：通知链表，调用 Wait() 方法的 Goroutine 会放到这个链表中，从这里获取需要被唤醒的 Goroutine 列表

**方法**
- sync.NewCond(l Locker): 创建一个Cond变量，Locker是一个必填参数，这在 cond.Wait 里面涉及到 locker 的锁操作
- Cond.Wait(): 阻塞等待被唤醒，调用 Wait 函数前需要先加锁；并且由于Wait函数被唤醒的时候存在虚假唤醒等情况，导致被唤醒后发现，条件依旧不成立，因此需要for来循环等待，直到条件成立
- Cond.Signal(): 只唤醒一个最先 Wait 的 goroutine，可以不用加锁
- Cond.Broadcast(): 唤醒所有Wait的goroutine，可以不用加锁

### Go 由那些方式安全读写共享变量
|方法|并发原语|备注|
|:---:|:---:|:---:|
|不要修改变量|sync.Once|不要去写变量，变量只初始化一次|
|只允许一个goroutine访问变量|channel|不要通过共享变量来通信，通过通信(channel)来共享变量|
|允许多个goroutine访问变量，但是同一时间只允许一个goroutine访问呢|sync.Mutex sync.RWMutex 原子操作|实现锁机制，同时只有一个线程能拿到锁|

### Go 如何排查数据竞争问题

只要有两个以上的 goroutine 并发访问同一变量，且至少其中的一个是写操作的时候就会发生数据竞争；全是读的情况下是不会发生数据竞争的

```go
package main

import "fmt"

func main() {
	i := 0
	go func() {
		i++
	}()
	fmt.Println(i)
}
```
go 命令行有个参数 race 可以帮助检测代码中的数据竞争

go run -race main.go