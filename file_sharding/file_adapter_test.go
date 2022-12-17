package file_sharding

import (
	"fmt"
	"testing"
)

func Test_fileAdapter_Merge(t *testing.T) {
	adapter := NewFileAdapter(
		WithNumber(5),
		WithFilename("/Users/liushiming/Desktop/1.mp4"))
	merge, err := adapter.Merge()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(merge)
}
