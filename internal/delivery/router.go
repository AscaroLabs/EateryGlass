package delivery

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.GET("/restaurants", MakeHandler("getRestaurants", db))
	router.GET("/tables", MakeHandler("getTables", db))
	router.POST("/reservations", MakeHandler("postReservations", db))
	return router
}
