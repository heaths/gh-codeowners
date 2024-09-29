package cmd

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/heaths/go-console"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestView(t *testing.T) {
	content := heredoc.Doc(`
		# comment
		* @heaths
		docs/ @writers
	`)
	tests := []struct {
		name       string
		tty        bool
		fs         fs.FS
		mocks      func()
		wantStdout string
	}{
		{
			name: "no errors",
			fs: fstest.MapFS{
				"CODEOWNERS": {Data: []byte(content)},
			},
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON("{}")
			},
			wantStdout: content,
		},
		{
			name: "no errors (tty)",
			tty:  true,
			fs: fstest.MapFS{
				"CODEOWNERS": {Data: []byte(content)},
			},
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON("{}")
			},
			wantStdout: heredoc.Docf(`
				%[1]s[0;38;2;0;255;0m# comment%[1]s[0m
				* @heaths
				docs/ @writers
			`, "\033"),
		},
		{
			name: "errors",
			fs: fstest.MapFS{
				".github/CODEOWNERS": {Data: []byte{}},
				"CODEOWNERS":         {Data: []byte(content)},
			},
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repository": {
								"codeowners": {
									"errors": [
										{
											"path": "CODEOWNERS",
											"kind": "Unknown owner",
											"line": 3,
											"column": 7,
											"source": "docs/ @writers"
										}
									]
								}
							}
						}
					}`)
			},
			wantStdout: content,
		},
		{
			name: "errors (tty)",
			tty:  true,
			fs: fstest.MapFS{
				".github/CODEOWNERS": {Data: []byte{}},
				"CODEOWNERS":         {Data: []byte(content)},
			},
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repository": {
								"codeowners": {
									"errors": [
										{
											"path": "CODEOWNERS",
											"kind": "Unknown owner",
											"line": 3,
											"column": 7,
											"source": "docs/ @writers"
										}
									]
								}
							}
						}
					}`)
			},
			wantStdout: heredoc.Docf(`
				%[1]s[0;38;2;0;255;0m# comment%[1]s[0m
				* @heaths
				docs/ %[1]s[0;38;2;255;0;0m@writers%[1]s[0m
			`, "\033"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(gock.Off)

			fake := console.Fake(console.WithStdoutTTY(tt.tty))
			repo, err := repository.Parse("heaths/gh-codeowners")
			require.NoError(t, err)

			opts := viewOptions{
				GlobalOptions: &GlobalOptions{
					Color: ColorOptions{
						Comment: "#00FF00",
						Error:   "#FF0000",
					},
					Console: fake,
					Repo:    repo,

					host:      "github.com",
					authToken: "***",
					fs:        tt.fs,
				},
			}

			if tt.mocks != nil {
				tt.mocks()
			}

			err = view(&opts)
			require.NoError(t, err)

			stdout, _, _ := fake.Buffers()
			assert.Equal(t, tt.wantStdout, stdout.String())
		})
	}
}
