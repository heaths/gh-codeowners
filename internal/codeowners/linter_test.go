package codeowners

import (
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/MakeNowJust/heredoc"
	"github.com/heaths/go-console"
	"github.com/stretchr/testify/assert"
)

func TestError_UnknownOwner(t *testing.T) {
	const source = "testdata/** @foo @bar"
	tests := []struct {
		name   string
		kind   ErrorKind
		column int
		want   string
	}{
		{
			name: "empty",
		},
		{
			name:   "unknown owner",
			kind:   ErrorKindUnknownOwner,
			column: 13,
			want:   "@foo",
		},
		{
			name:   "other",
			kind:   "Other",
			column: 18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := Error{
				Kind:   tt.kind,
				Line:   1,
				Column: tt.column,
				Source: source,
			}
			assert.Equal(t, tt.want, sut.UnknownOwner())
		})
	}
}

func TestErrors_UnknownOwners(t *testing.T) {
	const source = "testdata/** @foo @bar"
	tests := []struct {
		name    string
		kind    ErrorKind
		columns []int
		want    []string
	}{
		{
			name: "empty",
		},
		{
			name:    "unknown owner",
			kind:    ErrorKindUnknownOwner,
			columns: []int{13},
			want:    []string{"@foo"},
		},
		{
			name:    "unknown owners",
			kind:    ErrorKindUnknownOwner,
			columns: []int{13, 18},
			want:    []string{"@bar", "@foo"},
		},
		{
			name:    "other",
			kind:    "Other",
			columns: []int{18},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := make(Errors, 0, len(tt.columns))
			for _, v := range tt.columns {
				sut = append(sut, Error{
					Kind:   tt.kind,
					Line:   1,
					Column: v,
					Source: source,
				})
			}
			assert.Equal(t, tt.want, sut.UnknownOwners())
		})
	}
}

func TestLint(t *testing.T) {
	var source = heredoc.Doc(`
		# License

		* @default # Default owner(s)
		docs/** @writers @unknown
	`)

	const path = ".github/CODEOWNERS"
	mockFS := fstest.MapFS{
		path: {Data: []byte(source)},
	}

	tests := []struct {
		name    string
		errors  Errors
		tty     bool
		want    string
		wantErr string
	}{
		{
			name: "unknown owner",
			errors: Errors{
				{
					Kind:   ErrorKindUnknownOwner,
					Line:   4,
					Column: 18,
					Source: "docs/** @writers @unknown",
					Path:   path,
				},
			},
			want: source,
		},
		{
			name: "unknown owner (tty)",
			errors: Errors{
				{
					Kind:   ErrorKindUnknownOwner,
					Line:   4,
					Column: 18,
					Source: "docs/** @writers @unknown",
					Path:   path,
				},
			},
			tty: true,
			want: heredoc.Docf(`
				%[1]s[0;38;2;0;255;0m# License%[1]s[0m

				* @default %[1]s[0;38;2;0;255;0m# Default owner(s)%[1]s[0m
				docs/** @writers %[1]s[0;38;2;255;0;0m@unknown%[1]s[0m
			`, "\033"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			con := console.Fake(
				console.WithStdout(stdout),
				console.WithStdoutTTY(tt.tty),
			)

			opts := LintOptions{
				Console: con,
				Color: struct {
					Comment string
					Error   string
				}{
					Comment: "#00FF00",
					Error:   "#FF0000",
				},
			}

			err := Lint(mockFS, tt.errors, opts)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, stdout.String())
		})
	}
}
