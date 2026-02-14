package create

import (
	"time"

	"github.com/AlecAivazis/survey/v2"
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
	answers   *SurveyAnswers
	questions []*survey.Question
}

func NewSurvey(flags Flags) *Survey {
	var questions = []*survey.Question{}

	requiredFieldsAreEmpty := len(flags.Title) == 0

	if flags.Title == "" {
		questions = append(questions, &survey.Question{
			Name:     "title",
			Prompt:   &survey.Input{Message: "Title"},
			Validate: survey.Required,
		})
	}

	if flags.Description == "" && requiredFieldsAreEmpty {
		questions = append(questions, &survey.Question{
			Name: "description",
			Prompt: &editor.GhEditor{
				BlankAllowed: true,
				Editor: &survey.Editor{
					FileName: "*.md",
					Message:  "Description",
				},
			},
		})
	}

	if flags.DueDate == nil && requiredFieldsAreEmpty {
		questions = append(questions, &survey.Question{
			Name:   "dueDate",
			Prompt: &survey.Input{Message: "Due date [yyyy-mm-dd]"},
			Validate: survey.Validator(func(ans interface{}) error {
				switch dueDate := ans.(type) {
				case string:
					parseDueDate, _ := github.NewDueDate(dueDate)

					if parseDueDate == nil {
						return nil
					}

					return nil
				default:
					return nil
				}
			}),
			Transform: survey.Transformer(func(ans interface{}) (newAns interface{}) {
				switch dueDate := ans.(type) {
				case string:
					parsedDate, err := github.NewDueDate(dueDate)

					if err != nil {
						return time.Time{}
					}

					return parsedDate.Time
				default:
					return nil
				}
			}),
		})
	}

	if len(questions) > 0 {
		questions = append(questions, &survey.Question{
			Name: "confirm",
			Prompt: &survey.Confirm{
				Message: "Do you want create the Milestone?",
			},
		})
	}

	return &Survey{
		answers: &SurveyAnswers{
			Confirm:     !requiredFieldsAreEmpty,
			Description: flags.Description,
			DueDate:     flags.getDueDate(),
			Title:       flags.Title,
		},
		questions: questions,
	}
}

func (s Survey) Ask() (SurveyAnswers, error) {
	err := survey.Ask(s.questions, s.answers)

	return *s.answers, err
}
