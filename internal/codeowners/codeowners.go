package codeowners

import (
	_fs "io/fs"
)

func Find(fs _fs.StatFS) string {
	// Based on https://docs.github.com/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners
	lookup := []string{
		".github/CODEOWNERS",
		"CODEOWNERS",
		"docs/CODEOWNERS",
	}
	for _, path := range lookup {
		if fileExists(fs, path) {
			return path
		}
	}
	return ""
}

func fileExists(fs _fs.StatFS, path string) bool {
	if stat, err := fs.Stat(path); err == nil && !stat.IsDir() {
		return true
	}
	return false
}
