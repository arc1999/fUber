package Response
import (
"fUber/Request"
"net/http"
)

type ListBookingResponse struct {
	*Request.CreateBookingRequest

}


func (ListBookingResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func ListBooking(b1 *Request.CreateBookingRequest) *ListBookingResponse {
	resp := &ListBookingResponse{CreateBookingRequest: b1}
	return resp
}