package delete

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	ghub "github.com/google/go-github/v44/github"
)

type SurveyAnswers struct {
	Confirm bool
}

type Config struct {
	Confirm   bool
	Milestone *ghub.Milestone
}

type Survey struct {
	answers   *SurveyAnswers
	questions []*survey.Question
}

func NewSurvey(config Config) *Survey {
	var questions = []*survey.Question{}

	if !config.Confirm {
		questions = append(questions, &survey.Question{
			Name: "confirm",
			Prompt: &survey.Input{
				Message: fmt.Sprintf("You're going to delete milestone #%d (%s). This action cannot be reversed. To confirm, type the milestone number:", *config.Milestone.Number, *config.Milestone.Title),
			},
			Transform: func(ans interface{}) (newAns interface{}) {
				switch milestoneNumber := ans.(type) {
				case string:
					return strconv.Itoa(*config.Milestone.Number) == milestoneNumber
				default:
					return false
				}
			},
		})
	}

	return &Survey{
		answers: &SurveyAnswers{
			Confirm: config.Confirm,
		},
		questions: questions,
	}
}

func (s Survey) Ask() (SurveyAnswers, error) {
	err := survey.Ask(s.questions, s.answers)

	return *s.answers, err
}
