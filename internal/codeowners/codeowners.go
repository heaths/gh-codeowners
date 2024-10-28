package codeowners

import (
	_fs "io/fs"

	"github.com/hairyhenderson/go-codeowners"
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
	var err error
	var stat _fs.FileInfo

	if statFS, ok := fs.(_fs.StatFS); ok {
		stat, err = statFS.Stat(path)
	} else if f, e := fs.Open(path); e != nil {
		return false
	} else {
		defer f.Close()
		stat, err = f.Stat()
	}

	return err == nil && !stat.IsDir()
}

type Codeowners struct {
	source *codeowners.Codeowners
}

func (c Codeowners) Owners(path string) []string {
	return c.source.Owners(path)
}

func Open(fs _fs.FS, path string) (*Codeowners, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, err
	}

	if c, err := codeowners.FromReader(f, ""); err != nil {
		return nil, err
	} else {
		return &Codeowners{
			source: c,
		}, nil
	}
}
