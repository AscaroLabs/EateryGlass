package structures

// Пакет structures содержит используемые в
// приложении структуры

// Данные структуры отвечает за клиентов,
// которые бронируют столики

// CREATE TABLE clients (
// 	id SERIAL PRIMARY KEY,				уникальный идентификатор
//
// 	name TEXT NOT NULL,					имя
//
// 	phone_number TEXT NOT NULL, 		номер телефона
//
// 	UNIQUE (name, phone_number)			констрейнт, не позволяющий
// );									добавлять одинаковых клиентов

// Структура, которая возвращается из
// таблицы clients
type Client struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone_number"`
}

// Структура, которая передается как поле
// при запросе на бронирование столика
type RawClient struct {
	Name  string `json:"name"`
	Phone string `json:"phone_number"`
}
