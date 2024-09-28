package cmd

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/cli/go-gh/pkg/jsonpretty"
)

func printJson(opts *GlobalOptions, v any) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r := bytes.NewBuffer(buf)
	if opts.Console.IsStdoutTTY() {
		return jsonpretty.Format(opts.Console.Stdout(), r, indent, opts.IsColorEnabled())
	}

	_, err = io.Copy(opts.Console.Stdout(), r)
	return err
}
