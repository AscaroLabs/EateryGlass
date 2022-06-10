package structures

// Пакет structures содержит используемые в
// приложении структуры

// Данные структуры отвечает за брони.

// CREATE TABLE reservations (
// 	id SERIAL PRIMARY KEY,						уникальный идентификатор
//
// 	table_id integer REFERENCES tables (id),	столик, за которым закреплениа бронь
//
// 	start_time timestamp NOT NULL,				начало действия брони
//
// 	end_time timestamp NOT NULL,				конец действия брони (начало + 2 часа)
//
// 	reserved_by int REFERENCES clients (id),	id клиента за которым закреплена бронь
//
// 	UNIQUE (table_id, start_time)				констрейнт, не позволяющий:
//												во-первых, создавать одинаковые брони
//												во-вторых, двум клиентам бронировать
// );											один и тот же столик на одно и то же время

// Структура, которая возвращается из
// таблицы reservations
type Reservation struct {
	ID          string `json:"id"`
	Table_id    string `json:"table_id"`
	Start_time  string `json:"start_time"`
	End_time    string `json:"end_time"`
	Reserved_by string `json:"reserved_by"`
}

// Структура, которая передается при
// запросе на бронирование столика
type RawReservation struct {
	Table_id    string    `json:"table_id"`
	Start_time  string    `json:"start_time"`
	Reserved_by RawClient `json:"reserved_by"`
}
