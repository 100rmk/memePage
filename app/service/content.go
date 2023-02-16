package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"memePage/app/conf"
	"memePage/app/database/postgresDB"
	"net/http"
	"os"
	"time"
)

const (
	Today   PostgresPeriod = "today"
	Week                   = "week"
	Month                  = "month"
	AllTime                = "all"
)

type (
	FilePathResp struct {
		Ok     bool `json:"ok"`
		Result struct {
			FileId       string `json:"file_id"`
			FileUniqueId string `json:"file_unique_id"`
			FileSize     int    `json:"file_size"`
			FilePath     string `json:"file_path"`
		} `json:"result"`
	}

	FailFilePathResp struct {
		Ok          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}

	Post struct {
		LikesCount    int
		DislikesCount int
		FileId        string
		TgPath        string
	}

	ResultPost struct {
		Likes       int    `json:"likes"`
		Dislikes    int    `json:"dislikes"`
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
	}

	PostgresPeriod string
)

func CheckFiles(posts []postgresDB.PgPost) {
	for _, post := range posts {
		filePath := fmt.Sprintf("%s/%s", conf.AppConf.ContentPath, post.FileID)
		if !IsFileExist(filePath) {
			postTgUrl := make(chan string)
			go GetFilePath(post.FileID, postTgUrl)
			go DownloadTgFile(<-postTgUrl, filePath)
			close(postTgUrl)
		}
	}
}

func GetFilePath(postId string, tgUrl chan string) {
	var client http.Client

	reqUrl := fmt.Sprintf(
		"https://api.telegram.org/bot%s/getFile?file_id=%s",
		conf.Tg.Token,
		postId)

	res, err := client.Get(reqUrl)
	defer res.Body.Close()
	if err != nil {
		log.Printf("could not create request: %s\n", err)
	}

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("could not read response body: %s\n", err)
	}

	if res.StatusCode != http.StatusOK {
		parsedResp := FailFilePathResp{}
		json.Unmarshal(respBody, &parsedResp)
		log.Printf("Failed to get file path %s", parsedResp.Description)
	} else {
		parsedResp := FilePathResp{}
		json.Unmarshal(respBody, &parsedResp)

		tgUrl <- parsedResp.Result.FilePath
	}

}

func DownloadTgFile(fileUrl string, filePath string) {
	var client http.Client
	reqUrl := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", conf.Tg.Token, fileUrl)
	resp, err := client.Get(reqUrl)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("could not create request: %s\n", err)
	}

	out, err := os.Create(filePath)
	defer out.Close()
	if err != nil {
		log.Printf("Failed to create file: %s\n", err)
	}

	_, err = io.Copy(out, resp.Body)
}

func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// Golang считает, что начало недели это воскресенье, я так не считаю
func getWeekday(now time.Weekday) int {
	if int(now) == 0 {
		return 6
	} else {
		return int(now) - 1
	}
}

func (r PostgresPeriod) GetSearchPeriodParams() (int64, time.Time, error) {
	now := time.Now()
	switch r {
	case Today:
		return 3,
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local),
			nil
	case Week:
		return 10,
			time.Date(now.Year(), now.Month(), now.Day()-getWeekday(now.Weekday()), 0, 0, 0, 0, time.Local),
			nil
	case Month:
		return 10,
			time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local),
			nil
	case AllTime:
		return 10,
			time.Time{},
			nil
	default:
		return 0, time.Time{}, errors.New("unsupported period")
	}
}

func GetPosts(posts []postgresDB.PgPost) []ResultPost {
	result := make([]ResultPost, len(posts))

	for i, post := range posts {
		p := ResultPost{
			Likes:       post.LikesCount,
			Dislikes:    post.DislikesCount,
			URL:         fmt.Sprintf("%s/content/%s", conf.AppConf.ServerUrl, post.FileID),
			ContentType: post.ContentType,
		}

		result[i] = p
	}
	return result
}
