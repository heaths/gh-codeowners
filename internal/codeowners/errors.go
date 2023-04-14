package codeowners

import (
	"sort"
	"strings"
	"unicode"

	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/shurcooL/graphql"
)

type ErrorKind string

const (
	ErrorKindUnknownOwner ErrorKind = "Unknown owner"
)

type Error struct {
	Kind    ErrorKind `json:"kind"`
	Path    string    `json:"path"`
	Line    int       `json:"line"`
	Column  int       `json:"column"`
	Source  string    `json:"source"`
	Message string    `json:"message"`
}

func (e Error) UnknownOwner() string {
	if e.Kind == ErrorKindUnknownOwner && e.Column > 0 {
		owner := e.Source[e.Column-1:]
		if idx := strings.IndexFunc(owner, func(r rune) bool {
			return unicode.IsSpace(r)
		}); idx > 0 {
			return owner[:idx]
		}
		return owner
	}

	return ""
}

type Errors []Error

func (e Errors) Path() string {
	for _, e := range e {
		return e.Path
	}

	return ""
}

func (e Errors) UnknownOwners() []string {
	unknown := make(map[string]bool)
	for _, e := range e {
		if owner := e.UnknownOwner(); owner != "" {
			unknown[owner] = true
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	owners := make([]string, 0, len(unknown))
	for k := range unknown {
		owners = append(owners, k)
	}
	sort.Strings(owners)

	return owners
}

func (e Errors) indexUnknownOwners() map[int][]string {
	index := make(map[int][]string, len(e))
	for _, e := range e {
		index[e.Line] = append(index[e.Line], e.UnknownOwner())
	}

	return index
}

func QueryErrors(client api.GQLClient, repo repository.Repository, ref string) (Errors, error) {
	var query struct {
		Repository struct {
			Codeowners struct {
				Errors Errors
			} `graphql:"codeowners(refName: $ref)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner": graphql.String(repo.Owner()),
		"repo":  graphql.String(repo.Name()),
		"ref":   graphql.String(ref),
	}

	err := client.Query("CodeownersErrors", &query, variables)
	if err != nil {
		return nil, err
	}

	return query.Repository.Codeowners.Errors, nil
}
