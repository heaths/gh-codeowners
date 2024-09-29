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

func TestPR(t *testing.T) {
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
			name: "single page",
			fs: fstest.MapFS{
				"CODEOWNERS": {Data: []byte(content)},
			},
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repository": {
								"pullRequest": {
									"files": {
										"nodes": [
											{
												"path": "main.go",
												"changeType": "MODIFIED"
											},
											{
												"path": "docs/README.md",
												"changeType": "ADDED"
											}
										],
										"pageInfo": {
											"hasNextPage": false
										}
									}
								}
							}
						}
					}`)
			},
			wantStdout: `[{"path":"main.go","changeType":"MODIFIED","owners":["@heaths"]},{"path":"docs/README.md","changeType":"ADDED","owners":["@writers"]}]`,
		},
		{
			name: "multiple pages (tty)",
			tty:  true,
			fs: fstest.MapFS{
				"CODEOWNERS": {Data: []byte(content)},
			},
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repository": {
								"pullRequest": {
									"files": {
										"nodes": [
											{
												"path": "main.go",
												"changeType": "MODIFIED"
											}
										],
										"pageInfo": {
											"hasNextPage": true,
											"endCursor": "abcd1234"
										}
									}
								}
							}
						}
					}`)
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"data": {
							"repository": {
								"pullRequest": {
									"files": {
										"nodes": [
											{
												"path": "docs/README.md",
												"changeType": "ADDED"
											}
										],
										"pageInfo": {
											"hasNextPage": false
										}
									}
								}
							}
						}
					}`)
			},
			wantStdout: heredoc.Doc(`
				[
				  {
				    "path": "main.go",
				    "changeType": "MODIFIED",
				    "owners": [
				      "@heaths"
				    ]
				  },
				  {
				    "path": "docs/README.md",
				    "changeType": "ADDED",
				    "owners": [
				      "@writers"
				    ]
				  }
				]
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(gock.Off)

			fake := console.Fake(
				console.WithStdoutTTY(tt.tty),
				console.WithColorScheme(nil),
			)
			repo, err := repository.Parse("heaths/gh-codeowners")
			require.NoError(t, err)

			opts := prOptions{
				GlobalOptions: &GlobalOptions{
					Console: fake,
					Repo:    repo,

					host:          "github.com",
					authToken:     "***",
					colorDisabled: true,
					fs:            tt.fs,
				},
			}

			if tt.mocks != nil {
				tt.mocks()
			}

			err = pr(&opts)
			require.NoError(t, err)

			stdout, _, _ := fake.Buffers()
			assert.Equal(t, tt.wantStdout, stdout.String())
		})
	}
}
