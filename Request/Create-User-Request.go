package Request
import (
	"fUber/types"
	"errors"
	"net/http"
)

type CreateUserRequest struct {
	*types.User
}

func (c *CreateUserRequest) Bind(r *http.Request) error {
	if c.U_name == "" {
		return errors.New("Username is either empty or invalid")
	}

	if c.U_pass == ""{
		return errors.New("Password is either empty or invalid")
	}

	if c.Email_id == ""{
		return errors.New("Email id is either empty or invalid")
	}

	return nil
}