package git

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/execabs"
)

func Exec(args ...string) (stdout, stderr bytes.Buffer, err error) {
	path, err := execabs.LookPath("git")
	if err != nil {
		err = fmt.Errorf("could not find git executable: %w", err)
		return
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run git: %w, error: %s", err, stderr.String())
		return
	}

	return
}

func RootFS() (fs.FS, error) {
	stdout, _, err := Exec("rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("failed to find git root: %w", err)
	}

	path := strings.TrimSpace(stdout.String())
	return os.DirFS(path), nil
}

func BranchRef() (string, error) {
	stdout, _, err := Exec("rev-parse", "--symbolic-full-name", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get branch name: %w", err)
	}

	name := strings.TrimSpace(stdout.String())
	return name, nil
}
