package main

import "fmt"

func main() {
	Test(nil)
}

func Test(key interface{}) {
	if key == nil {
		fmt.Println("xxxx")
	}
}
