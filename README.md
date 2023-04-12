# {{param "name" (param "github.repo") "What is your project name?" | titlecase}}

{{if (param "badges" true "Do you need badges?") -}}
[![releases](https://img.shields.io/github/v/release/{{param "github.owner"}}/{{param "github.repo"}}.svg?logo=github)](https://github.com/{{param "github.owner"}}/{{param "github.repo"}}/releases/latest)
[![reference](https://pkg.go.dev/badge/github.com/{{param "github.owner"}}/{{param "github.repo"}}.svg)](https://pkg.go.dev/github.com/{{param "github.owner"}}/{{param "github.repo"}})
[![ci](https://github.com/{{param "github.owner"}}/{{param "github.repo"}}/actions/workflows/ci.yml/badge.svg?event=push)](https://github.com/{{param "github.owner"}}/{{param "github.repo"}}/actions/workflows/ci.yml)
{{- end -}}

<!-- {{if 0}} -->
To create a new repository from this template repository for Go projects,
using the [GitHub CLI](https://github.com/cli/cli) run:

```bash
gh extension install heaths/gh-template
gh template clone <name> --template heaths/template-golang --public

# Recommended
cd <name>
git commit -a --amend
```

The `gh template` command will:

1. Create a new repository with the given `<name>` on GitHub.
2. Copy the `heaths/template-golang` files into that repo.
3. Clone the new repository into a directory named `<name>` in the current directory.
4. Apply built-in and passed parameters, or prompt for undefined parameters, to format template files.

This will create a new repo with the given `<name>` in GitHub, copy the
`heaths/template-golang` files into that repo, and clone it into a
subdirectory of the current directory named `<name>`.

See [heaths/gh-template](https://github.com/heaths/gh-template) for more information
about this GitHub CLI extension.
<!-- {{end}} -->

## License

Licensed under the [MIT](LICENSE.txt) license.
