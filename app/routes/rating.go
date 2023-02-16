package routes

import (
	"github.com/gin-gonic/gin"
	"memePage/app/database/postgresDB"
	"memePage/app/service"
	"net/http"
)

func GetMemes(c *gin.Context) {
	cPeriod := c.Query("period")
	count, date, err := service.PostgresPeriod(cPeriod).GetSearchPeriodParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	posts := postgresDB.GetPostsByRating(count, date)
	if posts == nil {
		c.JSON(http.StatusOK, gin.H{"result": "Nothing to show you right now"})
		return
	}
	service.CheckFiles(posts)
	c.JSON(http.StatusOK, gin.H{"result": service.GetPosts(posts)})
}
