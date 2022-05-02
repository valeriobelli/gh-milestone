package create

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/google/go-github/v44/github"
)

func handleResponseError(response *github.Response) error {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	var errorResponse github.ErrorResponse

	err = json.Unmarshal(body, &errorResponse)

	if err != nil {
		return err
	}

	if response.StatusCode == 400 {
		return errors.New(errorResponse.Message)
	}

	if response.StatusCode == 422 {
		for _, e := range errorResponse.Errors {
			switch e.Code {
			case "missing":
				return errors.New(fmt.Sprintf("The requested milestone does not esist."))
			case "missing_field":
				return errors.New(fmt.Sprintf("The required field \"%s\" has not been set.", e.Field))
			case "invalid":
				return errors.New(fmt.Sprintf("The content set in the field \"%s\" is invalid.", e.Field))
			case "already_exists":
				return errors.New("The milestone with the same title already exists.")
			case "unprocessable":
				return errors.New("The data sent was unprocessable.")
			case "custom":
				return errors.New(e.Message)
			default:
				return errors.New("Unknown error")
			}
		}
	}

	return errors.New(errorResponse.Message)
}
