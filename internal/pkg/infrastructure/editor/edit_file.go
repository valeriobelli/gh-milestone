package editor

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"

	shellquote "github.com/kballard/go-shellquote"
)

type showable interface {
	Show() error
}

func needsBom() bool {
	return runtime.GOOS == "windows"
}

func editFile(editorCommand, fn, initialValue string, stdin io.Reader, stdout io.Writer, stderr io.Writer, cursor showable) (string, error) {
	pattern := fn
	if pattern == "" {
		pattern = "survey*.txt"
	}
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())

	if needsBom() {
		if _, err := f.Write(bom); err != nil {
			return "", err
		}
	}

	if _, err := f.WriteString(initialValue); err != nil {
		return "", err
	}

	if err := f.Close(); err != nil {
		return "", err
	}

	if editorCommand == "" {
		editorCommand = defaultEditor
	}
	args, err := shellquote.Split(editorCommand)
	if err != nil {
		return "", err
	}
	args = append(args, f.Name())

	editorExe, err := exec.LookPath(args[0])
	if err != nil {
		return "", err
	}
	args[0] = editorExe

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if cursor != nil {
		_ = cursor.Show()
	}

	if err := cmd.Run(); err != nil {
		return "", err
	}

	raw, err := os.ReadFile(f.Name())
	if err != nil {
		return "", err
	}

	return string(bytes.TrimPrefix(raw, bom)), nil
}
