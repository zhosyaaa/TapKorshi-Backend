package domain

type Location struct {
	ID              uint    `json:"id"`
	Country         string  `json:"country"`
	City            string  `json:"city"`
	Street          string  `json:"street"`
	HouseNumber     string  `json:"houseNumber"`
	Floor           int     `json:"floor"`
	ApartmentNumber string  `json:"apartmentNumber"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}
