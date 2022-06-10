package structures

// Пакет structures содержит используемые в
// приложении структуры

// Данная структура отвечает за столики,
// которые бронируют

// CREATE TABLE tables (
// 	id SERIAL PRIMARY KEY,								уникальный идентификатор
//
// 	restaurant_id integer REFERENCES restaurants (id),	ресторан в котором
//														находится столик
// 	capacity integer NOT NULL							вместимость столика
// );

// Структура, которая возвращается из
// таблицы tables
type Table struct {
	ID            string `json:"id"`
	Restaurant_id string `json:"restaurant_id"`
	Capacity      int    `json:"capacity"`
}
