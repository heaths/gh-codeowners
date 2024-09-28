package codeowners

import (
	_fs "io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		fs   _fs.StatFS
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Find(tt.fs)
			assert.Equal(t, tt.want, got)
		})
	}
}
