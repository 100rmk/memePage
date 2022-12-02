package routes

import "github.com/gin-gonic/gin"

func Memes(router *gin.RouterGroup) {
	router.GET("/memes", GetMemes)
}
