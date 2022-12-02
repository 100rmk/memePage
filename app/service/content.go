package service

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"memePage/app/conf"
	"net/http"
	"os"
	"sync"
)

type (
	GetFilePath struct {
		Ok     bool `json:"ok"`
		Result struct {
			FileId       string `json:"file_id"`
			FileUniqueId string `json:"file_unique_id"`
			FileSize     int    `json:"file_size"`
			FilePath     string `json:"file_path"`
		} `json:"result"`
	}

	FailGetFilePath struct {
		Ok          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	Post struct {
		LikesCount int
		FileId     string
		TgPath     string
	}
)

func GetFilePaths(wg *sync.WaitGroup, posts []bson.M, resultPosts *[]Post) {
	var client http.Client
	defer wg.Done()

	for i, elem := range posts {
		reqUrl := fmt.Sprintf(
			"https://api.telegram.org/bot%s/getFile?file_id=%s",
			conf.Tg.Token,
			elem["file_id"])

		res, err := client.Get(reqUrl)
		if err != nil {
			log.Printf("could not create request: %s\n", err)
		}

		respBody, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("could not read response body: %s\n", err)
		}

		if res.StatusCode != http.StatusOK {
			parsedResp := FailGetFilePath{}
			json.Unmarshal(respBody, &parsedResp)
			log.Printf("Failed to get file path %s", parsedResp.Description)
		} else {
			parsedResp := GetFilePath{}
			json.Unmarshal(respBody, &parsedResp)
			(*resultPosts)[i] = Post{
				LikesCount: int(elem["count"].(int32)),
				FileId:     parsedResp.Result.FileId,
				TgPath:     parsedResp.Result.FilePath}
		}

		res.Body.Close()
	}
}

func DownloadTgFile(post Post) {
	var client http.Client
	reqUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", conf.Tg.Token, post.TgPath)
	resp, err := client.Get(reqUrl)
	if err != nil {
		log.Printf("could not create request: %s\n", err)
	}
	defer resp.Body.Close()
	// TODO: поменять на нормальный путь "/opt/content/files/%s"
	out, err := os.Create(fmt.Sprintf("./%s", post.FileId))
	if err != nil {
		log.Printf("Failed to create file: %s\n", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
}

// get path https://api.telegram.org/bot1925154606:AAE28-g4bxzBdKc3eIXMvDJr01PqSAqsihU/getFile?file_id=AgACAgIAAxkBAAJA6GNNUaBcB-9SuuBOX2Ax4h5iZJaVAAKyxDEbz-hoSuYGheVT5kH3AQADAgADeQADKgQ
// download by path https://api.telegram.org/file/bot1925154606:AAE28-g4bxzBdKc3eIXMvDJr01PqSAqsihU/photos/photos/file_6120
