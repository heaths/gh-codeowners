package codeowners

import (
	"bufio"
	"fmt"
	_fs "io/fs"
	"strings"

	"github.com/heaths/go-console"
)

type RenderOptions struct {
	Console console.Console
	Fix     bool

	Color struct {
		Comment string
		Error   string
	}
}

func Render(fs _fs.FS, errors Errors, opts RenderOptions) error {
	path := errors.Path()
	if path == "" {
		return nil
	}

	file, err := fs.Open(path)
	if err != nil {
		return err
	}

	cs := opts.Console.ColorScheme()
	remove := cs.ColorFunc(opts.Color.Error)
	comment := cs.ColorFunc(opts.Color.Comment)

	linenum := 0
	missing := errors.indexUnknownOwners()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linenum++

		line := scanner.Text()

		if opts.Console.IsStdoutTTY() {
			if owners, ok := missing[linenum]; ok {
				for _, owner := range owners {
					line = strings.ReplaceAll(line, owner, remove(owner))
				}
			}

			if idx := strings.IndexRune(line, '#'); idx >= 0 {
				line = line[:idx] + comment(line[idx:])
			}
		}

		fmt.Fprintln(opts.Console.Stdout(), line)
	}

	return nil
}
