package file_sharding

// Sharding ๆไปถๅ็
type Sharding interface {
	Sharding() ([]string, error)
}
