package structures

// CREATE TABLE clients (
// 		id integer PRIMARY KEY,
// 		name char(50) NOT NULL,
// 		phone_number char(50) NOT NULL);

type Client struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone_number"`
}

type RawClient struct {
	Name  string `json:"name"`
	Phone string `json:"phone_number"`
}
