package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memePage/app/conf"
	"memePage/app/database/postgresDB"
	"memePage/app/routes"
)

func Run() {
	postgresDB.ConnectDb(conf.Postgres.Url)
	postgresDB.CreateViews()
	defer postgresDB.DisconnectDb()

	router := initRoutes()
	err := router.Run(fmt.Sprintf("%s:%s", conf.AppConf.Host, conf.AppConf.Port))
	if err != nil {
		return
	}
}

func initRoutes() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api")
	routes.Memes(v1.Group("/content"))
	return router
}
