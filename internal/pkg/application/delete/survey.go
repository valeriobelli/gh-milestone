package delete

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	ghub "github.com/google/go-github/v68/github"
)

type SurveyAnswers struct {
	Confirm bool
}

type Config struct {
	Confirm   bool
	Milestone *ghub.Milestone
}

type Survey struct {
	config Config
}

func NewSurvey(config Config) *Survey {
	return &Survey{config: config}
}

type model struct {
	textInput textinput.Model
	milestone *ghub.Milestone
	confirm   bool
	quitting  bool
	err       error
}

func initialModel(milestone *ghub.Milestone) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		milestone: milestone,
	}
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
			val := m.textInput.Value()
			expected := strconv.Itoa(*m.milestone.Number)
			if val == expected {
				m.confirm = true
			} else {
				m.confirm = false
			}
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"You're going to delete milestone #%d (%s). This action cannot be reversed. To confirm, type the milestone number:\n\n%s\n\n(esc to quit)",
		*m.milestone.Number, *m.milestone.Title,
		m.textInput.View(),
	)
}

func (s Survey) Ask() (SurveyAnswers, error) {
	if s.config.Confirm {
		return SurveyAnswers{Confirm: true}, nil
	}

	p := tea.NewProgram(initialModel(s.config.Milestone))
	m, err := p.Run()

	if err != nil {
		return SurveyAnswers{}, err
	}

	finalModel := m.(model)
	if finalModel.quitting && finalModel.err != nil {
		return SurveyAnswers{}, finalModel.err
	}

	if finalModel.quitting {
		return SurveyAnswers{}, fmt.Errorf("aborted")
	}

	return SurveyAnswers{Confirm: finalModel.confirm}, nil
}
