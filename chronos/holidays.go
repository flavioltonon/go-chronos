package chronos

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
