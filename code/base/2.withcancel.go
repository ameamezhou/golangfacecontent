package main

// 很常见的一个案例，假设有一个获取ip的协程，但是这是一个非常耗时的操作每用户随时可能会取消
// 如果用户取消了，那么之前那个获取协程的函数就要停止了


