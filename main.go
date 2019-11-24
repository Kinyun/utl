package main

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
)

func task() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

func main() {
	// also, you can create a new scheduler
	// to run two schedulers concurrently
	s := gocron.NewScheduler()
	s.Every(1).Seconds().Do(task)
	<-s.Start()
}
