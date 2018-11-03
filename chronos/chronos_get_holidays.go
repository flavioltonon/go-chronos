package chronos

import (
	"encoding/json"
	"strconv"

	"github.com/go-resty/resty"
)

type Holidays map[string][]Holiday

type Holiday struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Date    string `json:"date"`
}

type holidaysResponse struct {
	status   int
	Holidays Holidays `json:"holidays"`
}

func (h *Chronos) GetHolidays(year int) (Holidays, error) {
	query := map[string]string{
		"country": "BR",
		"year":    strconv.Itoa(year),
	}

	resp, err := resty.R().SetQueryParams(query).Get(HOLIDAY_API_URL)
	if err != nil {
		return nil, ErrUnableToSendGetHolidaysRequest
	}

	var res holidaysResponse
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, ErrUnableToUnmarshalGetHolidaysResponse
	}

	return res.Holidays, nil
}
