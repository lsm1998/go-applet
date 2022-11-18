package main

type Copy interface {
	Copy(src, dist string) error
}
