package structures

// Пакет structures содержит используемые в
// приложении структуры

// Данная структуры отвечает за рестораны,
// столики в которых бронируются

// CREATE TABLE restaurants (
// 	id SERIAL PRIMARY KEY,			уникальный индентификатор

// 	name TEXT NOT NULL,				название ресторана

// 	avg_time interval NOT NULL,		среднее время ожидания

// 	avg_price integer NOT NULL,		средний чек

// 	UNIQUE(name)					констрейнт, не позволяющий
// );								добавлять одинаковые рестораны

// Структура, которая возвращается из
// таблицы restaurants
type Restaurant struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Avg_time  string `json:"avg_time"`
	Avg_price int    `json:"avg_price"`
}
