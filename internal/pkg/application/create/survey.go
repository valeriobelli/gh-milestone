package create

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/valeriobelli/gh-milestone/internal/pkg/domain/github"
	"github.com/valeriobelli/gh-milestone/internal/pkg/infrastructure/editor"
)

type SurveyAnswers struct {
	Confirm     bool
	Description string
	DueDate     time.Time
	Title       string
}

func (sa SurveyAnswers) getTime() *time.Time {
	if sa.DueDate.IsZero() {
		return nil
	}
	return &sa.DueDate
}

type Flags struct {
	Description string
	DueDate     *time.Time
	Title       string
}

func (f Flags) getDueDate() time.Time {
	if f.DueDate == nil {
		return time.Time{}
	}
	return *f.DueDate
}

type Survey struct {
	flags Flags
}

func NewSurvey(flags Flags) *Survey {
	return &Survey{flags: flags}
}

type step int

const (
	stepTitle step = iota
	stepDescription
	stepDueDate
	stepConfirm
	stepDone
)

type model struct {
	step      step
	answers   SurveyAnswers
	flags     Flags
	textInput textinput.Model
	err       error
	quitting  bool
}

func initialModel(flags Flags) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := model{
		flags:     flags,
		textInput: ti,
		answers: SurveyAnswers{
			Description: flags.Description,
			Title:       flags.Title,
			DueDate:     flags.getDueDate(),
			Confirm:     true, // Default to true if not interactive
		},
	}

	m.answers.Confirm = len(flags.Title) == 0

	// Determine initial step
	if flags.Title == "" {
		m.step = stepTitle
		m.textInput.Placeholder = "Title"
	} else {
		// If title is provided via flags, we skip everything else according to original logic
		m.step = stepDone
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			m.err = fmt.Errorf("aborted")

			return m, tea.Quit
		case tea.KeyEnter:
			return m.handleKeyEnter()
		case tea.KeyRunes:
			if m.step == stepDescription && msg.String() == "e" {
				return m, m.openEditor
			}
		}

	case editorResultMsg:
		if msg.err != nil {
			m.err = msg.err

			return m, nil
		}

		m.answers.Description = msg.content
		m.nextStep()

		if m.step == stepDone {
			return m, tea.Quit
		}

		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

type editorResultMsg struct {
	content string
	err     error
}

func (m model) openEditor() tea.Msg {
	content, err := editor.Open(m.answers.Description, "*.md")

	return editorResultMsg{content, err}
}

func (m *model) validateCurrentStep() error {
	val := m.textInput.Value()

	switch m.step {
	case stepTitle:
		if strings.TrimSpace(val) == "" {
			return fmt.Errorf("title validates: value is required")
		}
	case stepDueDate:
		if strings.TrimSpace(val) == "" {
			return nil
		}

		_, err := github.NewDueDate(val)

		return err
	}
	return nil
}

func (m *model) commitValue() {
	val := m.textInput.Value()

	switch m.step {
	case stepTitle:
		m.answers.Title = val
	case stepDueDate:
		if strings.TrimSpace(val) != "" {
			d, _ := github.NewDueDate(val)

			if d != nil {
				m.answers.DueDate = d.Time
			}
		}
	}
}

func (m *model) nextStep() {
	m.textInput.Reset()
	m.err = nil

	current := m.step

	switch current {
	case stepTitle:
		m.step = stepDescription
	case stepDescription:
		m.step = stepDueDate
		m.textInput.Placeholder = "Due date [yyyy-mm-dd]"
	case stepDueDate:
		m.step = stepConfirm
		m.textInput.Placeholder = "Do you want create the Milestone? (y/N)"
	case stepConfirm:
		m.step = stepDone
	}
}

func (m model) View() string {
	if m.step == stepDone {
		return ""
	}

	var s strings.Builder

	if m.err != nil {
		s.WriteString(fmt.Sprintf("Error: %s\n", m.err))
	}

	switch m.step {
	case stepTitle:
		s.WriteString("Title\n")
		s.WriteString(m.textInput.View())
	case stepDescription:
		s.WriteString("Description\n")
		s.WriteString(fmt.Sprintf("[(e) to launch %s, enter to skip]", editor.ReadDefaultEditor()))
	case stepDueDate:
		s.WriteString("Due date [yyyy-mm-dd]\n")
		s.WriteString(m.textInput.View())
	case stepConfirm:
		s.WriteString("Do you want create the Milestone? (y/N)\n")
		s.WriteString(m.textInput.View())
	}

	s.WriteString("\n\n(esc to quit)")

	return s.String()
}

func (s Survey) Ask() (SurveyAnswers, error) {
	p := tea.NewProgram(initialModel(s.flags))
	m, err := p.Run()

	if err != nil {
		return SurveyAnswers{}, err
	}

	finalModel := m.(model)

	if finalModel.quitting && finalModel.err != nil {
		return SurveyAnswers{}, finalModel.err
	}

	// If aborted with ctrl+c but no explicit error set (other than default)
	if finalModel.quitting {
		return SurveyAnswers{}, fmt.Errorf("aborted")
	}

	return finalModel.answers, nil
}

func (m model) handleKeyEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case stepDescription:
		m.nextStep()

		return m, nil
	case stepConfirm:
		val := strings.ToLower(strings.TrimSpace(m.textInput.Value()))

		if val == "y" || val == "yes" {
			m.answers.Confirm = true
		} else {
			m.answers.Confirm = false
		}

		m.step = stepDone

		return m, tea.Quit
	}

	// Validate and move next
	if err := m.validateCurrentStep(); err != nil {
		m.err = err

		return m, nil
	}

	m.err = nil
	m.commitValue()
	m.nextStep()

	if m.step == stepDone {
		return m, tea.Quit
	}

	return m, nil
}
