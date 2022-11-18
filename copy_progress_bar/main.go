package main

import (
	"fmt"
	"time"
)

func main() {
	bar := NewProgressBar(100, WithTail(">"), WithFiller("="), WithInterval(time.Second/100))
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second / 10)
		if bar.Done(1.75) {
			break
		}
	}
	fmt.Println("Done!!!")
}
