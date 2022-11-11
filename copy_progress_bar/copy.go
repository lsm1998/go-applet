package copy_progress_bar

type Copy interface {
	Copy(src, dist string) error
}
