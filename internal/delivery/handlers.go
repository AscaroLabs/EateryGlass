package delivery

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AscaroLabs/EateryGlass/pkg/database/postgresql/storage"
	"github.com/AscaroLabs/EateryGlass/pkg/selection"
	"github.com/AscaroLabs/EateryGlass/pkg/structures"
	"github.com/gin-gonic/gin"
)

func MakeHandler(handler_name string, db *sql.DB) func(c *gin.Context) {
	switch handler_name {
	case "getRestaurants":
		return func(c *gin.Context) {
			rest, err := storage.GetRestaurants(db)
			if err != nil {
				log.Fatal(err)
			}
			c.IndentedJSON(http.StatusOK, rest)
		}
	case "getTables":
		return func(c *gin.Context) {
			log.Printf("\n-----* /tables *-----\n")
			log.Printf("\n--* Query: %v *--\n", c.Request.URL.RawPath)
			volume, ok_volume := c.GetQuery("volume")
			appropriate_time, ok_time := c.GetQuery("time")

			log.Printf("\nvolume(%v)[%v] and time(%v)[%v] getted!\n",
				volume, ok_volume,
				appropriate_time, ok_time)

			if !(ok_volume && ok_time) {
				log.Printf("\n!!!---* ERROR *---!!!\n")
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Not enough parameters!"})
				return
			}
			TablesByRestaurants, err := selection.SelectTables(db, volume, appropriate_time)
			if err != nil {
				log.Fatal(err)
			}
			c.IndentedJSON(http.StatusOK, TablesByRestaurants)
		}
	case "postReservations":
		return func(c *gin.Context) {
			var newReservation structures.RawReservation
			// type RawReservation struct {
			// 	Table_id    string `json:"table_id"`
			// 	Start_time  string `json:"start_time"`
			// 	Reserved_by RawClient `json:"reserved_by"`
			// }
			if err := c.BindJSON(&newReservation); err != nil {
				log.Printf("Smth wrong with parse JSON %v", err)
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "wrong JSON!"})
				return
			}
			res, err := storage.PostReservations(db, newReservation)
			if err != nil {
				log.Fatal(err)
			}
			c.IndentedJSON(http.StatusCreated, res)
		}
	default:
		return nil
	}
}
