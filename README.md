# golangfacecontent
## 基础部分

### goLang 的优势

- 部署简单，不依赖其他的库
- 静态类型语言(看运行的时候是编译后运行还是解释器运行)，强类型方便阅读和重构
- 语言层面天生支持高并发，充分利用多核 -> GMP
- 工程化优秀，GoFmt可以统一代码格式
- 强大的标准库， runtime 系统调度机制，高效的GC垃圾回收 -> Go的垃圾回收机制
- 因为是开发k8s的语言，go社区与k8s社区良性互动和发展

## 数据类型
https://github.com/ameamezhou/go-data-structure

这里做了部分常用数据类型详细的记录

这个仓库要重点看，比如slice和map的扩容，函数内的调用都设计它们的底层结构，说得越清楚越好

简单来说
1. bool
2. 数字类型 uint int float32 float64 byte rune
3. 字符串类型
4. 复合类型
    - 数组
    - 切片
    - map
    - 管道
    - 结构体 struct
5. 指针 pointer
6. 接口 interface
7. 函数
8. 方法类型 method (注意和函数进行区分)

其实和关键字一样 属于基础提问，但是可以针对你回答的问题进行深入探究

### 方法和函数的区别
在go语言里面，函数和方法不太一样，有明确的概念区分。在其他语言中，比如java，一般来说函数就是方法，但是在go语言中，函数是指不属于任何结构体、类型的方法，也就是说
函数是没有接收者的，但是方法有接收者

```go

func (t *T) add (a, b int) int {
	return a + b
}

// 其中T是自定义类或者说结构体，不能是基础数据类型 int 等

func add (a, b int) int {
    return a + b
}

```

### 方法接收者和指针接收者的区别
这两者都属于能把函数内的修改带回到本身的一个使用方法，所以要区分好他们之间有什么区别

如果方法的接收者是指针类型，无论调用者是对象还是对象指针，修改的都是对象本身，会影响调用者

如果方法接收者是值类型，无论调用者是对象还是对象指针，修改的都是对象的副本，不影响调用者

我们通常使用指针类型作为方法的接收者的理由：
- 使用指针类型能够修改调用者的值
- 使用指针类型可以避免在每次调用方法的时候复制该值，在值的类型为大型结构体时，这个做法会更高效

### Go函数返回局部变量的指针是否是安全的
一般来说，局部变量会在函数返回后被销毁，因此被返回的引用就成了“无所指”的引用，陈旭会进入未知状态。

但这在Go中是安全的，Go编译器将会对每个局部变量进行逃逸分析.如果发现局部变量的作用域超出该函数，则不会将内存分配在栈上，而是分配在堆上，因为他们不在栈区，所以即使释放函数，其内容也不会受影响

(这里关于堆栈的内容，详细的要看golang的垃圾回收机制和操作系统的堆栈分配区的区别)

```go
package main

import (
   "fmt"
   "net/http"
)

func add(x, y int) *int {
   res := x + y
   return &res
}

func main(){
	fmt.Println(add(1, 2))
}
```
这个例子中，add函数的局部变量 res 发生了逃逸，res作为返回值，在main函数中继续使用，因此res指向的内存不能够分配在栈上，随着函数结束而会回收，因此只能分配在堆上

编译的时候可以用 -gcflags=-m 查看变量逃逸的情况
![img.png](./基础篇/img.png)
我们看到 res escapes to heap 代表res从栈区分配到了堆区，发生了逃逸

### go函数中参数传递到底是值传递还是引用传递
Go语言中所有的传参都是值传递，都是一个副本一个拷贝

参数如果是非引用类型 int string struct 这些，这样就在函数中无法修改原内容数据；如果是引用类型 (指针、map、slice、chan等这些)，这样就可以修改原内容数据

是否可以修改原内容数据，和传值、传引用没有必然的关系，在c++中，传引用肯定是可以修改原内容数据的，但是在Go中虽然只有传值，但是我们也可以修改原内容数据，因为参数是引用类型。

引用类型和引用传递是两个概念
- 值传递: 将实际参数的值传递给形参，形式参数是实际参数的一份拷贝，实际参数和形式参数的内存地址不同。函数内堆形式参数值的内容进行修改，至于是否影响实际参数的值的内容，取决于参数是否是引用类型
- 引用传递: 将实际参数的地址传递给形式参数，函数内堆形式参数内容的修改，将会影响实际参数的值的内容。GO语言中是没有引用传递的，在c++中函数参数的传递方式又引用传递。

**int 类型**
```go
package main

import "fmt"

func main(){
	var i = 1
	fmt.Printf("原内存地址 %p \n", &i)
	modifyInt(i)
	fmt.Printf("改动后值 %v \n", i)
}

func modifyInt(i int){
	fmt.Printf("函数内内存地址 %p \n", &i)
	i = 10
}
```
**指针类型**
```go
package main

import "fmt"

func main(){
   var args = 1
   p := &args
   fmt.Printf("原指针内存地址 %p \n", &p)
   fmt.Printf("原指针变量内存地址 %p \n", p)
   modifyPointer(p)
   fmt.Printf("改动后值 %v \n", *p)
}

func modifyPointer(i *int){
   fmt.Printf("函数内内存地址 %p \n", &i)
   *i = 10
}
```
**Slice 类型**

形式参数和实际参数内存地址一样，不代表是引用类型；下面进行详细说明slice还是值传递，传递的是指针
![img_1.png](./基础篇/img_1.png)

slice 这里的结构体参考golang数据结构那个仓库   后面的map同

**map**

map 形式参数和实际参数内存地址不同，所以其实还是值传递
![img_2.png](./基础篇/img_2.png)
因为这里我们通过make创建的map变量的本质是一个hmap类型的指针，所以函数内堆形参的修改还是会返回原来的内容数据

**channel**

![img_3.png](./基础篇/img_3.png)

因为通过make创建的chan本质也是一个hchan类型的指针，所以堆形参的修改会修改原内容数据

**struct**
形参和实际参数内存地址不一样  是值传递，只要内部的元素不是指针类型的  函数内对形参的修改就不会修改原来的内容数据

## 关键字
### 声明相关
package: 包声明

import: 引入包

var: 变量声明

const: 常量声明

interface: 接口声明

struct: 结构体声明

map、chan: 类型声明

type: 自定义类型声明

### 函数相关
func: 函数定义

return: 从函数返回

### 流程控制
break case continue for fallthrough else if switch goto default: 流程控制

go: 创建goroutine

range: 遍历读取 slice chan map 的数据

select: IO

这个关键字类型的问题属于发散性的面试问题，比如说我是面试官，我问你golang常见的关键字有哪些，你回答了interface那些，就可能会被针对性的提问这个方向，所以要对自己回答的每个关键字要做到心里有数

### defer关键字的实现原理
defer能让我们推迟执行某些函数调用，推迟到当前函数返回前才会实际执行。

defer与panic和recover结合  形成了Go语言风格的异常与捕获机制

使用场景：
defer语句进场被用于处理承兑的操作，比如文件句柄关闭，关闭连接、释放锁等等

优点：方便开发者使用

缺点: 又性能损耗

实现原理： Go1.14中编译器会将defer函数直接插入到函数尾部，无需链表和栈上参数拷贝，性能大幅提升。把defer函数在当前函数内展开并直接调用，这种方式被称为open coded defer

源代码：
```go
func A(i int) {
	defer A1(i, 2*i)
	if (i > 1) {
		defer A2("hello")
    }
    // code
    return
}

func A1 (a, b int) {
    //	
}

func A2(m string) {
	
}

// 编译后  会变成这样

func A(i int){
	// code
    if (i > 1) {
        A2("hello")
    }
    A1(i, 2*i)
}
```

1. 函数退出钱 按照现金后厨的顺序执行defer函数
2. panic 后的defer函数是不会执行的
3. panic没有被recover的时候，抛出的panic 到当前goroutine最上层函数的时候，最上层程序直接异常终止
4. panic 有被recover的时候，当前goroutine最上层函数正常执行

### recover
Recover是在defer中的，它只能捕获自身协程内的异常，不能跨协程捕获，然后实际上的实现原理应该是再函数栈上调用的时候触发panic就会在推出的时候调用，输出panic内容，不因为一个协程挂了就影响main

然后recover并不是所有的错误都能获取到，它只能获取一些panic，更严重的fatal是不能被获取的。比如map是一个非线程安全的map，不能直接进行并发写，会触发fatal，这个是不能被recover捕获的

### new 和 make 的区别
纠正一下，make和new是内置函数，不是关键字

变量初始化一般包括两步，变量声明+变量内存分配。

new和make函数主要是用来分配内存的。

var声明值类型的变量的时候，系统会默认为他分配内存空间，并赋该类型的零值

比如bool，数字，字符串，结构体

如果指针类型或者引用类型的变量，系统不会为它分配内存，默认就是nil。此时如果你想直接使用的话，系统会抛出异常，必须进行内存分配之后，才能使用

new和make两个内置函数主要是用来进行内存空间的分配，有了内存空间，变量才能使用，二者主要有以下两点的区别:
1. 使用场景区别
    - make只能用来分配以及初始化类型为 slice map chan 的数据
    - new可以分配任意类型的数据，并且置0
2. 返回值的区别
    - make 函数原型如下，返回的是slice map chan 本身
    - new 函数原型如下，返回一个指向该类型内存地址的指针
```go
func make(t Type, size ...IntegerType) Type
```
```go
func new(Type) *Type
```

## slice head
slice的具体原理已经在golang数据结构那个仓库有总结了，这里再梳理一边

源码在 src/runtime/slice.go 里面定义了slice的数据结构

```go
type slice struct{
	array   unsafe.Pointer
	len     int
	cap     int
}
```
slice 占用24个字节

- array： 指向底层数组的指针 占用8字节
- len： 切片的长度 占用8字节
- cap： 切片的容量，cap总是大于等于len的，占用8个字节

slice 有四种初始化方式
```go
// var 直接声明
var slice1 []int

// 字面量初始化
slice2 := []int{1, 2, 3, 4}

// make
slice3 := make([]int, 3, 5)

// 从切片或者数组截取
slice4 := arr[1:3]
```

通过一个简单的程序看看slice初始化调用的底层函数
```go
package main

import "fmt"

func main(){
	slice := make([]int, 0)
	slice = append(slice, 1)
	fmt.Println(slice, len(slice), cap(slice))
}

```
通过 `go tool compile -S test.go | grep CALL` 得到汇编代码
![img.png](./part_slice/img.png)

初始化slice调用的是runtime.makeslice, makeslice 函数的主要工作就是计算slice所需的内存大小，然后调用mallocgc进行内存分配

所需内存大小 = 切片中元素大小*切片的容量
```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
	mem, overflow := math.MulUintptr(et.size, uintptr(cap))
	if overflow || mem > maxAlloc || len < 0 || len > cap {
		// NOTE: Produce a 'len out of range' error instead of a
		// 'cap out of range' error when someone does make([]T, bignumber).
		// 'cap out of range' is true too, but since the cap is only being
		// supplied implicitly, saying len is clearer.
		// See golang.org/issue/4085.
		mem, overflow := math.MulUintptr(et.size, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panicmakeslicelen()
		}
		panicmakeslicecap()
	}

	return mallocgc(mem, et, true)
}
```

### array 和 slice的区别
数据结构那个仓库有  再梳理一边

1. 数组长度不同
    - 数组初始化必须指定长度 并且长度就是固定的
    - 切片的唱剫不固定的，可以追加元素，再追加的时候可能使切片的容量增大
2. 函数传参不同
    - 数组是值类型，将一个数组赋值给另一个数组的时候，传递的是一份深拷贝，函数传参操作都会复制整个数组数据，会占用额外的内存，函数内对数组元素值的修改不会修改原数组的内容
    - 切片是引用类型，将一个切片赋值给另一个切片的时候，传递是一份浅拷贝，函数传参操作不会直接拷贝整个切片，只会复制len和cap，底层共用同一个数组，不会占用额外的内存。函数内对数组元素值的修改会修改原数组内容。
3. 计算数组长度方式不同
    - 数组需要变例计算数组长度，时间复杂度为 O(n)
    - 切片底层包含len字段，可以通过len()计算切片长度，时间复杂度为 O(1)

### slice 的深浅拷贝
深拷贝: 拷贝的是数据本身，创造一个新对象，新创建的对象与原对象不共享内存，新创建的对象再内存中开辟一个新的内存地址，新对象值修改的时候不会影响原对象的值

实现深拷贝的方式：
- copy(slice2, slice1)
- 遍历append赋值
```go
package main

import "fmt"

func main(){
	slice1 := []int{1, 2, 3, 4}
	slice2 := make([]int, 5, 5)
	fmt.Printf("slice1: %v, %p \n", slice1, slice1)
	copy(slice2, slice1)
	fmt.Printf("slice2: %v, %p \n", slice2, slice2)
    slice3 := make([]int, 0, 5)
    for _, v := range slice1{
    	slice3 = append(slice3, v)
    }
    fmt.Printf("slice2: %v, %p \n", slice3, slice3)
}
```
浅拷贝：拷贝的是数据地址，只复制指向对象的指针，此时新对象和老对象指向的内存地址是一样的，新对象修改值时老对象也会发生变化

实现浅拷贝的方式:
1. slice2 := slice1

引用类型的变量，默认复制操作就是浅拷贝

```go
package main

import "fmt"

func main(){
   slice1 := []int{1, 2, 3, 4}
   fmt.Printf("slice1: %v, %p \n", slice1, slice1)
   slice2 := slice1
   fmt.Printf("slice2: %v, %p \n", slice2, slice2)
}
```
### slice 的扩容机制
扩容会发生在slice append的时候，当slice的cap不足以容纳新元素的时候就会触发扩容，扩容规则如下:
- 如果新申请的容量比两倍原有的容量大，那么扩容后容量大小为新申请容量
- 如果原有slice长度小于1024，那么每次扩容为原来的两倍
- 如果原有slice长度大于等于1.24，那么每次扩容就为原来的1.25倍
  ![img_1.png](./part_slice/img_1.png)

```go
var inta = [5]int{1, 2, 3, 4, 5}
ints := inta[0:]
fmt.Println(cap(ints))
// 这个时候 使用 cap(ints)查看  会发现 cap==5, 此时如果要append一个元素会怎么样
ints = append(ints, 6)
fmt.Println(cap(ints))
// 这里我们会发现cap变成了10
// 先不去考虑cap的问题, 我们知道array在内存中是一个连续的一段,并且不能扩大;
// 那么当slice需要表示的len超过了array就会重新给slice创建一个新的array, 再将元数据拷贝过去
// 至此就能理解为什么会出现cap变成10的原因了;
// 因为slice是可以扩大的, 如果没append一次就要重新创建数组再copy回来, 那么对于性能的损耗就会比较大
// 所以 Go 对slice的扩容做了优化
```

### slice为什么不是线程安全的
线程安全的定义:
- 多个线程访问一个对象的时候，可以调用这个对象的行为，并且都能获得正确的结果，那么这个对象就是线程安全的
- 若有多个线程同时执行写操作，一般都要考虑线程同步，否则的话就可能影响线程安全。

Go实现线程安全的常用方式
1. 互斥锁
2. 读写锁
3. 原子操作
4. sync.once
5. sync.atomic
6. channel

slice底层结构并没有用加锁等方式，不支持并发读写，所以并不是线程安全的，使用多个goroutine对slice变量进行操作的时候，每次输出的值大概率都不会一样，与预期的不一致；slice在并发执行的过程中不会报错，但是数据会丢失

```go
/*
切片非并发安全
多次执行看到的结果不同
可以考虑用channel不呢神的特性(阻塞)来实现安全的并发读写
 */

package sliceTest

import (
	"sync"
	"testing"
)

func TestSliceConcurrencySafy(t *testing.T){
	a := make([]int, 0)
	var wg sync.WaitGroup
	for i :=0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			a = append(a, i)
			wg.Done()
        }(i)
    }
    wg.Wait()
	t.Log(len(a))
}
```

[toc]
## Go map
### Go map 的实现原理
Go中的map是一个指针，占用8个字节，指向hmap结构体

源码在 src/runtime/map.go 中定义了hmap的数据结构

hmap半酣若干给结构为bmap的数组，每个bmap底层都采用链表结构，bmap通常称其为bucket，也就是go数据结构中桶的那个说法

具体用到的一些数据结构
```go
// Go map 的底层结构体表示
type hmap struct {
    count     int    // map中键值对的个数，使用len()可以获取 
	flags     uint8  // 这事一个状态标记位，标记是否处于正在写入
	B         uint8  // 哈希桶的数量的log2，比如有8个桶，那么B=3
	noverflow uint16 // 溢出桶的数量
	hash0     uint32 // 哈希种子

	buckets    unsafe.Pointer // 指向哈希桶数组的指针，数量为 2^B 
	oldbuckets unsafe.Pointer // 扩容时指向旧桶的指针，当扩容时不为nil 
	nevacuate  uintptr        // 扩容进度，小鱼此处地址的buckets表示已经迁移完成了

	extra *mapextra  // 可选字段  存储溢出桶，这个字段是为了优化GC扫描而设计的
}

const (
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits     // 桶数量 1 << 3 = 8
)

// Go map 的一个哈希桶，一个桶最多存放8个键值对
type bmap struct {
    // tophash存放了哈希值的最高字节
	tophash [bucketCnt]uint8
	// 用于实现快速定位key的位置，在实现过程中会使用key的哈希值的高八位作为tophash存放在tophash字段中
	// tophash字段不仅存储key哈希值的高八位，还会存储一些状态来表明当前桶的状态
    // 特殊标记相关在下面tophash有明确说明
    // 在这里有几个其它的字段没有显示出来，因为k-v的数量类型是不确定的，编译的时候才会确定
    // keys: 是一个数组，大小为bucketCnt=8，存放Key
    // elems: 是一个数组，大小为bucketCnt=8，存放Value
    // 你可能会想到为什么不用空接口，空接口可以保存任意类型。但是空接口底层也是个结构体，中间隔了一层。因此在这里没有使用空接口。
    // 注意：之所以将所有key存放在一个数组，将value存放在一个数组，而不是键值对的形式，是为了消除例如map[int64]所需的填充整数8（内存对齐）
    keys    [bucketCnt]keyType
	values  [bucketCnt]valueType
    // overflow: 是一个指针，指向溢出桶，当该桶不够用时，就会使用溢出桶
    overflow uintptr
}
```
**图像如下**
![img.png](./part_map/img.png)
这里还没有画出溢出桶，找个图
![img_1.png](./part_map/img_1.png)
这里绿色部分就是溢出桶

```go
//可能的tophash值。我们保留了一些特殊标记的可能性。
//每个存储桶（包括其溢出存储桶，如果有的话）将有其全部或全部
//vacuum*状态中的条目（在evacuate（）方法期间除外，该方法只发生
//在映射写入期间并且因此在该时间期间没有其他人能够观察到映射）。
emptyOne       = 1 // this cell is empty
evacuatedX     = 2 // key/elem is valid.  Entry has been evacuated to first half of larger table.
evacuatedY     = 3 // same as above, but evacuated to second half of larger table.
evacuatedEmpty = 4 // cell is empty, bucket is evacuated.
minTopHash     = 5 // minimum tophash for a normal filled cell.

// tophash calculates the tophash value for hash.
func tophash(hash uintptr) uint8 {
    top := uint8(hash >> (goarch.PtrSize*8 - 8))
    if top < minTopHash {
        top += minTopHash
    }
    return top
}

这里我们可以看到，为了防止高八位和这些状态值相等，都自动加上了minTopHash这些值
```

**mapextra结构体**
```go
// mapextra holds fields that are not present on all maps.
type mapextra struct {
    //如果key和elem都不包含指针并且是内联的，那么我们标记bucket
    //类型为不包含指针。这样可以避免扫描此类地图。
    //但是，bmap.overflow是一个指针。为了保持水桶溢出
    //活着时，我们将指向所有溢出存储桶的指针存储在hmap.extra.overflow和hmap.extra.oldoverflow中。
    //只有当key和elem不包含指针时，才会使用overflow和oldoverflow。
    //overflow包含用于hmapbucket的溢出bucket。
    //oldoverflow包含hmap.oldbuckets的溢出bucket。
    //间接寻址允许在hiter中存储指向切片的指针。	overflow    *[]*bmap
	overflow    *[]*bmap
	oldoverflow *[]*bmap

	// nextOverflow保存一个指向空闲溢出存储桶的指针。
	nextOverflow *bmap
}
```
当map的key和value都不是指针类型对象的时候，bmap将完全不包含指针，那么gc的时候就不用扫描bmap，bmap指向一处同的字段overflow是uintptr类型，为了防止这些overflow
桶被gc掉，所以需要mapextra.overflow将它保存起来。如果bmap的overflow是*bmap类型，那么gc扫描的是一个个拉链表，效率明显不如直接扫描一段内存

**总结**

注意这里上面提到的key和value都是各自存放在一起的，并不是key-value/key-value 这种存储形式，当key和value类型不同的时候二者占用的字节大小不一样，这样可能会因为
考虑内存对齐而造成内存空间浪费，所以go采用key和value分开存储的设计，这样更节省内存空间
![img_2.png](./part_map/img_2.png)

这一段是map里面必须要弄懂的，后面扩容相关规则参考这篇博客 说得还挺清楚的
https://blog.csdn.net/Peerless__/article/details/125458742

### Go map的遍历为什么是无序的
使用range多次遍历map的时候输出的key和value顺序有可能不同，这事Go语言的设计者们有意为之，旨在告诉开发者们，Go底层实现并不保证map遍历顺序稳定，请打架
不要依赖range遍历结果顺序

主要原因有两点：
- map在遍历的时候并不是从固定的0号bucket开始遍历的，每次遍历都会从一个随机值序号的bucket，在从其中随机的cell开始遍历
- map遍历时，是按序遍历bucket，同时按需遍历bucket中和其他overflow bucket中的cell。但是map在扩容后会发生key的搬迁，这造成原来落在一个bucket中的key，搬迁后，有可能落到其他bucket中了，从这个角度看遍历map的结果就不可能是按照原来的顺序了

map本身是无序的，且遍历的时候顺序还会被随机化，如果想顺序遍历map，需要对map key 先排序，再按照key的顺序遍历map。

### 为什么map不是线程安全的
map默认是并发不安全的，同时对map进行并发读写，程序会出现panic

Go官方在经过长时间讨论后认为map更适配典型使用场景，不需要从多个goroutine中进行安全访问，而不是为了小部分情况(并发访问)，导致大部分程序付出枷锁的代价(性能)，决定了不支持

如果两个协程同时读写，会出现致命错误：fatal error: concurrent map writes

**注意！** 这个fatal是不能被recover进行异常捕获的

如果想要实现map的线程安全
- 方法1：使用读写锁 --- map + sync.RWMutex

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var lock sync.RWMutex
	s := make(map[int]int)
	for i := 0; i < 100; i++ {
		go func(i int) {
			lock.Lock()
			s[i] = i
			lock.Unlock()
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			lock.RLock()
			fmt.Printf("map 元素 %v    %v \n", i, s[i])
			lock.RUnlock()
		}(i)
	}
	time.Sleep(1 * time.Second)
}
```

- 方法2 使用 Go提供的 sync.map
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var m sync.Map
	for i := 0; i < 100; i++ {
		go func(i int) {
			m.Store(i, i)
		}(i)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			v, ok := m.Load(i)
			fmt.Printf("load: %v, %v", v, ok)
		}(i)
	}
	time.Sleep(1 * time.Second)
}
```

### Go map 如何查找
Go语言中读取map有两种语法：带comma和不带comma，当要查询的key不在map里面，带comma的用法会返回一个bool型的变量提示key是否在map中，而不殆comma的语句则会返回一个value类型的零值。
如果value是int就会返回0，如果value是string类型就会返回空字符串。

```go
// 不带 comma
value := a["name"]
fmt.Printf("value %s", value)

// 带 comma
value, ok := a["name"]
```

map的查找可以通过生成会变吗可以知道，根据key的不同类型/返回参数，编译器会将查找函数用具体的函数替换，优化效率
![img_3.png](./part_map/img_3.png)

**查找流程**
![img_4.png](./part_map/img_4.png)

1. 写保护检测

函数首先会检查map的标志位flags，如果flags的写标志位此时被置为1了，说明有其他的协程正在进行写操作，进而导致程序panic，这也说明了map不是线程安全的

![img_5.png](./part_map/img_5.png)

2. 计算hash值
```go
hash := t.hasher(key, uintptr(h.hash0))
```
key 经过哈希函数计算之后，得到的哈希值如下（主流64位机下共六十四个bit位） 不同类型的key会有不同搞得hash函数：

10010111|000011110110110010001111001010100010010110010101010|01010

3. 找到hash值对应的bucket

bucket定位：哈希值的低B个bit位，用来定位key锁存放的bucket

如果当前正在扩容中，并且定位到的旧的bucket数据还未完成迁移，就会使用就的bucket(扩容前的bucket)

```go
// 计算hash值
hash := t.hasher(key, uintptr(h.hash0))
// 桶的个数n-1，即 1 << B-1, B=5时，则有0-31号桶
m := bucketMask(h.B)
// 计算hash值对应的bucket
// t.bucketsize 为一个bmap的大小，通过对哈希值和桶个数取模得到桶的编号，通过对桶编号和buckets其实地址进行运算，获取哈希值对应的bucket
b := (*bmap)(add(h.buckets, (hash&m)*uintptr(t.bucketsize)))
// 是否在扩容
if c := h.oldbuckets; c != nil {
	// 桶的个数已经发生增长，则就bucket的桶个数为当前桶个数的一半
    if !h.sameSizeGrow() {
        // There used to be half as many buckets; mask down one more power of two.
        m >>= 1
    }
    // 计算哈希值对应的旧的bucket
    oldb := (*bmap)(add(c, (hash&m)*uintptr(t.bucketsize)))
    // 如果就的bucket数据还没有完成迁移，则使用旧的bucket查找
    if !evacuated(oldb) {
        b = oldb
    }
}
top := tophash(hash)
```

4. 遍历bucket查找

tophash值定位：哈希值的高八个bit位，用来快速判断key是否已经存在当前bucket中，如果不在的画则需要取bucket的overflow中查找

用步骤2中的hash值得到高八个bit位，也就是10010111，转化为10进制也就是151
```go
top := tophash(hash)

// tophash calculates the tophash value for hash.
func tophash(hash uintptr) uint8 {
    top := uint8(hash >> (goarch.PtrSize*8 - 8))
    if top < minTopHash {
        top += minTopHash
    }
    return top
}
```
上面函数中的hash是六十四位的，但是sys.PtrSize 的值是8，所以 `top := uint8(hash >> (goarch.PtrSize*8 - 8))` 等同于 `top = uint8(hash >> 56)`
最后top取出来的值就是hash的高八位值

在bucket以及bucket的overflow中寻找tophash值（HOB hash）为151*的曹魏，即key所在的位置，如果找到空槽或者2号槽位，这样整个查找过程就结束了，其中空槽为代表没找到
```go
bucketloop:
	for ; b != nil; b = b.overflow(t) {
		for i := uintptr(0); i < bucketCnt; i++ {
			if b.tophash[i] != top {
				// 未使用的槽位，插入
				if b.tophash[i] == emptyRest {
					break bucketloop
				}
				continue
			}
			// 找到tophash值对应的key
			k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))
			if t.indirectkey() {
				k = *((*unsafe.Pointer)(k))
			}
			if t.key.equal(key, k) {
				e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
				if t.indirectelem() {
					e = *((*unsafe.Pointer)(e))
				}
				return e
			}
		}
	}
```
（这里顶格写的是标签，然后break + 标签 是跳出整个标签，相关的关键字用法还有goto 是跳转到标签段执行，这里可以取搜索下相关的资料）

5. 返回keyh对应的指针

如果上面的步骤找到了key对应的槽位下标i，我们再来看如何获取到key和value的

```go
dataOffset = unsafe.Offsetof(struct {
    b bmap
    v int64
}{}.v)

bucketCntBits = 3
bucketCnt     = 1 << bucketCntBits

// key 定位公式
k := add(unsafe.Pointer(b), dataOffset+i*uintptr(t.keysize))

// value elem 定位公式
e := add(unsafe.Pointer(b), dataOffset+bucketCnt*uintptr(t.keysize)+i*uintptr(t.elemsize))
```
bucket里面keys的起始地址就是unsafe.Pointer(b)+dataOffset

第i个下标key的地址就要在此基础上跨过i个key的大小；

而我们还知道 value 的地址是在所有的key之后，因此第i个下标的value地址还要加上所有key的偏移

### Go map 解决冲突的方式
比较常见的解决hash冲突的方法有链地址发和开放寻址法

**链地址法：** 当哈希冲突发生的时候，创建新的单元，并将新单元添加到冲突单元所在链表的尾部

**开放寻址法：** 当哈希冲突发生的时候，从发生冲突的那个单元起，按照一定的次序，从哈希表中寻找一个空闲的单元，然后把发生冲突的元素存入到该单元。开放寻址发需要的表长度要大于等于所需要存放的元素数量

开放寻址法有多种方式：线性探测法，平方探测法，随机探测法和双重哈希法。这里以线性探测法来说明

**线性探测法**

设 Hash(key) 表示关键字 key 的哈希值，表示哈希表的槽位数（哈希表大小）

线性探测法可以表示为：

如果 `Hash(x) % M` 已经有数据，则尝试 `(Hash(x) + 1) % M`;

如果 `Hash(x + 1) % M` 已经有数据，则尝试 `(Hash(x) + 2) % M`;

以此类推

**两种方法比较**

对于链地址法，基于数组+链表进行存储，链表结点可以在需要的时候再创建，不必像开放地址法那样先申请号足够的内存，因此链地址法对于内存的利用率会比开放地址法的利用率高。
链地址法对装载因子的容忍度会更高，并且适合存储大对象、大数据量的哈希表。而且相较于开放寻址法则，它更加灵活，支持更多的优化策略，比如可以采用红黑树来代替链表。但是链地址法需要额外的空间来存储指针

对于开放地址法，它只有数组一种数据结构完成存储，继承了数组的优点，比如对cpu缓存友好，易于序列号操作，但是它对于内存的利用率不如链地址法，并且发生冲突的时候代价更高。当数据量明确，装载因子小的时候更适合用开放寻址法

**总结**

在发生哈希冲突的时候，python中的dict采用的开放寻址发，java的hashmap采用的是链地址法，具体就是插入key到map中，当key定位的桶充满八个元素后(这里的单元就是桶，不是元素)，将会创建出一个溢出桶，并且将溢出桶插入当前桶所在链表尾部

```go
	if inserti == nil {
		// The current bucket and all the overflow buckets connected to it are full, allocate a new one.
		newb := h.newoverflow(t, b)
		// 创建要给新的溢出桶
		inserti = &newb.tophash[0]
		insertk = add(unsafe.Pointer(newb), dataOffset)
		elem = add(insertk, bucketCnt*uintptr(t.keysize))
	}
```

### Go 的负载因子为什么是6.5
在Go的数据结构那一章里面提到了Go map的负载因子

负载因子就是用于衡量当前哈希表中空间占用率的核心指标，也就是每个bucket桶存储的平均元素个数。

负载因子 = 哈希表存储的元素个数/桶个数

另外负载因子与扩容、迁移等重新散列(rehash)行为有直接关系：
- 在程序运行到时候会不断地进行插入、删除等，会导致bucket不均，内存利用率低，需要迁移
- 在程序运行的时候会出现负载因子过大，需要做扩容，解决bucket过大的问题。

负载因子是哈希表中的一个重要指标，在各种版本的哈希表视线中都有类似的东西，主要目的是为了平衡buckets的存储空间大小和查找元素时的性能高低。

在接触各种哈希表的时候都可以关注一下，做不同的对比，看看各家的考量。

```go
func overLoadFactor(count int, B uint8) bool {
    return count > bucketCnt && uintptr(count) > loadFactorNum*(bucketShift(B)/loadFactorDen)
}
const loadFactorNum = 13
// 扩容规则的意思是：如果map中键值对的数量 count> 8，也就是说，至少要能装满一个bmap；
// 且 count > 13*桶的数量/2，也就是说 count/bucketCount >6.5；两个条件都满足才会允许扩容；
```

**为什么是6.5？**

为什么Go语言中的哈希表的负载因子是6.5，为什么不是8也不是1，这里面有可靠的数据支撑吗？

**测试报告**

实际上这是Go官方经过认真测试得出的数字，在官方报告中一共包含4个关键指标

loadFactor、%overflow、bytes/entry、hitprobe、missprobe

- loadFactor: 负载因子，也有叫装载因子
- %overflow: 溢出率，有溢出bucket的百分比
- bytes/entry: 平均每对key value的开销字节数
- hitprobe: 查找一个存在的key时，需要查找的平均个数
- missprobe: 查找一个不存在的key的时候，需要查找的平均个数

**选择数值**

Go官方发现：装在因子越大，填入的元素越多，空间利用率就越高，但是发生哈希冲突的可能性就变大。反之，转载因子越小，填入的元素越少，冲突发生的几率就越小，但是空间浪费就会变得更多，而且会提高扩容操作的次数

根据测试结果，Go官方选取了一个相对适中的值，把Go中的map的负载因子编码为6.5，这就是6.5的理由

这意味着在Go语言中，当map存储的元素个数大于或者等于 6.5 * 桶个数的时候，就会触发扩容行为

### map 如何扩容
**扩容规则**

1. 条件1：超过负载。 map元素个数 > 6.5 * 桶个数
```go
func overLoadFactor(count int, B uint8) bool {
    return count > bucketCnt && uintptr(count) > loadFactorNum*(bucketShift(B)/loadFactorDen)
}
const loadFactorNum = 13
// 扩容规则的意思是：如果map中键值对的数量 count> 8，也就是说，至少要能装满一个bmap；
// 且 count > 13*桶的数量/2，也就是说 count/bucketCount >6.5；两个条件都满足才会允许扩容；
```
2. 条件2： 溢出桶太多。

当桶总数 < 2^15 时， 如果溢出桶总数 >= 桶总数，则会认为溢出桶过多。

当桶总数 >= 2^15 时， 直接与 2^15 比较，当溢出桶总数 >= 2^15，则会认为溢出桶太多
```go
// tooManyOverflowBuckets reports whether noverflow buckets is too many for a map with 1<<B buckets.
// Note that most of these overflow buckets must be in sparse use;
// if use was dense, then we'd have already triggered regular map growth.
func tooManyOverflowBuckets(noverflow uint16, B uint8) bool {
	// If the threshold is too low, we do extraneous work.
	// If the threshold is too high, maps that grow and shrink can hold on to lots of unused memory.
	// "too many" means (approximately) as many overflow buckets as regular buckets.
	// See incrnoverflow for more details.
	if B > 15 {
		B = 15
	}
	// The compiler doesn't see here that B < 16; mask B to generate shorter shift code.
	return noverflow >= uint16(1)<<(B&15)
}
```
对于条件2，其实算是对条件1的补充。因为在负载因子较小的情况下，有可能map的查找效率也低，而第一点识别不出来这种情况。

表面现象就是负载因子比较小，即map里元素总数少，但是桶的数量多(真是分配的桶数量多，包括大量的溢出桶)。比如不断的增删，这样会造成overflow的bucket数量增多
，但是负载因子又不高，达不到第一点的临界值，就不能触发扩容机制来环节这种情况，这样会造成桶的使用率不高，值存储的比较系数，查找插入效率会变得比较低，因此有了第二扩容条件

**扩容机制:**

- 双倍扩容：针对条件1，新建一个buckets数组，新的buckets大小是原来的两倍，然后旧的buckets数据搬迁到新的buckets。这种方法称为双倍扩容
- 等量扩容：针对条件2，并不扩大容量，buckets数量维持不变，重新做一遍类似双倍扩容的搬迁动作，把松散的键值对重新排列一次，是的同一个bucket中的key排列得更加紧密，更节约空间，提高bucket的利用率。

**扩容函数:**

在golang数据结构中提到的hashGrow并没有实现真正的”搬迁“，它只是分配好了新的buckets，并将老的buckets挂到了oldbuckets字段上。真正搬迁buckets动作在growWork()函数中，
而调用 growWork() 函数的动作是在 mapassign 和 mapdelete 函数中，也就是插入或者修改、删除key的时候都会尝试进行搬迁buckets的工作，先检查oldbuckets是否搬迁完，具体来说就是检查oldbuckets是否为nil

```go
func hashGrow(t *maptype, h *hmap) {
    // bigger为需要扩充的数量
    bigger := uint8(1)
    // 判断是否满足扩容条件
    if !overLoadFactor(h.count+1, h.B) {
        // 不满足bigger为0
        bigger = 0
        h.flags |= sameSizeGrow
    }
    // oldbuckets和 按照修改后的数组创建 newbuckets
    // 记录老的 buckets
    oldbuckets := h.buckets
    // 申请新的buckets
    newbuckets, nextOverflow := makeBucketArray(t, h.B+bigger, nil)
    // 注意 &^ 运算符，这块代码的逻辑是转移标志位
    flags := h.flags &^ (iterator | oldIterator)
    if h.flags&iterator != 0 {
        flags |= oldIterator
    }
    // 修改h的buckets数量，也就是翻倍，例如原来B=2，数量为 1<<2 == 4，1<<(2+1) == 8；
    // 修改flag，把oldbuckets、newbuckets修改，将rehash进度置为0，将溢出桶的数量置为0
    h.B += bigger
    h.flags = flags
    h.oldbuckets = oldbuckets
    h.buckets = newbuckets
    // 搬迁进度
    h.nevacuate = 0
    h.noverflow = 0
    // 修改 extra字段中的 oldoverflow 和 overflow 
    if h.extra != nil && h.extra.overflow != nil {
        // Promote current overflow buckets to the old generation.
        if h.extra.oldoverflow != nil {
            throw("oldoverflow is not nil")
        }
        h.extra.oldoverflow = h.extra.overflow
        h.extra.overflow = nil
    }
    if nextOverflow != nil {
        if h.extra == nil {
            h.extra = new(mapextra)
        }
        h.extra.nextOverflow = nextOverflow
    }
}
```

由于map扩容需要将原有的key/value 重新搬迁到新的内存地址，如果map存储了数以亿记的key-value，一次性搬迁将会造成比较大的时延，因此 Go map的扩容采取了
一种被称为 **渐进式** 的方式，原有的key并不会一次性搬迁完毕，每次搬迁只会搬迁两个bucket

```go
func growWork(t *maptype, h *hmap, bucket uintptr) {
	// make sure we evacuate the oldbucket corresponding
	// to the bucket we're about to use
	evacuate(t, h, bucket&h.oldbucketmask())

	// evacuate one more oldbucket to make progress on growing
	if h.growing() {
		evacuate(t, h, h.nevacuate)
	}
}
```

### map 和 sync.Map 性能比较

**sync.Map**
```go
type Map struct {
	mu Mutex
	read atomic.Value // readOnly
	dirty map[any]*entry
	misses int
}
```
对比原始map：和原始map+RWLock的实现并发的方式相比，减少了加锁对性能的影响，它做了一些优化：可以无锁访问read map，而且会优先操作 read map，倘若只操作read map 就可以满足要求，
那就不用去操作 write map(dirty)，所以在某些特定场景它发生锁竞争的频率会远远小于 map + RWLock 的实现方式

优点：适合多读写少的场景

缺点：写多的场景，会导致read map 缓存失效，需要加锁，冲突变多，性能急剧下降

## channel
### channel 的底层实现原理

**概念:**  
Go中的channel是一个队列，遵循先进先出的原则，分则协程之间的通信（Go语言体长不要通过共享内存来通信，而要通过通信来实现内存共享，CSP communicating sequential process 并发模型就是通过goroutine 和 channel 来实现的）

**使用场景:**
- 停止信号监听
- 定时任务
- 生产方和消费方解耦
- 控制并发数

**底层数据结构:**  
通过var声明或者make函数创建的channel变量是一个存储在 **函数栈** 上的 **指针**，占用8个字节，指向堆上的hchan结构体

源码中 src/runtime/chan.go 定义了 hchan 的数据结构:

![img.png](img.png)

```go
type hchan struct {
	/*
	channel 分为无缓冲和有缓冲两种
	对于有缓冲的channel存储数据，使用了 ring buffer 环形缓冲区 来写入数据，本质是循环数组
	为啥是循环数组？普通数组不行吗？ 普通数组容量固定更适合指定的空间，弹出元素的时候普通数组需要廍向前移动。
	当下标超过数组容量后会回到第一个位置所以需要两个字段记录当前读写的下表位置
	*/
	qcount   uint           // total data in the queue 循环数组的元素数量
	dataqsiz uint           // size of the circular queue 循环数组的胀肚
	buf      unsafe.Pointer // 指向底层的循环数组指针  也就是ring buffer
	elemsize uint16         // 元素大小
	
	closed   uint32 // channel 是否关闭的标志
	elemtype *_type // element type 元素类型
	
	// 这里就可以看成循环指针
	sendx    uint   // send index 下一次写下标的位置
	recvx    uint   // receive index 下一次读下标的位置
	
	// 尝试向channel读或者向channel写入数据而被阻塞的goroutine
	recvq    waitq  // list of recv waiters 读等待队列
	sendq    waitq  // list of send waiters 写等待队列

	// lock protects all fields in hchan, as well as several
	// 保护所有hchan的部分  甚至包括 sudog 被组织的几个字段
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	/*
	在持有这个锁的时候不要随便修改G的状态
	*/
	lock mutex // 互斥锁  保证读写channel的时候不发生并发竞争
}
```

等待队列：

`waitq`包含一个头节点一个尾结点，是个双向链表

每个结点都是一个sudog结构体变量，记录哪个协程在等待，等待的是哪个channel，等待发送、接收的数据在哪

```go
type waitq struct {
	first *sudog
	last  *sudog
}


type sudog struct {
// The following fields are protected by the hchan.lock of the
// channel this sudog is blocking on. shrinkstack depends on
// this for sudogs involved in channel ops.
g *g

next *sudog
prev *sudog
elem unsafe.Pointer // data element (may point to stack)
c        *hchan // channel
...
}
```

操作：
- **创建**

使用`make(chan T, cap)` 来创建channel，make语法在编译的时候会转化为 `makechan64` 和 `makechan`

```go
func makechan64(t *chantype, size int64) *hchan {
	if int64(int(size)) != size {
		panic(plainError("makechan: size out of range"))
	}

	return makechan(t, int(size))
}

func makechan(t *chantype, size int) *hchan { ... }
```

创建channel有两种，一种是带缓冲的channel，一种是不带缓冲的channel

```go
// buffer
ch := make(chan int, 2)
// no buffer
ch := make(chan int)
```

创建的时候会做一些检查:
1. 元素大小不超过64k
2. 元素对齐大小不能超过 maxAlign 也就是8字节
3. 计算出来的内存是否超过限制

创建时的策略:
1. 如果是无缓冲的 channel，会直接给 hchan 分配内存
2. 如果是有缓冲的 channel，并且元素不包含指针，那么会为 hchan 和底层数组分配一段连续的地址
3. 如果是有缓冲的 channel，并且元素包含指针，那么会为 hchan 和底层数组分别分配地址
```go
// makechan
var c *hchan
switch {
case mem == 0:
    // Queue or element size is zero.
    c = (*hchan)(mallocgc(hchanSize, nil, true))
    // Race detector uses this location for synchronization.
    c.buf = c.raceaddr()
case elem.ptrdata == 0:
    // Elements do not contain pointers.
    // Allocate hchan and buf in one call.
    c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
    c.buf = add(unsafe.Pointer(c), hchanSize)
default:
    // Elements contain pointers.
    c = new(hchan)
    c.buf = mallocgc(mem, elem, true)
}
```
---
- **发送**

发送操作，编译时会转换为 `runtime.chansend` 函数

```go
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool
```

**阻塞式**

调用chansend函数 并且block = true

```go
ch <- 10
```

**非阻塞式**

调用 chansend 函数，并且block=false

```go
select {
    case ch <- 10:
    	...
    default
    	
}
```

向 channel中发送数据时大概分为两大块：检查和数据发送，数据发送流程如下：
- 如果channel的读等待队列存在接收者goroutine
   - 将数据直接发送给第一个等待的 goroutine，唤醒接收的goroutine
- 如果channel的读等待队列不存在接收者goroutine
   - 如果循环数组 buf 未满，那么将会把数据发送到循环数组buf队尾
   - 如果循环数组 buf 已满，这个时候就会走阻塞发送的流程，将当前goroutine加入写等待队列，并挂起等待唤醒

---

- **接收**

发送操作，编译时转换为 `runtime.chanrecv` 函数

```go
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool)
```
**阻塞式**

调用 chanrecv 函数，并且 block = true
```go
<- ch

v := <-ch

v, ok := <-ch

// 当channel 关闭到时候，for循环会自动退出，无需主动监测channel是否关闭，可以防止读取已经关闭的channel，造成督导数据为通道所存储数据类型的零值
for i := range ch {
	fmt.Println(i)
}
```    

**非阻塞式**

调用 chanrecv 函数，并且 block = false
```go
select {
    case <- ch:
    	...
    default
    	
}
```

向 channel 中接收数据时大概分为两大块，检查和数据发送，而数据接收流程如下：
- 如果是channel的写等待队列存在发送者goroutine
   - 如果是无缓冲channel，直接从第一个发送者 goroutine 那里把数据拷贝给接收变量，唤醒发送的 goroutine
   - 如果是有缓冲channel(已满)，将循环数组buf的队首元素拷贝给接收变量，将第一个发送者goroutine的数据拷贝到buf循环数组队尾，唤醒发送的 goroutine
- 如果 channel 的写等待队列不存在发送者 goroutine
   - 如果循环数组buf非空，将循环数组 buf 的队首元素拷贝给接收变量
   - 如果循环数组buf非空，这个时候就会走阻塞接收的流程，将当前 goroutine 加入读等待队列，病挂起等待唤醒

---

- **关闭**

关闭操作，调用 close 函数，编译时转换为 `runtime.closechan` 函数
```go
close(ch)

func closechan(c *hchan) 
```

```go
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
```

---

- **总结**
   - 用来保存goroutine之间传递数据的循环数组: buf
   - 用来记录循环数组当前发送或者接收数据的下标值： sendx 和 recvx
   - 用来保存向该chan发送和从该chan接收数据被阻塞的goroutine队列：sendq和recvq
   - 保证channel写入和读取数据时线程安全的锁: lock

### channel 的特点

channel 的类型：无缓冲、有缓冲

channel 有3种模式：写操作模式(单向通道)、读操作模式(单项通道)、读写操作模式(双向通道)

```go
// 写操作模式
make(chan <-int)
// 读操作模式
make(<- chan int)
// 读写操作模式
make(chan int)
```

channel 有三种状态：未初始化、正常、关闭
![img_1.png](./part_channel/img_1.png)

**注意点**
1. 一个channel不能多次关闭，会导致panic
2. 如果多个 goroutine 都监听头一个 channel，那么channel上的数据都可能随机被某一个goroutine取走进行消费
3. 如果多个 goroutine 都监听同一个 channel，如果这个channel被关闭，则所有 goroutine 都能接收到退出信号

### Go channel 有无缓冲的区别

无缓冲：就类似于这个东西我递给你，你不接我就一直举着手，知道你拿走了，我才会收走

有缓冲：只要你桌子上有空余的地方，我就直接放到你空余的地方就好了，除非你的桌子堆满了，我就要等到你空出一个位置之后我放下才会走

|     | 无 缓 冲 | 有 缓 冲 |
|:---:|:---:|:---:|
|创建方式|make(chan T)|make(chan T, size)|
|发送阻塞|数据接收前发送阻塞|缓冲满的时候发送阻塞|
|接收阻塞|数据发送前接收阻塞|缓冲空的时候接收阻塞|

```go
package main

import (
	"fmt"
	"time"
)
// 非缓冲 channel
func loop(ch chan int){
	for {
		select {
		case i := <- ch:
			fmt.Println()
        }
    }
}

func main(){
	ch := make(chan int)
	ch <- 1
	go loop(ch)
	time.Sleep(1 * time.Second)
}
```
这里回报错 `fatal error` 这是因为 ch<-1 发送了，但是没有接收者，所以出现了阻塞

不过这里我们可以把 ch <- 1 放到 go loop 下面，也能够正常执行


如果希望能正常发送和接受，那我们要做一个缓冲 channel

这样程序也能正常运行，这里不做demo了  就把上面的改成缓冲channel 然后多塞几个进去就好了

### channel 为什么是线程安全的

不同协程通过 channel 进行通信，本身的使用场景就是多线程，为了保证数据的一致性，必须实现线程安全

channel的底层实现中，hchan结构体中就采用了 mutex 锁来保证读数据读写安全，在对循环数组buf中数据进行入队和出队操作的时候，必须先获取互斥锁才能操作channel

### channel 如何控制 goroutine 并发执行顺序

多个 goroutine 并发执行的时候，每一个 goroutine 强盗处理器的时间点不一致，goroutine 的执行本身并不能保证顺序，即代码中险些的 goroutine 并不能保证限制性

思路：用channel进行通知，用channel去传递信息，从而控制并发执行顺序

```go
var wg sync.WaitGroup

func main(){
	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)
	ch3 := make(chan struct{}, 1)
	ch1 <- struct{}{}
	wg.Add(3)
	start := time.Now().Unix()
	go outPut("goroutine1", ch1, ch2)
	go outPut("goroutine2", ch2, ch3)
	go outPut("goroutine3", ch3, ch1)
	wg.Wait()
	end := time.Now().Unix()
	fmt.Printf("duration: %d \n", end - start)
}

func outPut(s string, inch, outch chan struct{}){
	time.Sleep(1 * time.Second)
	select {
	case <- inch:
		fmt.Printf("%s \n", s)
		outch <- struct{}{}
	}
	wg.Done()
}
```

### channel共享内存的优劣

无论是通过共享内存来通信还是通过通信来共或内存，最终我们应用程序都是读取的内存当中的数据，只是前者是直接读取内存的数据，而后者是通过发送消息的方式来
进行同步。而通过发送消息来同步的这种方式常见的就是 Go 采用的 CSP(Ccommunication SequentialProcess)模型以及 Eang 采用的 Actor 模型，这两种方式都是通过
通信来共享内存。

![img_2.png](./part_channel/img_2.png)

大部分的语言采用的都是第一种方式直接去操作内存，然后通过互斥锁，CAS等操作来保证并发安全。Go引入了Channel和 Goroutine 实现 CSP模型将生产者和消费者进行了解耦，Channel其实和消息队列很相似。而Actor 模型和 CSP模型都是通过发送消息来共享内存，但是它们之间最大的区别就是 Actor 模型当中并没有一个独立的Channel组件，而是 Actor与 Actor 之间直接进行消息的发送与接收，每个 Actor 都有一个本地的"信箱"消息都会先发送到这个"信箱当中"。

- **优点**
   - 使用channel可以帮助我们解耦生产者喝消费者，可以降低并发当中的耦合

- **缺点**
   - 容易出现死锁

### channel 死锁情况
- **死锁**
   - 单个协程永久阻塞
   - 两个或者两个以上的协程的执行过程中，由于竞争资源或由于彼此通信而造成的一种阻塞的现象

---

- **channel死锁场景**
   - 非缓存 channel 只写不读
   - 非缓存 channel 读在写后面
   - 缓存 channel 写入超过缓冲区的数量
   - 空读
   - 多个协程互相等待

1. 非缓存channel只读不写
```go
func deadlock1(){
	ch := make(chan int)
	ch <- 3 // 这里会一直阻塞
}
```

2. 非缓存channel读在写后
```go
func deadlock2(){
    ch := make(chan int)
    ch <- 3 // 这里会一直阻塞
    num := <-ch
    fmt.Println(num)
}

func deadlock2(){
    ch := make(chan int)
    ch <- 100
    go func(){
  	    num := <-ch
  	    fmt.Println(num)
    }
    time.Sleep(time.Second)
}
```
3. 缓存 channel 写入超过缓冲区数量
```go
func deadlock3(){
	ch := make(chan int, 3)
	ch <- 1
    ch <- 2
    ch <- 3
    ch <- 4 // 这里会一直阻塞
}
```
4. 空读
```go
func deadlock4(){
    ch := make(chan int)
    num := <-ch
    fmt.Println(num)
}
```
5. 多个协程相互等待
```go
func deadlock5(){
    ch1 := make(chan int)
    ch2 := make(chan int)
    go func(){
    	for {
    	    select {
    	    case num := <- ch1:
    	    	fmt.Println(num)
    	    	ch2 <- 100
            }       	
        }   
    }()
    go func(){
        for {
            select {
                case num := <- ch2:
                fmt.Println(num)
                ch1 <- 100
            }       	
        }   
    }()
}
```

### 空 chan 和 关闭的 chan 进行读写会怎么样
1. 空chan
   - 读会读到该chan类型的零值
   - 写会直接写到chan中
2. 关闭的chan
   - 读已经关闭的chan能一直读到东西，但是读到的内容根据通道内关闭前是否有元素而不同，如果有元素就继续读剩下的元素，没有元素就会读零值
   - 写已经关闭的chan会panic

## Mutex
### Mutex 的实现原理
Go sync包提供了两种锁类型：互斥锁 sync.Mutex 和读写互斥锁 sync.RWMutex，都属于悲观锁。

**概念:**

Mutex是互斥锁，当一个goroutine获得了锁后，其他 goroutine 不能获取锁 (只能存在一个写者或者读者，不能同时读写)

**使用场景:**

多个线程同时访问临界区，为保证数据的安全，所著一些共享资源，以防止并发访问这些共享数据时可能导致的数据不一致的问题。

获取锁的线程可以正常访问临界区，未获取到锁的线程等待锁释放后可以尝试获取锁

![img.png](./part_channel/img.png)

**底层实现结构:**

互斥锁对应的是底层结构是 sync.Mutex 结构体，位于 src/sync/mutex.go 中

```go
type Mutex struct {
	state int32
	sema  uint32
}
```
state 表示锁的状态，有锁定、被唤醒、饥饿模式等。并且是用state的二进制位来标识的，不同模式下会有不同的处理方式

![img_1.png](./part_mutex/img_1.png)

Sema表示信号量，mutex阻塞队列的定位就是通过这个变量来实现的，从而实现goroutine的阻塞和唤醒

![img_2.png](./part_mutex/img_2.png)

(引入sudog结构体)
```go
type sudog struct {
    // The following fields are protected by the hchan.lock of the
    // channel this sudog is blocking on. shrinkstack depends on
    // this for sudogs involved in channel ops.
    g *g
    
    next *sudog
    prev *sudog
    elem unsafe.Pointer // data element (may point to stack) 指向sema变量
    waitlink *sudog // g.waiting list or semaRoot
    waittail *sudog // semaRoot
    ...
}
```

**操作**

锁的实现一般会依赖于原子操作、信号量，通过atomic包中一些原子操作来实现锁的锁定，通过信号量来实现线程阻塞与唤醒

**加锁**

通过原子操作 cas 加锁，如果加锁不成功就会根据不同场景选择自旋重试加锁或者阻塞等待被唤醒后加锁

![tu](./part_mutex/img_3.png)

```go
func (m *Mutex) Lock() {
    // Fast path: grab unlocked mutex.
    if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
        if race.Enabled {
            race.Acquire(unsafe.Pointer(m))
        }
        return
    }
    // Slow path (outlined so that the fast path can be inlined)
    m.lockSlow() // 尝试自选或者阻塞获取锁
}
```

**解锁**

通过原子操作add解锁，如果任有 goroutine 在等待，唤醒等待的goroutine

![img_4.png](./part_mutex/img_4.png)

```go
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}

	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		// Outlined slow path to allow inlining the fast path.
		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
		m.unlockSlow(new) // 唤醒等待的goroutine
	}
}
```

**注意点:**
- 在 Lock() 之前使用 Unlock() 会导致 panic 异常
- 使用 Lock() 加锁后，再次 Lock() 会导致死锁（不支持重入），需要 Unlock() 解锁后才能再加锁
- 锁定状态与 goroutine 没有关联，一个goroutine 可以 Lock，另一个 goroutine 可以 Unlock

### Go互斥锁正常模式和饥饿模式的区别
Go有两种抢锁的模式，一种是正常模式，另一种是饥饿模式

**正常模式(非公平锁)**

在刚开始的时候，是处于正常模式(Barging)，也就是，当一个G1持有一个锁的时候，G2会自旋的去尝试获取这个锁

当自旋四次还没有能获取到锁的时候，这个G2就会被浇入到获取锁的等待队列里，并阻塞等待唤醒

    正常模式下，所有等待锁的goroutine按照 FIFO 顺序等待，唤醒的 goroutine 不会直接拥有锁，而是回合请求锁的 goroutine 竞争，新请求锁的 gotoutine 具有优势: 它正在CPU上执行，而且可能有好几个，所以刚刚唤醒的 goroutine 有很大可能在竞争中失败，长时间获取不到锁会进入饥饿模式

**饥饿模式(公平锁)**

当一个 goroutine 等待锁时间超过 1 毫秒的时候，它可能会遇到接问题。 在版本 1.9 中，这个场景下 Go Mutex 切换到饥饿模式 handoff 解决接问题

```go
starving = runtime_nanotime()-waitStartTime > 1e6
```

    饥饿模式下，直接把锁交给等待队列中排在第一位的 goroutine(队头)，同时饥饿模式下，新进来的goroutine不会参与抢锁，也不会进入自旋状态，会直接进入等待队列的队尾，这样很好地解决了老的 goroutine 一直抢不到锁的问题

那么也不可能说永远保持一个饥饿状态，总归有要有吃饱的时候，也就是说有那么一个 Mutex 要回到正常模式，那么回归正常模式必须具备的条件有以下几种：

    1. G的执行时间小于 1ms
    2. 等待队列已经全部清空了

当满足上述两个条件的任意一个的时候，Mutex会切换回正常模式，而Go的抢锁过程，就是在这个正常模式和饥饿模式中来回切换进行的。
```go
delta := int32(mutexLocked - 1<<mutexWaiterShift)
if !starving || old>>mutexWaiterShift == 1 {
    // Exit starvation mode.
    // Critical to do it here and consider wait time.
    // Starvation mode is so inefficient, that two goroutines
    // can go lock-step infinitely once they switch mutex
    // to starvation mode.
    delta -= mutexStarving
}
atomic.AddInt32(&m.state, delta)
```

总结:

对于两种模式，正常模式下的性能都是最好的，goroutine 可以连续多次获取锁，饥饿模式就是为了解决锁公平的问题，但是性能会下降

### 互斥锁允许自旋的条件

线程没有获取到锁的时候常见有两种处理方式：
- 一种是没有获取到锁的线程就一直循环等待判断该资源是否已经释放锁，这种锁也叫自旋锁，它不用将线程阻塞起来，适用于并发低且程序执行时间短的场景，缺点是cpu占用高
- 另一种处理方式就是把自己阻塞起来，会释放CPU给其他线程，内核会将线程置为 _睡眠_ 状态, 等到锁被释放后，内核会在合适的实际唤醒该线程，适用于高并发场景，缺点是有线程上下文切换的开销

Go语言中的Mutex实现了自旋与阻塞两种场景，当满足不了自旋条件的时候就会进入阻塞

允许自旋的条件:
1. 锁已经被占用，且锁不处于饥饿状态
2. 积累的自旋次数小鱼最大自旋次数 (active_spin=4)
3. CPU核数大于1
4. 有空闲的P
5. 当前 goroutine 所挂在的P下，本地待运行队列为空

```go
if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
// Active spinning makes sense.
// Try to set mutexWoken flag to inform Unlock
// to not wake other blocked goroutines.
    if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
        atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
        awoke = true
    }
    runtime_doSpin()
    iter++
    old = m.state
    continue
}

func sync_runtime_canSpin(i int) bool {
	if i >= active_spin || ncpu <= 1 || gomaxprocs <= int32(sched.npidle+sched.nmspinning) + 1 {
		return false
    }   
    if p := getg().m.p.ptr(); !runqempty(p){
        return false	
    }
    return false
}
```

**自旋:**
```go
func sync_runtime_doSpin(){
	procyield(active_spin_cnt)
}
```
如果可以进入自旋状态后就调用上面这个方法来进入自旋，doSpin 方法会调用 procyield(30) 执行 30 次 PAUSE 指令，什么都不做，但是会消耗CPU时间

### Go 读写锁的实现原理
读写互斥锁 RWMutex, 是对 Mutex 的一个扩展，当一个 goroutine 获得了读锁之后，其他 goroutine 可以获取读锁，但是不能获取写锁；当一个goroutine获得了写锁后，其他goroutine既不能获取读锁也不能获取写锁（只能存在一个写者或者多个读者，可以同时读）

**使用场景:**

读多余写的情况(既保证线程安全，又保证性能不太差)

**底层实现结构:**

互斥锁对应的底层结构在 src/sync/rwmutex.go 中

```go
type RWMutex struct {
	w           Mutex  // held if there are pending writers
	writerSem   uint32 // semaphore for writers to wait for completing readers
	readerSem   uint32 // semaphore for readers to wait for completing writers
	readerCount int32  // number of pending readers
	readerWait  int32  // number of departing readers
}
```

**操作**

####读锁的加锁与释放

```go
func (rw *RWMutex) RLock()

func (rw *RWMutex) RUnlock()
```
**加读锁**
```go
func (rw *RWMutex) RLock() {
	if race.Enabled {
		_ = rw.w.state
		race.Disable()
	}
	// 为什么readerCount会小于0 因为 writer的lock会对readerCount做减法(原子操作)
	if atomic.AddInt32(&rw.readerCount, 1) < 0 {
		// A writer is pending, wait for it.
		runtime_SemacquireMutex(&rw.readerSem, false, 0)
	}
	if race.Enabled {
		race.Enable()
		race.Acquire(unsafe.Pointer(&rw.readerSem))
	}
}
```
`atomic.AddInt32(&rw.readerCount, 1)` 调用这个原子方法,对当前在读的数量加1，如果返回负数那就说明当前有其他写缩，就调用 `runtime_SemacquireMutex`
休眠当前goroutine等待被唤醒

**释放读锁**

解锁的时候对正在读的操作减1，如果返回值小鱼0那么说明当前有在写的操作，这个时候调用 rUnlockSlow 进入慢速通道

```go
func (rw *RWMutex) RUnlock(){
	if r:= atomic.AddInt32(&rw.readerCount, -1); r < 0 {
		rw.rUnlockSlow(r)
    }   
}
```
被阻塞的准备读的goroutine的数量减一，readerWait 为0，就表示当前没有正在准备读的 goroutine 这时候调用 runtime_Semrelease 唤醒写操作

```go
func (rw *RWMutex) rUnlockSlow(r int32) {
	if r+1 == 0 || r+1 == -rwmutexMaxReaders {
		race.Enable()
		throw("sync: RUnlock of unlocked RWMutex")
	}
	// A writer is pending.
	if atomic.AddInt32(&rw.readerWait, -1) == 0 {
		// The last reader unblocks the writer.
		runtime_Semrelease(&rw.writerSem, false, 1)
	}
}
```

#### 写锁的加锁与释放
```go
func (rw *RWMutex) Lock()

func (rw *RWMutex) Unlock()
```

**加写锁**
```go
func (rw *RWMutex) Lock() {
    if race.Enabled {
        _ = rw.w.state
        race.Disable()
    }
    // First, resolve competition with other writers.
    rw.w.Lock()
    // Announce to readers there is a pending writer.
    r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
    // Wait for active readers.
    if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
        runtime_SemacquireMutex(&rw.writerSem, false, 0)
    }
    if race.Enabled {
        race.Enable()
        race.Acquire(unsafe.Pointer(&rw.readerSem))
        race.Acquire(unsafe.Pointer(&rw.writerSem))
    }
}
```

首先调用互斥锁的lock，获取到互斥锁之后，如果计算之后当前仍然又其他的goroutine持有读锁，那么就调用`runtime_SemacquireMutex` 休眠当前的goroutine
等待所有读锁操作完成

这里的read count 原子性加上一个很大的负数是防止后面的协程能拿到读锁，阻塞读

**释放写锁**

```go
func (rw *RWMutex) Unlock() {
	if race.Enabled {
		_ = rw.w.state
		race.Release(unsafe.Pointer(&rw.readerSem))
		race.Disable()
	}

	// Announce to readers there is no active writer.
	r := atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)
	if r >= rwmutexMaxReaders {
		race.Enable()
		throw("sync: Unlock of unlocked RWMutex")
	}
	// Unblock blocked readers, if any.
	for i := 0; i < int(r); i++ {
		runtime_Semrelease(&rw.readerSem, false, 0)
	}
	// Allow other writers to proceed.
	rw.w.Unlock()
	if race.Enabled {
		race.Enable()
	}
}
```

解锁的操作会西安调用 `atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)` 将回复之前写入的负数，然后根据当前有多少个读操作在等待，循环唤醒

**注意点**

- 读锁或写锁Lock()之前使用 Unlock() 会导致 panic
- 使用 Lock() 加锁后再次 Lock() 会导致死锁（不支持重入），需要 Unlock() 之后才能再加锁
- 锁定状态与 goroutine 没有关联，一个 goroutine 可以 Rlock(Lock), 另一个 goroutine 可以 RUnlock(Unlock)

**互斥锁和读写锁的区别：**

- 读写锁区分读者和写者，而且互斥锁不区分
- 互斥锁同一时间只允许一个线程访问该对象，无论读写；读写锁同一时间只允许一个写者，但是允许多个读者同时读对象

### Go 可重入锁如何实现

可重入锁又称为递归锁，是指在同一个线程在外层方法获取锁的时候，在进入该线程的内层方法时会自动获取锁，不会因为之前已经获取过还没释放再次加锁导致死锁

#### 为什么Go语言种没有可重入锁？
Mutex不是可重入的锁，mutex实现种没有记录哪个 goroutine 拥有这把锁，理论上任何 goroutine 都可以随意地 unlock 这把锁，所以没办法计算重入条件，并且 Mutex 重复Lock 会导致死锁

如何实现可重入锁？

两个重点

- 记住持有锁的线程
- 统计重入的次数

```go
package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type ReentrantLock struct {
	sync.Mutex
	recursion	int32 // goroutine 可重入的次数
	owner 		int64 // 当前持有锁的 goroutine id
}

func GetGoroutineID() int64 {

	var buf [64]byte
	// 获取栈信息
	n := runtime.Stack(buf[:], false)
	// 抽取id
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine"))[0]
	// 转为64位整数
	gid, _ := strconv.Atoi(idField)
	return int64(gid)
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
```

### Go 的原子操作有哪些

Go atomic 包是最轻量级的锁(也称无锁结构)，可以在不形成临界区和创建互斥量的情况下完成并发安全值的替换操作，不过这个包值支持 int32/int64/uint32/uint64/uintptr
这集中数据类型的一些基础操作(增减、交换、载入、存储)

概念：原子操作仅会由一个独立的CPU指令和代表完成，原子操作是无锁的，常常直接通过CPU指令直接实现。事实上，其他同步技术的实现依赖于原子操作

使用场景：

当我们想要对某个变量并发安全的修改，除了使用官方提供的 mutex，还可以使用 sync/atomic 包的原子操作，它能够保证对变量的读取或修改期间不被其他的协程所影响。

atomic包提供的原子操作能够确保任一时刻只有一个goroutine对变量进行操作，善用atomic能够避免程序中出现大量锁操作。

**常见操作：**
- 增减 Add
- 载入 Load
- 比较并交换 CompareAndSwap
- 交换 Swap
- 存储 Store

atomic 操作的对象是一个地址，你需要把可寻址的变量的地址作为参数传递给方法，而不是把变量的值传递给方法

下面分别介绍这些操作

**增减**
```go
func AddInt32(addr *int32, delta int32)(new int32)
func AddInt64(addr *int64, delta int64)(new int64)
func AddUInt32(addr *uint32, delta uint32)(new uint32)
func AddUInt32(addr *uint64, delta uint64)(new uint64)
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)
```

需要注意的是，第一个参数必须是指针类型的值，通过指针变量可以获取被操作数在内存中的地址，从而施加特殊的CPU指令，确保同一时间只有一个goroutine能够操作

```go
fund add(addr *int64, delta int64){
	atomic.AddInt64(addr, delta)
	fmt.Println("add opts: ", *addr)
}
```

**载入**
```go
func LoadInt32(addr *int32) (val int32)
func LoadInt64(addr *int64) (val int64)
func LoadUint32(addr *uint32) (val uint32)
func LoadUint64(addr *uint64) (val uint64)
func LoadUintptr(addr *uintptr) (val uintptr)
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)
```
载入操作能保证原子的读变量的值，当读取的时候，任何其他CPU操作都无法对该变量进行读写，其实是吓死你机制受到底层硬件的支持

**比较并交换**

此类操作的前缀为 CompareAndSwap 简称为 CAS，可以实现乐观锁
```go
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
```
该操作在进行交换前首先确保变量的值未被更改，即依然保持参数old所记录的值，满足此前提条件下才进行交换。CAS的做法类似操作数据库时常见的乐观锁机制

需要注意的是，当由大量的 goroutine 对变量进行读写操作的时候，可能导致 CAS 操作无法成功，这时可以利用 for 循环多次尝试

**交换**
```go
func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)
func SwapUint32(addr *uint32, new uint32) (old uint32)
func SwapUint64(addr *uint64, new uint64) (old uint64)
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
```
相对于CAS，明显此类操作更直接暴力，不管变量的旧值是否被改变，直接赋予新值然后返回替换的值

**存储**
```go
func StoreInt32(addr *int32, val int32)
func StoreInt64(addr *int64, val int64)
func StoreUint32(addr *uint32, val uint32)
func StoreUint64(addr *uint64, val uint64)
func StoreUintptr(addr *uintptr, val uintptr)
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)

func (v *Value)Store(x interface{}){}
```
此类操作确保了写变量的原子性，避免其他操作督导了修改变量过程中的脏数据

### Go的原子操作和锁

1. 原子操作由底层硬件支持，而锁是基于原子操作+信号量完成。若实现相同功能，前者通常会更有效率
2. 原子操作是单个指令的互斥操作；互斥锁/读写锁是一种数据结构，可以完成临界区(多个指令)的互斥操作，扩大原子操作的反胃
3. 原子操作是无锁操作，属于乐观锁；说起锁的时候一般都是悲观锁
4. 原子操作存在于各个指令/语言层级，比如"机器指令层级的原子操作"，“汇编指令层级的原子操作”，“Go语言层级的原子操作”
5. 锁也存在于各个指令/语言层级，比如 "机器指令层级的锁"，“汇编指令层级的锁”，“Go语言层级的锁”等

## 协程 Goroutine
### 底层原理

Goroutine 可以理解为一种 Go 语言的协程(轻量级线程)，是 Go 支持高并发的基础，属于用户态的线程，由 Go runtime 管理而不是操作系统

**底层数据结构**
```go
type g struct {
	stack       stack   // offset known to runtime/cgo
	sched     gobuf     // goroutine 切换的时候用于保存 g 的上下文
	goid         int64  // goroutine 唯一id
	gopc           uintptr         // pc of go statement that created this goroutine
	startpc        uintptr         // pc of goroutine function
    ...
}

type gobuf struct {
    sp   uintptr        // 栈指针的位置
    pc   uintptr        // 运行到的程序的位置
    g    guintptr       // 指向 goroutine
    ctxt unsafe.Pointer
    ret  uintptr        // 保存系统调用的返回值
    lr   uintptr
    bp   uintptr // for framepointer-enabled architectures
}

type stack struct {
    lo uintptr  // 栈的下界内存地址
    hi uintptr  // 栈的上界内存地址
}
```
最终有一个 runtime.g 对象放入调度队列

**状态流转**

|**状态**|**含义**|
|:---:|:---:|
|空闲中_Gidle|G刚刚新建，任未初始化|
|待运行_Grunnable|就绪状态，G在运行队列中，等待M取出并运行|
|运行中_Grunning|M正在运行这个G，这个时候M会拥有一个P|
|系统调用中_Gsyscall|M正在运行这个G发起的系统调用，这个时候M并不拥有P|
|等待中_Gwaiting|G在等待某些条件完成，这时候G不在运行也不在运行队列中(可能在channel的等待队列中)|
|已终止_Gdead|G未被使用，可能已执行完毕|
|栈复制中_Gcopystack|G正在获取一个新的栈空间，并把原来的内容复制过去(用于防止GC扫描)|

![img.png](./part_goroutine/img.png)

#### 创建
通过 `go` 关键字调用底层函数 `runtime.newproc()` 创建一个 `goroutine`

当调用该函数之后， goroutine 会被设置成 runnable 状态

创建好的这个goroutine 会申请一个自己的栈空间，同时在 G 的 sched 中维护栈地址与程序计数器这些信息。

！！！每个 G 在被创建之后，都会被优先放到本地队列，如果本地队列已经满了，就会被放到全局队列中！！！‘

#### 运行
goroutine 本身只是一个数据结构，真正让 goroutine 运行起来的是调度器，Go实现了一个用户态的调度器模型(GMP)
这个调度器充分利用现代计算机的多核心特性，同时让多个 goroutine 运行，同时 goroutine 设计的很轻量级，调度和上下文切换的代价都比较小

![img_1.png](./part_goroutine/img_1.png)

#### 调度时机
- 新启一个协程和协程执行完毕
- 会阻塞的系统调用，比如文件IO、网络 IO
- channel、mutex 等阻塞操作
- time.sleep
- gc
- 主动调用 runtime.Gosched()
- 运行过久或者系统调用过久等

1. 每个 M 开始执行 P 的本地队列中的 G 时，goroutine 会被设置成 `running` 状态
2. 如果某个 M 把本地队列中的 G 都执行完成之后，然后就会去全局队列中拿 G，这里需要注意，每次去全局队列拿 G 的时候都需要上锁，避免同样的任务被多次拿
3. 如果全局队列都被拿完了，而且当前 M 也没有更多的G可以执行的时候，它就回去其他P的本地队列中拿任务，这个机制被称之为 work stealing 机制，每次会拿走一半的任务，向下取整，比如另一个p中有3个任务，那么拿走一半就是拿走一个
4. 当全局队列为空，M也没办法从其他的 P 中拿任务的时候，就会让自身进入自选状态，等待有新的 G 进来。最多只有 GOMAXPROCS 个 M 在自旋状态，过多的 M 自旋会浪费 CPU 资源

#### 阻塞
channel 的读写操作、等待锁、等待网络数据、系统调用等都有可能发生阻塞，会调用底层函数 `runtime.gopark()` 会让出 cpu 时间片，让调度器安排其他等待的任务进行，并在下次某个时候从该位置恢复执行

当调用该函数之后，goroutine 会被设置成 waiting 状态

#### 唤醒
处于 waiting 状态的 goroutine， 在调用 `runtime.goready()` 函数后会被唤醒，唤醒的goroutine 会被放到 M 对应的上下文P对应的 runqueue 中等待被调度

当调用该函数后 goroutine 会处于 runnable 的状态

#### 退出
当 goroutine 执行完成后，会调用底层函数 `runtime.Goexit()`

当调用该函数之后，goroutine会被设置为 dead 状态

### goroutine 和线程的区别

| |goroutine|线程|
|:---:|:---:|:---:|
|内存占用|创建一个goroutine的栈内存消耗为2KB，实际运行过程中如果栈空间不足会自动进行扩容|创建一个线程的栈内存消耗为 1MB|
|创建和销毁|goroutine因为是由Go runtime 负责管理的，创建和销毁的消耗非常小，是用户级|线程创建和销毁都有巨大的内存消耗，因为要和操作系统打交道，是内核级的，通常解决办法就是线程池|
|切换|goroutine切换秩序保存三个寄存器: PC、SP、BP(程序计数器、堆栈指针寄存器，内容为栈顶的偏移地址、基地址寄存器，常用于访问内存时存放内存单元的偏移地址)goroutine切换约 200ns，相当于2400-3600条指令|当线程切换到时候需要保存各种寄存器，方便恢复现场。线程切换会消耗 1000-1500ns 相当于 12000-18000 条指令|

### goroutine 泄露场景
**泄露原因**
- Goroutine 内进行 channel/mutex 等读写操作一直被阻塞
- Goroutine 内的业务逻辑进入死循环，资源一直无法释放
- Goroutine 内的业务逻辑进入长时间等待，有不断新增的 goroutine 入场

**泄露场景**

如果输出的 goroutines 数量是不断在增加的，就说明存在泄露

- **nil channel**

channel 如果忘记初始化，无论你是读还是写都会阻塞

```go
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main(){
	fmt.Println("before goroutine: ", runtime.NumGoroutine())
	block1()
	time.Sleep(time.Second * 1)
	fmt.Println("after goroutine: ", runtime.NumGoroutine())
}

func block1(){
	var ch chan int
	for i := 0; i < 10; i++ {
		go func() {
			<- ch
		}()
	}
}
```
- **发送不接收**

channel发送数量超过channel接收数量，也会造成阻塞

```go
func block2(){
	ch := make(chan int)
	for i := 0; i < 10; i++ {
        go func() {
            ch <- 1
        }()
    }
}
```
- **接收不发送**

channel接收数量超过channel发送数量，也会造成阻塞

```go
func block2(){
	ch := make(chan int)
	for i := 0; i < 10; i++ {
        go func() {
            <-ch
        }()
    }
}
```
- http request body 未关闭

`resp.Body.Close()` 未被调用的时候，goroutine 不会退出

```go
func requestWithNoClose(){
	_, err := http.Get("www.baidu.com")
	if err != nil {
		fmt.Println("err: ", err.Error())
	}
}

func requestWithClose(){
	resp, err := http.Get("www.baidu.com")
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}
	defer resp.Body.Close()
}

var wg sync.WaitGroup

func block4(){
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			requestWithClose()
		}()
	}
}

func main(){
	block4()
	wg.Wait()
}
```
一般发起request 请求 需要关闭 body

- **互斥锁忘记解锁**

第一个协程获取 `sync.Mutex` 加锁了，但是他可能在处理业务逻辑，又或是忘记了 Unlock 了。

因此导致后面的协程想加锁，却因锁未释放被阻塞了
```go
func block5(){
	var mutex sync.Mutex
	for i := 0; i < 10; i++ {
	    go func(){
	    	mutex.Lock()
        }()   	
    }   
}
```

- **sync.waitGroup 使用不当**

由于 `wg.Add` 的数量与 `wg.Done` 的数量不匹配，因此在调用 wg.Wait 方法后一直阻塞等待

```go
func block6(){
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func(){
		    wg.Add(2)
		    wg.Done()
		    wg.Wait()
		}()
    }
}
```

#### 如何排查
单个函数：调用 runtime.NumGoroutine 方法来打印执行代码前后 Goroutine 的运行数量，进行前后比较，就能知道有没有泄露了

生产测试环境：使用 pprof 实时检测Goroutine 的数量

### Go 如何查看正在执行的 goroutine 的数量

在程序中引入 pprof package

```go
package main

import "net/http"

func main(){
	for i := 0; i < 100; i++ {
		go func() {
			select {}
		}()
	}
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	select {}
}
```
**分析goroutine文件**

执行命令
```shell
go tool pprof -http=:1248 http://127.0.0.1:6060/debug/pprof/goroutine
```

会自动打开一个网页

可惜我看不到图  so sad

Could not execute dot; may need to install graphviz.


### 如何控制并发 goroutine 的数量

#### 为什么要控制goroutine并发的数量
在开发过程中，如果对goroutine不加以控制而进行滥用的话，可能会导致服务整体崩溃。比如耗尽系统资源导致程序崩溃，或者CPU使用率过高导致系统忙不过来

#### 用什么方法控制 goroutine 并发的数量
- 有缓冲 channel

利用缓冲满时发生阻塞的特性
```go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg = sync.WaitGroup{}

func main(){
	// 模拟用户请求
	requestCount := 10
	fmt.Println("goroutine_num", runtime.NumGoroutine())
	// 管道长度即最大并发数
	ch := make(chan bool, 3)
	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		ch <- true
		go Read(ch,i)
	}
	wg.Wait()
}

func Read(ch chan bool, i int){
	fmt.Printf("goroutine_num: %d, go func: %d", runtime.NumGoroutine(), i)
	<- ch
	wg.Done()
}
```
输出结果 默认不超过3(4-1)个 goroutine 并发执行

- 无缓冲channel

任务发送和执行分离，指定消费者并发协程数
```go
package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg = sync.WaitGroup{}

func main(){
	// 模拟用户请求
	requestCount := 10
	fmt.Println("goroutine_num", runtime.NumGoroutine())
	// 管道长度即最大并发数
	ch := make(chan bool)
	for i := 0; i < 3; i++ {
		go Read(ch,i)
	}
	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		ch <- true
	}
	wg.Wait()
}

func Read(ch chan bool, i int){
	for _ = range ch {
		fmt.Printf("goroutine_num: %d, go func: %d", runtime.NumGoroutine(), i)
		wg.Done()
	}
}
```

## 调度模型
### Go 线程实现模型
Go实现的是两级线程模型(M:N)，准确的说是GMP模型，是对两级线程模型的改进实现，使它能够更加灵活地进行线程之间的调度。

**背景**

| | 含义|缺点|
|:---:|:---:|:---:|
|单进程时代|每个程序就是一个进程，直到一个程序运行完才能运行下一个进程|无法并发，只能串行\\进程阻塞所带来的CPU时间浪费|
|多进程/线程时代|一个线程阻塞，CPU可以立即切换到其他线程中执行|进程/线程占用内存高\\进程/线程上下文切换成本高|
|协程时代|协程(用户态线程)绑定线程(内核态线程)CPU调度线程执行|实现起来比较复杂，协程和线程的绑定依赖调度器实现(调度器算法)|

线程->CPU由操作系统调度， 协程->线程 由Go调度器来调度，协程与线程的映射关系有三种线程模型

#### 三种线程模型
主要分为：
1. 内核级线程模型
2. 用户级线程模型
3. 两级线程模型

它们的区别在于用户线程与内核线程之间的对应关系

---

- 内核级线程模型  1:1

1个用户线程对应1个内核线程，这种最容易实现，协程的调度都由CPU完成了

优点：
1. 实现起来最简单
2. 能够利用多核
3. 如果进程中要给线程被阻塞不会影响其他线程，是能够切换同一个进程内其他线程继续执行

缺点
1. 上下文切换成本高

---

- 用户级线程模型(n:1)

1个进程中所有线程对应1个内核线程

优点
1. 上下文切换成本低，再用户态即可完成协程切换

缺点
1. 无法利用多核
2. 一旦协程阻塞造成线程阻塞，本线程的其他协程无法执行

---

- 两级线程模型 (M:N)

M个协程对应N个内核线程

优点
1. 利用多核
2. 上下文切换成本低
3. 如果一个线程被阻塞不会阻塞其他线程，是能够切换同一进程内的其他线程继续执行

缺点
1. 实现起来最复杂

### GMP 和 GM

调度器的好坏评价标准: 是否能在适当的实际将合适的协程分配到合适的位置，保证公平和效率

Go采用了 GMP 模型(对两级线程模型的改进实现)，使它能够更加灵活地进行线程之间的调度。

#### GMP
GMP 是 Go 运行时调度层面的实现，包含4个重要结构，分别是G、M、P、Sched

![img.png](./part_mode/img.png)

1. G：取Goroutine的首字母，即用户态的线程，也叫协程；
2. M：取Machine的首字母，和内核线程一一对应，为简单理解，我们可以认为其就是内核线程；
    - **M(Machine)**:Go对操作系统线程(0Sthread)的封装，可以看作操作系统内核线程，想要在 CPU上执行代码必须有线程，通过系统调用 clone 创建。M在绑定有效的P后，进入一个调度循环，而调度循环的机制大致是从P的本地运行队列以及全局队列中获取G，切换到G的执行栈上并执行G的函数，调用goexit 做清理工作并回到M，如此反复。M并不保留G状态，这是G可以跨M调度的基础。M的数量有限制，默认数量限制是10000，可以通过 debug.SetMaxThreads()方法进行设置，如果有M空闲，那么就会回收或者睡眠。
3. P：取Processor的首字母，表示处理器（可以理解成用户态的协程调度器），是G和M之间的中间层，负责协程调度。
    - **P(Processor)**:虚拟处理器，M执行G所需要的资源和上下文，只有将P和M绑定，才能让P的runq中的G真正运行起来。P的数量决定了系统内最大可并行的G的数量，**P的数量受本机的CPU核数影响，可通过环境变量$GOMAXPROCS或在runtime.GOMAXPROCS()来设置，默认为CPU核心数。
4.	Sched 调度器结构，它维护有存储M和G的全局队列以及调度器的一些状态信息

| |G|M|P|
|:---:|:---:|:---:|:---:|
|数量限制|无限制、受机器内存影响|有限制，默认最多 10000| 有限制，最多 GOMAXPROCS 个|
|创建时机|go func| 当没有足够的M来关联P并运行其中的可运行的G的时候就会请求创建新的M|再确定了P的最大数量n后，运行时系统会根据这个数据创建P|

**核心数据结构**
```go
type g struct {
	stack       stack   // offset known to runtime/cgo
	sched     gobuf     // goroutine 切换的时候用于保存 g 的上下文
	goid         int64  // goroutine 唯一id
	gopc           uintptr         // pc of go statement that created this goroutine
	startpc        uintptr         // pc of goroutine function
    ...
}

type p struct {
	lock        mutex
    id          int32
    status      uint32 // one of pidle/prunning/...
    // Queue of runnable goroutines. Accessed without lock.
    runqhead uint32     // 本地队列队头
    runqtail uint32     // 本地队列队尾
    runq     [256]guintptr  // 本地队列 大小为256的数组，数组往往都会被读到缓存中，对缓存友好，效率较高
    runnext guintptr    // 下一个优先执行的 goroutine （一定是最后生产出来的），为了实现局部性原理，runnext的G永远是会被最先调度的
    ...
}

type m struct {
    g0      *g // goroutine with scheduling stack
    // 每个M都有一个自己的G0，不知想任何函数 在系统调用或者调度的时候，M会切换到 G0，使用 G0的栈空间来调度
    curg    *g // 当前正在执行的G
    ...
}

type schedt struct {
    ...
    // Global runnable queue.
    runq     gQueue // 全局队列 链表 无长度限制
    runqsize int32  // 全局队列长度
    ...
}
```
GMP模型的实现算是Go调度器的一大进步，但调度器任然有一个令人头疼的问题，那就是不支持抢占式调度，这导致一旦某个G中出现死循环的代码逻辑，那么G将
永久占用分配给它的 P 和 M，而位于同一个P中的其他G将得不到调度，出现饿死的情况。

当只有要给 P (GOMAXZPROCS=1) 的时候，整个Go程序中其他的G都将被饿死，因此在 Go1.2 版本中实现了基于写作的 “抢占式” 调度，在 Go 1.14 版本中实现了基于信号的 “抢占式” 调度。

#### GM 模型
![img_1.png](./part_mode/img_1.png)

缺点：
1. 主要就是性能不行，全局队列的锁竞争
2. M和G切换的额外开销：当 M1 在执行 G1 的时候，M1 创建了 G2，为了继续执行 G1，需要把G2保存到全局队列中，无法保证 G2是被M1处理。因为M1 原本就保存了 G2 的信息，所以 G2 最好是在 M1 上执行，这样的话也不需要转移 G 到全局队列和线程上下文切换
3. 线程使用效率不能最大化，没有 work-stealing 和 hand-off 机制(这是啥  我也不知道)

### Go 调度的原理
goroutine 调度的本质就是将 **Goroutine(G)** 按照一定算法放到CPU上去运行。

CPU感知不到Goroutine，只知道内核线程，所以需要Go调度器将协程调度到内核线程上面去，然后操作系统调度器将内核级线程放到CPU上去执行

M是对内核级线程的封装，所以Go调度器的工作就是将G分配到M

Go调度器的实现不是一蹴而就的，它的调度模型与算法也是几经演化，从最初的GM模型、到GMP模型，从不支持抢占到支持协作式抢占，再到支持基于信号的异步抢占，经理了不断地优化与打磨

**设计思想**
- 线程复用(work stealing 和 hand off)
- 利用并行(利用多核CPU)
- 抢占调度(解决公平性问题)

**调度对象**

Go调度器

    Go 调度器是属于 Go runtime 的一部分，Go runtime 负责实现 Go 的并发调度、垃圾回收、内存堆栈管理等关键功能

**被调度对象**

G的来源
- P的runnext(只有1个G，局部性原理，永远会被最先调度执行)
- P的本地队列(数组，最多256个G)
- 全局G队列(链表，无限制)
- 网络轮询器 network poller (存放网络调用被阻塞的G)

P的来源
- 全局P队列(数组，GOMAXPROCS个P)

M的来源
- 休眠线程队列(未绑定P，长时间休眠会等待 GC 回收销毁)
- 运行线程(绑定P，指向P中的G)
- 自旋线程(绑定P，指向M的G0)

其中运行线程数+自旋线程数 <= P的数量(GOMAXPROCS), M个数 >= P个数

**调度流程**

协程的调度采用了生产者消费者模型，实现了用户任务与调度器的解耦

![img_2.png](./part_mode/img_2.png)

1. go creates a new goroutine
2. newly created goroutine being put into local or global queue
3. A M is being waken or created to execute goroutine
4. Schedule loop
5. Try its best to get a goroutine to execute
6. Clear, reenter schedule loop

![img_3.png](./part_mode/img_3.png)

生产端我们开启的每个协程都是一个计算任务，这些任务会被提交给go的runtime。如果计算任务非常多，有成千上万个，那么这些任务是不可能同时被立刻执行的
所以这个计算任务一定会被先暂存起来，一般的做法是放到内存的队列中等待被执行

G的生命周期：G从创建、保存、被获取、调度和执行、阻塞、销毁，步骤如下：
1. 步骤1：创建G，关键字 go func() 创建G
2. 步骤2：保存G，创建的G优先保存到本地队列P，如果P满了，则会平衡部分P到全局队列中
3. 步骤3：唤醒或者新建M执行任务，进入循环调度(4,5,6)
4. 步骤4：M获取G，M首先从P的本地队列获取G，如果P为空，则从全局队列中获取G，如果全局队列也为空，则从另一个本地队列中偷取一半的G，这种从其他P偷的方式称之为 work stealing
5. 步骤5：M调度和执行G，M调用 G.func() 函数执行G
    - 如果M在执行G的过程中发生系统调用阻塞(同步)，会阻塞G和M(操作系统限制)，此时P会和当前M解绑寻找新的M，如果没有空闲的M就会创建一个新的M，接管正在阻塞G所属的P，接着继续执行P中其余的G，这种阻塞后释放P的方式称之为hand off。当系统调用结束后，这个G会尝试获取一个空闲的P执行，优先获取之前绑定的P，并放入到这个P的本地队列，如果获取不到P那么这个线程M变成休眠状态，加入到空闲县城后，然后这个G会被放入到全局队列中
    - 如果M在执行G的过程中发生网络IO等操作阻塞时(异步)，阻塞G，不会阻塞M。M会寻找P中其他可执行的G继续执行，G会被网络轮询器 network poller 接手，当阻塞的G恢复后，G1从network poller 被移回到P的LRQ中，重新进入可执行状态。异步情况下，通过调度，Go scheduler 成功地将 I/O 的任务转变成为了 CPU任务，或者说将内核级线程切换转变成了用户级别的 goroutine 切换，大大提高了效率
6. 步骤6：M执行完G后清理现场，重新进入调度循环(将M上运行的goroutine切换为G0，G0负责调度时协程的切换)

其中步骤2中保存G的详细流程如下：
- 执行 go func 的时候，主线程M0回调用 newproc()
- 每个协程 G 都会被尝试先放到 P 中的 runnext，若 runnext 为空则放到 runnext 中，生产结束
- 若 runnext 已满，则将原来 runnext 中的 G 踢到本地队列中，将当前G放到 runnext 中，生产结束
- 若本地队列也满了，则将本地队列中的 G 拿出一半，放到全局队列中，生产结束

**调度时机**

什么时候进行调度(执行/切换)？

在一下情形下，回切换正在执行的goroutine

- 抢占式调度
    - sysmon检测到协程运行过久 （比如说sleep，死循环）
        - 切换到 g0， 进入循环调度

- 主动调度
    - 新起一个协程和协程执行完毕
        - 触发调度循环
    - 主动调用 runtime.Gosched()
        - 切换到g0，进入调度循环
    - 垃圾回收之后
        - stw之后，会重新选择g开始执行

- 被动调度
    - 系统调用(比如文件IO)阻塞(同步)
        - 阻塞G和M，P和M分离，将P交给其他M绑定，其他M执行P的剩余G
    - 网络IO调用阻塞(异步)
        - 阻塞G，G移动到NewPoller，M执行P的剩余G
    - atomic/mutex/channel 等阻塞(异步)
        - 阻塞G，G移动到 channel 的等待队列中，M执行P的剩余G

**调度策略**

使用什么调度策略来挑选下一个goroutine执行？

由于P中的G分布在 runnext、本地队列、全局队列、网络轮询器中，则需要挨个判断是否有可执行的G，大体逻辑如下：
- 每执行61次调度循环，从全局队列获取G，若有则直接返回
- 从P上的runnext看一下是否有G，若有则直接返回
- 从P上的本地队列看一下是否有G，若有则直接返回
- 上面都没查找到时，则去全局队列、网络轮询器查找或者从其他P中窃取，一直阻塞到获取一个可用的G为止

源码如下：
```go
func schedinit() {
	// 取出g0
	_g_ := getg()
	if raceenabled {
		_g_.racectx, raceprocctx0 = raceinit()
	}	

	// 设置最大线程数
	sched.maxmcount = 10000
	
	...	

	// 初始化m0, 前边已经将m0和g0的关系绑定好了
	// 只是检查一下各种变量，然后将m0挂到allm链表中
	mcommoninit(_g_.m)

	...

	sched.lastpoll = uint64(nanotime())
	// ncpu在osinit时已经获取
	procs := ncpu
	// 如果GOMAXPROCS设置并且合法就将procs的设置为GOMAXPROCS
	if n, ok := atoi32(gogetenv("GOMAXPROCS")); ok && n > 0 {
		procs = n
	}

	if procresize(procs) != nil {
		throw("unknown runnable goroutine during bootstrap")
	}

	...	
}
```

### Go 的 work stealing 机制
当线程M无可运行的G的时候，尝试从其他M绑定的P偷取G，减少空转，提高了线程利用率(避免闲着不干活)。

当从本线程绑定P本地队列、全局G队列、netpoller都找不到可执行的g，会从别的P里窃取G并放到当前P上面。

从 netpoller 中拿到的G 是 Gwating 状态(存放的是因为网络IO被阻塞的G)，从其他地方拿到的G是 Grunnable状态

从全局队列取G数量: N = min(len(GRQ)/GOMAXPROCS+1, len(GRQ/2)) (根据GOMAXPROCS负载均衡)

从其他P本地队列窃取的G数量： N = len(LRQ)/2 （平分）

**窃取流程**

源码见 runtime/proc.go stealWork函数 流程如下：如果金国多次努力一直找不到需要运行的goroutine则调用stopm进入睡眠状态，等待被其他工作线程唤醒
1. 选择要窃取的P
2. 从P中偷走一半G

**选择要窃取的P**

窃取的实质就说遍历allp中的所有P，查看其运行队列是否有goroutine，如果有，则窃取一半到当前工作线程的运行队列

为了保证公平性，遍历allP的时候不是固定从 allp[0] 开始，而是从随机位置上的P开始，而且遍历的顺序也随机化了，不是i到i+1，而是使用了一种伪随机的方式遍历每个p，防止每次遍历时使用同样的顺序访问allP中的元素

**从p偷走一半g**

源码见 runtime/proc.go runqsteal函数

挑选出盗取的对象p后，则调用runqsetal盗取p的运行队列中的goroutine，runqsteal函数再调用runqgrap从p的本地队列队尾拿走一半的G

### hand off 机制
也成为P分离机制，当本线程M因为G进行的系统调用阻塞时，线程释放绑定的P，把P转移给其他空闲的M执行，也提高了线程利用率(避免占着茅坑不拉屎)

**分离流程**

当前线程M阻塞，释放P，给其他空闲的M处理
![img_4.png](./part_mode/img_4.png)

### Go 抢占式调度
**基于协作的抢占式调度**

- 协作式：打架都按事先定义好的规则来，比如：一个 goroutine 执行完后，退出，让出P，然后下一个goroutine被调度到P上运行。这样做的缺点在于是否让出p的决定权在
  goroutine自身。一旦某个g不主动让出p或执行时间较长，那么后面的goroutine只能等着，没有方法让前者让出p，导致延迟甚至饿死

- 非协作式：就是由 runtime 来决定一个 goroutine 运行多长时间，如果你不主动让出，对不起，我有手段可以抢占你，把你踢出去，让后面的goroutine来运行

基于协作的抢占式调度流程：
1. 编译器会在函数调用前插入 runtime.morestack，让运行时有机会在这段代码中检查是否需要执行抢占式调度
2. Go语言运行时会在垃圾回收暂停程序、系统监控发现 Goroutine 运行超过 10ms，那么会在这个协程设置一个抢占标记
3. 当发生函数调用时，可能会执行编译器插入的 runtime.morestack，它调用的 runtime.newstack 会检查抢占标记，如果由抢占标记就会触发抢占让出cpu，切到调度主协程里

这种方案只能说局部解决了“饿死问题”，只在由函数调用的地方才能插入“抢占”代码(埋点)，对于没有函数调用而是纯算法循环计算的G，go调度器依然无法抢占

比如，死循环等并没有给编译器插入抢占代码的机会，以下程序在 go 1.14 之前的版本中会一直卡住，而不会打印 i got scheduled
```go
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main(){
	runtime.GOMAXPROCS(1)
	go func() {
		for  {
            
		}
    }()
    time.Sleep(time.Second)
	fmt.Println("i got scheduled")
}
```

为了解决这些问题，Go再1.14中添加了对非协作的抢占调度的支持，这种方式是基于系统信号的，也就是向线程向线程发送信号的方式来抢占运行的 goroutine

**基于信号的抢占式调度**

真正的抢占式调度是基于信号完成的，也成为"异步抢占"。不管协程有没有意愿主动让出 cpu 运行权，只要某个协程执行时间过长，就会发送信号强行夺取 cpu 运行权

- M注册一个 SIGURG 信号的处理函数：sighandler
- sysmon启动后会间隔性的进行监控，最长 10ms，最短 20us。如果发现协程独占超过 10ms 会给M发送抢占信号
- M收到信号后内核执行 sighandler 函数把当前协程状态从 Grunning 改为 Grunnable，把抢占的协程放到全局队列里，M继续寻找其他的goroutine来执行
- 被抢占的G再次调度过来执行，会继续原来的执行流

抢占式分为 _Prunning 和 _Psyscall, syscall 通常是由于阻塞性系统调用引起的，比如磁盘IO、cgo；Prunning一般是由四玄幻的计算逻辑引起的

### Go如何查看运行时的调度信息
由两种方式可以查看一个程序的调度GMP信息，分别是 go tool trace 和 GODEBUG

```go
// trace.go
package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

func main(){
	// 创建 trace 文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()
	// 启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	// main
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		fmt.Println("Hello")
	}
}
```
![img_5.png](./part_mode/img_5.png)

点击 `view trace` 能够看见可视化的调度流程

**GoDebug**

GODEBUG变量可以控制运行时内的调试变量。查看调度器信息，将会使用如下两个参数：
- schedtrace：设置 schedtrace=X 参数可以使运行时在每X毫秒发出一行调度器的摘要信息到标准err输出中。
- scheddtail：设置 schedtrace=X 和 scheddtail=1 可以使运行时在每个X毫秒发出一次详细的多行信息，信息内容主要包括调度程序、处理器、OS线程和Goroutine 状态

```shell
go build trace.go GODEBUG=schedtrace=1000 ./trace
```

