package Response
import (
	"fUber/types"
	"net/http"
)

type BookingsResponse struct {
*types.Bookings
}

func ListBookings(b1 *types.Bookings) *BookingsResponse{
return &BookingsResponse{ b1}

}

func (e *BookingsResponse) Render(w http.ResponseWriter, r *http.Request) error {

return nil
}

