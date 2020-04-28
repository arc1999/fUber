package Response

import (
	"fUber/Request"
	"net/http"
)

type ListUserResponse struct {
	*Request.CreateUserRequest

}


func (ListUserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListUser(u1 *Request.CreateUserRequest) *ListUserResponse {
	resp := &ListUserResponse{CreateUserRequest: u1}
	return resp
}