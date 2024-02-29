package main

import "fmt"

func escape1() *int{
	var a int = 1
	return &a
}

func escape2(){
	s := make([]int, 0, 10000)
	for index, _ := range s {
		s[index] = index
	}
}

func escape3() {
	number := 10
	s := make([]int, number) // 编译期间无法确定 slice 的长度
	for i := 0; i < len(s); i++ {
		s[i] = i
	}
}

func escape4() {
	fmt.Println(1111)
}

func escape5()func() int {
	var i int = 1
	return func() int {
		i++
		return i
	}
}

func main() {
	escape1()
	escape2()
	escape3()
	escape4()
	escape5()
}
