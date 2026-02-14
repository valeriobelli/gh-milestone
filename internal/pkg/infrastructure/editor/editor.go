// Package editor provides an extended survey.Editor prompt with more flexible behavior.
// This is derived from github.com/cli/cli/v2/pkg/surveyext to avoid pulling in the
// entire GitHub CLI as a dependency.
package editor

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
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

// ghConfigEditor reads the "editor" value from `gh` config (set via `gh config set editor`).
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

// GhEditor extends survey.Editor to give it more flexible behavior.
type GhEditor struct {
	*survey.Editor
	EditorCommand string
	BlankAllowed  bool
}

func (e *GhEditor) editorCommand() string {
	if e.EditorCommand == "" {
		return defaultEditor
	}

	return e.EditorCommand
}

var editorQuestionTemplate = `
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ .Config.HelpInput }} for help]{{color "reset"}} {{end}}
  {{- if and .Default (not .HideDefault)}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
  {{- color "cyan"}}[(e) to launch {{ .EditorCommand }}{{- if .BlankAllowed }}, enter to skip{{ end }}] {{color "reset"}}
{{- end}}`

type editorTemplateData struct {
	survey.Editor
	EditorCommand string
	BlankAllowed  bool
	Answer        string
	ShowAnswer    bool
	ShowHelp      bool
	Config        *survey.PromptConfig
}

func (e *GhEditor) prompt(initialValue string, config *survey.PromptConfig) (interface{}, error) {
	err := e.Render(
		editorQuestionTemplate,
		editorTemplateData{
			Editor:        *e.Editor,
			BlankAllowed:  e.BlankAllowed,
			EditorCommand: filepath.Base(e.editorCommand()),
			Config:        config,
		},
	)
	if err != nil {
		return "", err
	}

	rr := e.NewRuneReader()
	_ = rr.SetTermMode()
	defer func() { _ = rr.RestoreTermMode() }()

	cursor := e.NewCursor()
	_ = cursor.Hide()
	defer func() {
		_ = cursor.Show()
	}()

	for {
		r, _, err := rr.ReadRune()
		if err != nil {
			return "", err
		}
		if r == 'e' {
			break
		}
		if r == '\r' || r == '\n' {
			if e.BlankAllowed {
				return initialValue, nil
			}
			continue
		}
		if r == terminal.KeyInterrupt {
			return "", terminal.InterruptErr
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		if string(r) == config.HelpInput && e.Help != "" {
			err = e.Render(
				editorQuestionTemplate,
				editorTemplateData{
					Editor:        *e.Editor,
					BlankAllowed:  e.BlankAllowed,
					EditorCommand: filepath.Base(e.editorCommand()),
					ShowHelp:      true,
					Config:        config,
				},
			)
			if err != nil {
				return "", err
			}
		}
	}

	stdio := e.Stdio()
	text, err := editFile(e.editorCommand(), e.FileName, initialValue, stdio.In, stdio.Out, stdio.Err, cursor)
	if err != nil {
		return "", err
	}

	if len(text) == 0 && !e.AppendDefault {
		return e.Default, nil
	}

	return text, nil
}

// Prompt implements the survey.Prompt interface.
func (e *GhEditor) Prompt(config *survey.PromptConfig) (interface{}, error) {
	initialValue := ""
	if e.Default != "" && e.AppendDefault {
		initialValue = e.Default
	}
	return e.prompt(initialValue, config)
}
