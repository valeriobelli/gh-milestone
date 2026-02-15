package editor

import (
	"os"
	"runtime"

	"github.com/cli/go-gh/v2/pkg/config"
)

var (
	bom           = []byte{0xef, 0xbb, 0xbf}
	defaultEditor = "nano"
)

func init() {
	if gh := os.Getenv("GH_EDITOR"); gh != "" {
		defaultEditor = gh
	} else if ghCfg := ghConfigEditor(); ghCfg != "" {
		defaultEditor = ghCfg
	} else if runtime.GOOS == "windows" {
		defaultEditor = "notepad"
	} else if g := os.Getenv("GIT_EDITOR"); g != "" {
		defaultEditor = g
	} else if v := os.Getenv("VISUAL"); v != "" {
		defaultEditor = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		defaultEditor = e
	}
}

func ghConfigEditor() string {
	cfg, err := config.Read(nil)

	if err != nil {
		return ""
	}

	val, err := cfg.Get([]string{"editor"})

	if err != nil {
		return ""
	}

	return val
}

func Open(initialValue, fileName string) (string, error) {
	return editFile(defaultEditor, fileName, initialValue, os.Stdin, os.Stdout, os.Stderr)
}

func ReadDefaultEditor() string {
	return defaultEditor
}
