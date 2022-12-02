package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"memePage/app/database/mongo"
	"memePage/app/service"
	"net/http"
	"strconv"
	"sync"
)

func GetMemes(c *gin.Context) {
	cCount := c.Query("count")
	count, err := strconv.Atoi(cCount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count must be int"})
	}
	if count > 100 || count < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count must be less than 100 or greater than 1"})
	}

	postsChan := make(chan []bson.M)
	go mongo.GetPostsByRating(count, postsChan)
	dbPosts := <-postsChan
	close(postsChan)
	posts := make([]service.Post, len(dbPosts))

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go service.GetFilePaths(wg, dbPosts, &posts)
	wg.Wait()

	for _, elem := range posts {
		go service.DownloadTgFile(elem)
	}

	log.Println(posts)

	c.JSON(http.StatusOK, gin.H{"result": "true"})
}
