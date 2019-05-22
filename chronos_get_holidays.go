package chronos

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty"
)

var holidays Holidays

type ChronosGetHolidaysRequest struct {
	Country string
	Year    string
}

type ChronosGetHolidaysResponse struct {
	Holidays Holidays `json:"holidays"`
}

func (h Chronos) GetHolidays() (ChronosGetHolidaysResponse, error) {
	req := h.request.(ChronosGetHolidaysRequest)
	var response ChronosGetHolidaysResponse

	r, err := resty.R().SetQueryParams(map[string]string{
		"country": req.Country,
		"year":    req.Year,
	}).Get(os.Getenv("HOLIDAY_API_URL"))
	if err != nil {
		return ChronosGetHolidaysResponse{}, ErrUnableToSendGetHolidaysRequest
	}
	var resp holidaysResponse

	err = json.Unmarshal(r.Body(), &resp)
	if err != nil {
		return ChronosGetHolidaysResponse{}, ErrUnableToUnmarshalGetHolidaysResponse
	}

	response.Holidays = resp.Holidays

	return response, nil
}
