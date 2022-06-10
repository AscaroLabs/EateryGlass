package delivery

// Данный пакет отвечает за создание end-point'ов
// приложения и их обработку

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Создание gin.Engin обекта, который отвечает
// за маршрутизацию запростов и ответов
func NewRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/restaurants", MakeHandler("getRestaurants", db))
	router.GET("/tables", MakeHandler("getTables", db))
	router.POST("/reservations", MakeHandler("postReservations", db))
	return router
}
