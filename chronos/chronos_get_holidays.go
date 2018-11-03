package chronos

import (
	"encoding/json"
	"strconv"

	"github.com/go-resty/resty"
)

func (h *Chronos) GetHolidays(year int) error {
	query := map[string]string{
		"country": "BR",
		"year":    strconv.Itoa(year),
	}

	resp, err := resty.R().SetQueryParams(query).Get(HOLIDAY_API_URL)
	if err != nil {
		return ErrUnableToSendGetHolidaysRequest
	}

	var res holidaysResponse
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return ErrUnableToUnmarshalGetHolidaysResponse
	}

	h.holidays = res.Holidays

	return nil
}

type holidaysResponse struct {
	status   int
	Holidays Holidays `json:"holidays"`
}
