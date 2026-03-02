package spec

type SearchQuery struct {
	FromCityID  int    `json:"fromCity"`
	ToCityID    int    `json:"toCity"`
	Date        string `json:"date"`
}