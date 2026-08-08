package main

import "path/filepath"

func resolveLinkTarget(linkPath, target string) string {
	if filepath.IsAbs(target) {
		return filepath.Clean(target)
	}

	return filepath.Clean(filepath.Join(filepath.Dir(linkPath), target))
}
