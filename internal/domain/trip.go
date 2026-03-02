package domain

type Route struct {
	Origin      City
	Destination City
}
type Trip struct {
	ID             int
	Route          Route
	DepartureTime  string
	ArrivalTime    string
	PriceEUR       float64
	AvailableSeats int
}