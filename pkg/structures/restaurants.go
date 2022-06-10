package structures

type Restaurant struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Avg_time  string `json:"avg_time"`
	Avg_price int    `json:"avg_price"`
}
