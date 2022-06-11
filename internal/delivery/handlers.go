package delivery

// Пакет delivery отвечает за создание end-point'ов
// приложения и их обработку

import (
	"database/sql"
	"net/http"

	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/storage"
	"github.com/AscaroLabs/EateryGlass/pkg/selection"
	"github.com/AscaroLabs/EateryGlass/pkg/structures"
	"github.com/gin-gonic/gin"
)

// MakeHandler создает обработчики для end-point'ов
// Нужна для того что бы пробросить указатель на
// БД внуть обработчиков и обработчики удовленворяли
// сигнатуре gin.HandlerFunc
func MakeHandler(handlerName string, db *sql.DB) func(c *gin.Context) {
	switch handlerName {
	// Обработчик для "/restaurants"
	case "getRestaurants":
		return func(c *gin.Context) {

			// Забираем рестораны из БД
			restaurants, err := storage.GetRestaurants(db)
			if err != nil {
				// Отправляем сообщение об ошибке в качестве ответа
				c.IndentedJSON(http.StatusBadRequest,
					gin.H{"error": err.Error()})
				return
			}

			// Возвращаем данные о ресторанах
			c.IndentedJSON(http.StatusOK, restaurants)
		}
	// Обработчик для "/tables"
	case "getTables":
		return func(c *gin.Context) {

			// Забираем данные из GET запроса
			// e.g. tables?volume=7&time=2022-01-02T16:06:06Z
			volume, okVolume := c.GetQuery("volume")
			appropriateTime, okTime := c.GetQuery("time")

			// Возвращаем ошибку, если необходимые параметры не были переданы
			if !(okVolume && okTime) {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Not enough parameters!"})
				return
			}

			// Получаем доступные столики из БД
			tables, err := selection.SelectTables(db, volume, appropriateTime)
			if err != nil {
				// Отправляем сообщение об ошибке в качестве ответа
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Возвращаем данные о возможных вариантах
			c.IndentedJSON(http.StatusOK, tables)
		}

	// Обработчик для "/reservations"
	case "postReservations":
		return func(c *gin.Context) {

			// Создаем переменную для тела POST запроса
			var rawReservations []structures.RawReservation

			// Пытаемся спарсить JSON из тела ответа
			// Если не получилось, то отправляем ошибку
			if err := c.BindJSON(&rawReservations); err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Пытаемся положить в БД новые брони
			reservations, err := storage.PostReservations(db, rawReservations)
			if err != nil {
				// Отправляем сообщение об ошибке в качестве ответа
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.IndentedJSON(http.StatusCreated, reservations)
		}
	default:
		return nil
	}
}
