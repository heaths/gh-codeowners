package codeowners

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
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

func TestQueryErrors(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func()
		want    Errors
		wantErr error
	}{
		{
			name: "query error",
			mocks: func() {
				gock.New("https://api.github.com").
					Post("/graphql").
					Reply(200).
					JSON(`{
						"errors": [
							{
								"message": "no CODEOWNERS found"
							}
						]
					}`)
			},
			wantErr: assert.AnError,
		},
		{
			name: "validation errors",
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
											"kind": "Unknown owner",
											"line": 6,
											"column": 9,
											"source": "docs/** @writers"
										}
									]
								}
							}
						}
					}`)
			},
			want: []Error{
				{
					Kind:   ErrorKindUnknownOwner,
					Line:   6,
					Column: 9,
					Source: "docs/** @writers",
				},
			},
		},
	}

	repo, err := repository.Parse("heaths/gh-codeowners")
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(gock.Off)

			if tt.mocks != nil {
				tt.mocks()
			}

			client, err := gh.GQLClient(&api.ClientOptions{
				Host:      "github.com",
				AuthToken: "***",
			})
			assert.NoError(t, err)

			errors, err := QueryErrors(client, repo, "main")
			if tt.wantErr != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, errors)
			assert.True(t, gock.IsDone(), "pending mocks: %v", Mocks(gock.Pending()))
		})
	}
}

type Mocks []gock.Mock

func (m Mocks) String() string {
	paths := make([]string, len(m))
	for i, mock := range m {
		paths[i] = mock.Request().URLStruct.String()
	}

	return fmt.Sprintf("%d unmatched mocks: %s", len(paths), strings.Join(paths, ", "))
}
