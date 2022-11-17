package file_sharding

// Sharding 文件分片
type Sharding interface {
	Sharding() ([]string, error)
}
