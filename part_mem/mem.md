## GO 内存
### Go 内存分配机制
Go语言内置运行时 (runtime) 抛弃了传统的内存分配方式，改为自主管理。这样可以自主地实现更好的内存使用模式，比如内存池，预分配等等。这样，不会每次内存分配的时候都需要进行系统钓友。

**设计思想**

- 内存分配算法采用 Google 的 TCMalloc 算法，每个线程都会自行维护一个独立的内存池，进行内存分配时优先从该内存池中分配，当内存池不足时才会向加锁向全局内存池申请，减少系统调用并且避免不同线程对全局内存池的锁竞争
- 把内存切分的非常细小，分为多级管理，以降低锁的粒度
- 回收对象内存时，并没有将其真正释放掉，只是放回预先分配的大块内存中，以便复用。只有内存闲置过多的时候，才会尝试归还部分内存给操作系统，降低整体开销

**分配组件**

Go的内存管理组件主要有: mspan、mcache、mcentral 和 mheap

![img.png](img.png)

内存管理单元 mspan

`mspan` 是内存管理的基本单元，该结构体中包含 next 和 prev 两个字段，他们分别指向了前一个后一个 mspan，每个 mspan都管理 npages 个大小为 8kb的页，
一个span是由多个page组成的，这里的页不是操作系统中的内存页，它们是操作系统内存页的整数倍。

page是内存存储的基本单元，“对象”放到page中
```go
type mspan struct {
	next        *mspan // 后指针
	prev        *mspan // 前指针
	startAddr   uintptr // 管理页的起始地址 指向page
	npages      unitptr // 页面数量
	spanclass   spanClass // 规格
}

type spanClass uint8
```

Go 有 68 种不同大小的 spanClass, 用于小对象的分配

```go
const _NumSizeClasses = 68
var class_to_size = [_NumSizeClasses]uint16{0, 8, 16, 32, 64, 80, 96, 112, 128 ...}
```
如果按照序号为1的spanClass(对象规格为8B)分配，每个span占用的字节数: 8k, mspan 可以保存 1024 个对象

如果按照序号为2的spanClass(对象规格为16B)分配，每个span占用的字节数: 8k, mspan 可以保存 512 个对象

...

如果按照序号为67的spanClass(对象规格为32k)分配，每个span占用的字节数: 32k, mspan 可以保存 1 个对象

当大于32k的对象出现的时候，回直接从heap分配一个特殊的 span，这个特殊的span类型是0，只包含了一个大对象

**线程缓存:mcache**

mcache 管理线程在本地缓存的 mspan，每个goroutine绑定的P都有一个mcache字段
```go
type mcache struct {
	alloc [numSpanClasses]*mspan
}
_NumSizeClasses = 68
numSpanClasses = _NumSizeClasses << 1
```
mcache 用 span classes 作为索引管理多个用于分配的 mspan，它包含所有规格的mspan。它是 _NumSizeClasses 的两倍，其中 *2 是将spanClass 分成了又指针和没有指针两种，
方便垃圾回收。对于每种规格，又两个mspan，一个mspan不包含指针，另一个mspan则包含指针。对于无指针对象的mspan在进行垃圾回收的时候无需进一步扫描它是否引用了其他活跃的对象。

mcache 在初始化的时候是没有任何 mspan 资源的，在使用过程中回动态地从 mcentral 申请，只会回缓存下来，当对象小鱼等于32kb大小的时候，使用 mcache 的相应规格的 mspan 进行分配

**中心缓存: mcentral**

mcentral

```go
type mcentral struct {
	spanclass   spanClass// 当前规格大小
	partial [2]spanSet  // 有空闲的
	full    [2]spanSet  // 没有空闲object的mspan列表
}
```
每个mcentral管理一种spanClass的mspan，并将有空闲空间和没有空闲空间的貌似盘分开管理。

partial和full的数据类型为 spanSet，表示 mspans 集，可以通过 pop、push 来获得 mspans

```go
type spanSet struct {
	spineLock   mutex
	spine       unsafe.Pointer  // 指向 []span 指针
	spineLen    uintptr     // spin array length， accessed atomically
	spineCap    uintptr
	index       headTailIndex   // 前三十二位是头指针，后三十二位是尾指针
}
```
简单说下 mcache 从 mcentral 获取和归还 mspan 的流程:
- 获取; 加锁，从 partial 链表找到一个可用的 mspan；并将其从 partial 链表删除；将取出的 mspan 加入到 full 链表;将 mspan 返回给工作线程，解锁。
- 归还; 加锁，将 mspan 从 full 链表删除; 将 mspan 加入到 partial 链表，解锁

**页堆 mheap**

mheap管理Go的所有动态分配内存，可以认为是Go程序持有的整个堆空间，全局唯一
```go
type mheap struct {
	lock  mutex     // 全局锁
	pages pageAlloc // 页面分配的数据结构
	allspans []*mspan // 所有通过 mheap_ 申请的 mspans 也就是 堆
	arenas [1 << arenaL1Bits]*[1 << arenaL2Bits]*heapArena
	// 所有中心缓存 mcentral
	central [numSpanClasses]struct {
		mcentral mcentral
		pad      [cpu.CacheLinePadSize - unsafe.Sizeof(mcentral{})%cpu.CacheLinePadSize]byte
	}
	...
}
```
所有 mcentral 的集合则是存放于 mheap 中的。 mheap 里的 arens 区域是堆内存的抽象，运行时会将 8kb 看作一页，这些内存页中存储了所有在堆上初始化的对象。
运行时使用二维的 runtime.heapArena 数组管理所有的内存，每个 runtime.heapArena 都会管理64MB的内存。

当申请内存时，依次 金国 mcache 和 mcentral 都没有合适规格的大小就会向mheap 申请一个块内存，然后按照指定规格划分为列表，并添加到相同规格大小的 mcentral 的非空闲列表后面
