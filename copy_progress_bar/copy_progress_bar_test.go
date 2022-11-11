package copy_progress_bar

import (
	"fmt"
	"runtime"
	"testing"
)

func TestNewFileCopyProgressBar(t *testing.T) {
	fmt.Println(runtime.GOMAXPROCS(0))
	fmt.Println("Done!!!")
	printEmoji()

	fcpb := NewFileCopyProgressBar()
	if err := fcpb.Copy("E:\\gopath", "E:\\gopath_copy"); err != nil {
		panic(err)
	}
}
