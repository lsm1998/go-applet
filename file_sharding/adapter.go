package file_sharding

type Adapter interface {
	Merge
	Sharding
}
