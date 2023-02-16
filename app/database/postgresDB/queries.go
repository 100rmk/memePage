package postgresDB

import (
	"fmt"
	"log"
	"memePage/app/conf"
	"time"
)

func GetPostsByRating(limit int64, date time.Time) []PgPost {
	query := fmt.Sprintf(`
		SELECT file_id, timestamp, content_type, likes_count, dislikes_count
		FROM %s.rating
		WHERE timestamp >= $1
		ORDER BY likes_count DESC
		LIMIT $2
	`, conf.Postgres.DBName)

	rows, err := DB.Query(query, date, limit)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer rows.Close()

	var res []PgPost
	for rows.Next() {
		var p PgPost
		err = rows.Scan(&p.FileID, &p.Timestamp, &p.ContentType, &p.LikesCount, &p.DislikesCount)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		res = append(res, p)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
		panic(err)
	}

	return res
}
