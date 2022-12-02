package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memePage/app/conf"
	"memePage/app/database/mongo"
	"memePage/app/routes"
)

func Run() {
	mongo.ConnectDb(conf.Mongo.Url)
	mongo.InitCollections(conf.Mongo.DBName)
	defer mongo.DisconnectDb()

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
