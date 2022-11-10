package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("🍎")
	fmt.Println(runtime.GOMAXPROCS(0))
	bar := NewProgressBar(100)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second / 10)
		bar.Done(1)
	}
	bar.Finish()
}
