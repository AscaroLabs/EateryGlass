package structures

type Table struct {
	ID            string `json:"id"`
	Restaurant_id string `json:"restaurant_id"`
	Capacity      int    `json:"capacity"`
}
