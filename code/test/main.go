package main

import "fmt"

func main(){
	A1()
}

func A1(){
	defer func() {
		fmt.Println("defer test one")
	}()

	//err := fmt.Errorf("嘟嘟嘟嘟嘟嘟")
	//if err != nil {
	//	return
	//}

	defer func() {
		fmt.Println("defer test")
	}()
}
