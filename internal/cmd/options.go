package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"strconv"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/cli/go-gh/pkg/term"
	"github.com/heaths/gh-codeowners/internal/git"
	"github.com/heaths/go-console"
	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	Color   ColorOptions
	Console console.Console
	Log     *log.Logger
	Repo    repository.Repository
	Verbose bool

	// Test-only options.
	host          string
	authToken     string
	colorDisabled bool
	fs            fs.FS
}

type ColorOptions struct {
	Comment string
	Error   string
}

func (opts *GlobalOptions) EnsureRepository() (err error) {
	if opts.Repo == nil {
		opts.Repo, err = gh.CurrentRepository()
		if err != nil {
			return
		}
	}

	if opts.Repo == nil {
		return fmt.Errorf("no repository")
	}

	return
}

func (opts *GlobalOptions) IsAuthenticated() error {
	// Make sure the user is authenticated.
	host := opts.Repo.Host()
	if host == "" {
		host, _ = auth.DefaultHost()
	}

	token, _ := auth.TokenForHost(host)
	if token == "" {
		return fmt.Errorf("use `gh auth login` to authenticate with required scopes")
	}

	return nil
}

func (opts *GlobalOptions) IsColorEnabled() bool {
	return !opts.colorDisabled &&
		!term.IsColorDisabled() &&
		opts.Console != nil &&
		opts.Console.IsStdoutTTY()
}

func (opts *GlobalOptions) RootFS() (fs.FS, error) {
	if opts.fs == nil {
		var err error
		opts.fs, err = git.RootFS()
		if err != nil {
			return nil, err
		}
	}
	return opts.fs, nil
}

func StringEnumVarP(cmd *cobra.Command, p *string, name, shorthand, defaultValue string, values []string, usage string) {
	*p = defaultValue
	val := &enumValue{
		value:  p,
		values: values,
	}

	cmd.Flags().VarP(val, name, shorthand, fmt.Sprintf("%s: {%s}", usage, strings.Join(values, "|")))
	_ = cmd.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return values, cobra.ShellCompDirectiveNoFileComp
	})
}

type enumValue struct {
	value  *string
	values []string
}

func (v *enumValue) String() string {
	return *v.value
}

func (v *enumValue) Set(s string) error {
	if !stringSliceContains(s, v.values) {
		return fmt.Errorf("valid values are {%s}", strings.Join(v.values, "|"))
	}
	*v.value = s
	return nil
}

func (v *enumValue) Type() string {
	return "string"
}

func stringSliceContains(value string, values []string) bool {
	for _, v := range values {
		if strings.EqualFold(value, v) {
			return true
		}
	}

	return false
}

func parseNumberRef(number string) (int, error) {
	number = strings.TrimPrefix(number, "#")
	if i, err := strconv.ParseInt(number, 10, 32); err != nil {
		return 0, fmt.Errorf("parse issue or pull request number: %w", err)
	} else {
		return int(i), nil
	}
}
