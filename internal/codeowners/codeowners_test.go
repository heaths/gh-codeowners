package codeowners

import (
	_fs "io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		fs   _fs.FS
		want string
	}{
		{
			name: "missing",
			fs:   fstest.MapFS{},
		},
		{
			name: ".github/CODEOWNERS",
			fs: fstest.MapFS{
				".github/CODEOWNERS": {Data: []byte{}},
				"CODEOWNERS":         {Data: []byte{}},
				"docs/CODEOWNERS":    {Data: []byte{}},
			},
			want: ".github/CODEOWNERS",
		},
		{
			name: "CODEOWNERS",
			fs: fstest.MapFS{
				"CODEOWNERS":      {Data: []byte{}},
				"docs/CODEOWNERS": {Data: []byte{}},
			},
			want: "CODEOWNERS",
		},
		{
			name: "docs/CODEOWNERS",
			fs: fstest.MapFS{
				"docs/CODEOWNERS": {Data: []byte{}},
			},
			want: "docs/CODEOWNERS",
		},
		{
			name: "directory",
			fs: fstest.MapFS{
				".github/CODEOWNERS": {Mode: _fs.ModeDir},
				"CODEOWNERS":         {Data: []byte{}},
			},
			want: "CODEOWNERS",
		},
		{
			name: "base FS",
			fs: baseFS{
				".github/CODEOWNERS": baseFileInfo{isDir: true},
				"CODEOWNERS":         baseFileInfo{},
			},
			want: "CODEOWNERS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Find(tt.fs)
			assert.Equal(t, tt.want, got)
		})
	}
}

type baseFS map[string]baseFileInfo

func (fs baseFS) Open(name string) (_fs.File, error) {
	if f, ok := fs[name]; ok {
		f.name = name
		return f, nil
	}
	return nil, _fs.ErrNotExist
}

type baseFileInfo struct {
	name  string
	isDir bool
}

func (f baseFileInfo) Stat() (_fs.FileInfo, error) {
	return f, nil
}
func (f baseFileInfo) Read(buf []byte) (int, error) {
	return 0, nil
}
func (f baseFileInfo) Close() error {
	return nil
}
func (f baseFileInfo) Name() string {
	return f.name
}
func (f baseFileInfo) Size() int64 {
	return 0
}
func (f baseFileInfo) Mode() _fs.FileMode {
	if f.isDir {
		return _fs.ModeDir
	}
	return 0
}
func (f baseFileInfo) ModTime() time.Time {
	return time.Now()
}
func (f baseFileInfo) IsDir() bool {
	return f.isDir
}
func (f baseFileInfo) Sys() any {
	return nil
}
