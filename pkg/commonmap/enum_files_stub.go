//go:build !windows

package commonmap

func EnumFiles(basePath string, isRpf func(string) bool) (map[DWORDLONG]folderEntry, []folderEntry) {
	return nil, nil
}
