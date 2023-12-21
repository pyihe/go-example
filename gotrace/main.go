package main

import (
	"os"
	"runtime/trace"
)

// go run main.go 2> trace.out
// go tool trace trace.out
// https://mp.weixin.qq.com/s?__biz=MzA4ODg0NDkzOA==&mid=2247487157&idx=1&sn=cbf1c87efe98433e07a2e58ee6e9899e&source=41#wechat_redirect

func main() {
	trace.Start(os.Stderr)
	defer trace.Stop()

	ch := make(chan string)
	go func() {
		ch <- "TEST"
	}()
	<-ch
}
