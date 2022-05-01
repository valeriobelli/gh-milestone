package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cli/cli/v2/pkg/surveyext"
	"github.com/spf13/cobra"
	"github.com/valeriobelli/gh-milestones/internal/pkg/application/create"
	"github.com/valeriobelli/gh-milestones/internal/pkg/domain/constants"
)

func NewCreateCommand() *cobra.Command {
	var getDescription = func(command *cobra.Command) *string {
		description, _ := command.Flags().GetString("description")

		if description == "" {
			return nil
		}

		return &description
	}

	var parseDueDate = func(dueDate string) (*time.Time, error) {
		parsedDueDate, err := time.Parse(constants.DateFormat, dueDate)

		if err != nil {
			return nil, err
		}

		return &parsedDueDate, nil
	}

	var getDueDate = func(command *cobra.Command) (*time.Time, error) {
		dueDate, err := command.Flags().GetString("due-date")

		if err != nil {
			return nil, err
		}

		if dueDate == "" {
			return nil, nil
		}

		return parseDueDate(dueDate)
	}

	var getTitle = func(command *cobra.Command) *string {
		title, _ := command.Flags().GetString("title")

		if title == "" {
			return nil
		}

		return &title
	}

	var questions = []*survey.Question{
		{
			Name:      "title",
			Prompt:    &survey.Input{Message: "Title"},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
		{
			Name: "description",
			Prompt: &surveyext.GhEditor{
				BlankAllowed: true,
				Editor: &survey.Editor{
					FileName: "*.md",
					Message:  "Description",
				},
				EditorCommand: os.Getenv("EDITOR") + " --wait",
			},
		},
		{
			Name:   "dueDate",
			Prompt: &survey.Input{Message: "Due date [yyyy-mm-dd]"},
			Validate: survey.Validator(func(ans interface{}) error {
				switch dueDate := ans.(type) {
				case string:
					_, err := parseDueDate(dueDate)

					if err != nil {
						return err
					}

					return nil
				default:
					return nil
				}
			}),
			Transform: survey.Transformer(func(ans interface{}) (newAns interface{}) {
				switch dueDate := ans.(type) {
				case string:
					parsedDueDate, err := parseDueDate(dueDate)

					if err != nil {
						return nil
					}

					return parsedDueDate
				default:
					return nil
				}
			}),
		},
		{
			Name: "confirm",
			Prompt: &survey.Confirm{
				Message: "Do you want create the Milestone?",
			},
		},
	}

	createCommand := &cobra.Command{
		Use:   "create",
		Short: "Create Github Milestones",
		Run: func(command *cobra.Command, args []string) {
			description := getDescription(command)
			dueDate, err := getDueDate(command)
			title := getTitle(command)
			verbose, _ := command.Flags().GetBool("verbose")

			if err != nil {
				fmt.Println(err.Error())

				return
			}

			if title == nil {
				answers := struct {
					Confirm     bool
					Description string
					DueDate     *time.Time
					Title       string
				}{}

				err := survey.Ask(questions, &answers)

				if err != nil {
					fmt.Println(err.Error())

					return
				}

				description = &answers.Description
				title = &answers.Title
				dueDate = answers.DueDate
			}

			create.NewCreateMilestone(create.CreateMilestoneConfig{
				Description: description,
				DueDate:     dueDate,
				Title:       title,
				Verbose:     verbose,
			}).Execute()
		},
	}

	createCommand.Flags().String("description", "", "Set the description")
	createCommand.Flags().String("due-date", "", fmt.Sprintf("Set the due date [%s]", constants.DateFormat))
	createCommand.Flags().String("title", "", "Set the title")

	return createCommand
}

func init() {
	rootCommand.AddCommand(NewCreateCommand())
}
