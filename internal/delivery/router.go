package delivery

// Пакет delivery отвечает за создание end-point'ов
// приложения и их обработку

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// NewRouter создает gin.Engin обект, который отвечает
// за маршрутизацию запростов и ответов
func NewRouter(db *sql.DB) *gin.Engine {

	// Стандарнтый gin роутер
	router := gin.Default()

	// Запрос для получения данных о ресторанах
	// e.g. 172.17.0.3:8000/restaurants
	router.GET("/restaurants", MakeHandler("getRestaurants", db))

	// Запрос для получения информации о свободных столиках
	// e.g. 172.17.0.3:8000/tables?volume=7&time=2022-01-02T16:06:06Z
	router.GET("/tables", MakeHandler("getTables", db))

	// Запрос для брони столика
	// e.g. 172.17.0.3:8000/reservations
	// 		{
	// 			"table_id": "13",
	// 			"start_time": "2022-01-02T16:06:06Z",
	// 			"reserved_by": {
	// 				"name": "Maxim",
	// 				"phone_number": "88005553535"
	// 			}
	// 		}
	router.POST("/reservations", MakeHandler("postReservations", db))

	return router
}
