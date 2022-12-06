package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memePage/app/conf"
	"memePage/app/database/mongoDB"
	"memePage/app/routes"
)

func Run() {
	mongoDB.ConnectDb(conf.Mongo.Url)
	mongoDB.CreateViews(conf.Mongo.DBName)
	mongoDB.InitCollections(conf.Mongo.DBName)
	defer mongoDB.DisconnectDb()

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
