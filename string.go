package goshared

import "strings"

func FlatPath(path string) []string {
	var r []string

	// First element
	if len(path) == 0 || path[0:1] != "/" {
		panic("invalid absolute path")
	}
	r = append(r, path[0:1])

	// Middle elements
	tmpPath := strings.Replace(path, "/", " ", 1)
	for {
		offset := strings.Index(tmpPath, "/")
		if offset != -1 {
			r = append(r, path[:offset])
			tmpPath = strings.Replace(tmpPath, "/", " ", 1)
		} else {
			break
		}
	}

	// Latest element
	r = append(r, path)

	return r
}
