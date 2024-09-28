package codeowners

import (
	_fs "io/fs"
)

func Find(fs _fs.FS) string {
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

func fileExists(fs _fs.FS, path string) bool {
	if fs, ok := fs.(_fs.StatFS); ok {
		if stat, err := fs.Stat(path); err == nil && !stat.IsDir() {
			return true
		}
	}
	if f, err := fs.Open(path); err == nil {
		defer f.Close()
		if stat, err := f.Stat(); err == nil && !stat.IsDir() {
			return true
		}
	}
	return false
}
