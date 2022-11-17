package file_sharding

// Merge 文件合并
type Merge interface {
	Merge() (string, error)
}
