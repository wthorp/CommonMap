//go:build !windows

package commonmap

type DWORDLONG uint64

type folderEntry struct {
	name   string
	parent DWORDLONG
}
