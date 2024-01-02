package main

import (
	"fmt"
	"time"

	"github.com/jasonlvhit/gocron"
)

func job() {
	fmt.Println("job: xx")
}

func jobWithParams(a, b string) {
	fmt.Println("job with param: ", a, b)
}

func main() {
	gocron.Every(1).Second().Do(job)
	gocron.Every(1).Second().Do(jobWithParams, "a", "b")
	closeChan := gocron.Start()
	defer close(closeChan)

	time.Sleep(5 * time.Second)
}
