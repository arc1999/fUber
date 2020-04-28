package Request
import (
	"fUber/types"
	"errors"
	"net/http"
)

type CreateBookingRequest struct {
	*types.Booking
}

func (c *CreateBookingRequest) Bind(r *http.Request) error {
	if c.From == "" {
		return errors.New("From is either empty or invalid")
	}

	if c.To == ""{
		return errors.New("To is either empty or invalid")
	}
	if c.U_lat == 0.0{
		return errors.New("lat is either empty or invalid")
	}
	if c.U_long == 0.0{
		return errors.New("long is either empty or invalid")
	}
	return nil
}