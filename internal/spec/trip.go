package spec

type Route struct {
	OriginCity      City `json:"originCity"`
	DestinationCity City `json:"destinationCity"`
}

type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Trip struct {
	ID             int     `json:"id"`
	Route          Route  `json:"route"`
	Company        Company `json:"company"`

	DepartureDate  string `json:"departure_date"`
	DepartureTime  string `json:"departure_time"`
	ArrivalTime    string `json:"arrival_time"`

	Price          string `json:"price"`
	SeatCapacity   int    `json:"seat_capacity"`
	AvailableSeats int    `json:"available_seats"`
}
