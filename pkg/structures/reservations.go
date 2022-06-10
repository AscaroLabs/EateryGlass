package structures

type Reservation struct {
	ID          string `json:"id"`
	Table_id    string `json:"table_id"`
	Start_time  string `json:"start_time"`
	End_time    string `json:"end_time"`
	Reserved_by string `json:"reserved_by"`
}

type RawReservation struct {
	Table_id    string    `json:"table_id"`
	Start_time  string    `json:"start_time"`
	Reserved_by RawClient `json:"reserved_by"`
}
