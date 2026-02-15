package create

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/google/go-github/v68/github"
)

func handleResponseError(response *github.Response) error {
	if err := response.Body.Close(); err != nil {
		return err
	}

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
				return errors.New("the requested milestone does not exist")
			case "missing_field":
				return fmt.Errorf("the required field %q has not been set", e.Field)
			case "invalid":
				return fmt.Errorf("the content set in the field %q is invalid", e.Field)
			case "already_exists":
				return errors.New("the milestone with the same title already exists")
			case "unprocessable":
				return errors.New("the data sent was unprocessable")
			case "custom":
				return errors.New(e.Message)
			default:
				return errors.New("unknown error")
			}
		}
	}

	return errors.New(errorResponse.Message)
}
