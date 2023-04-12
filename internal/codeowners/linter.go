package codeowners

import (
	"bufio"
	"fmt"
	_fs "io/fs"
	"sort"
	"strings"
	"unicode"

	"github.com/heaths/go-console"
)

type ErrorKind string

const (
	ErrorKindUnknownOwner ErrorKind = "Unknown owner"
)

type Error struct {
	Kind   ErrorKind
	Line   int
	Column int
	Source string
	Path   string
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

type LintOptions struct {
	Console console.Console
	Fix     bool
}

func Lint(fs _fs.FS, errors Errors, opts LintOptions) error {
	path := errors.Path()
	if path == "" {
		return nil
	}

	file, err := fs.Open(path)
	if err != nil {
		return err
	}

	cs := opts.Console.ColorScheme()
	rem := cs.Red

	linenum := 0
	missing := errors.indexUnknownOwners()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linenum++

		line := scanner.Text()
		if owners, ok := missing[linenum]; ok {
			for _, owner := range owners {
				line = strings.ReplaceAll(line, owner, rem(owner))
			}
		}

		fmt.Fprintln(opts.Console.Stdout(), line)
	}

	return nil
}
