package copy_progress_bar

import (
	"fmt"
	"testing"
	"time"
)

func TestNewFileCopyProgressBar(t *testing.T) {
	fmt.Println("Done!!!")
	//printEmoji()

	newProgressBar := NewProgressBar(100)

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second / 10)
		newProgressBar.Done(1)
	}

	//bar := NewFileCopyProgressBar()
	//if err := bar.Copy("E:\\gopath", "E:\\gopath_copy"); err != nil {
	//	panic(err)
	//}
}