package postgresDB

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"memePage/app/conf"
	"time"
)

var DB *sql.DB

type PgPost struct {
	FileID        string    `db:"file_id"`
	Timestamp     time.Time `db:"timestamp"`
	ContentType   string    `db:"content_type"`
	LikesCount    int       `db:"likes_count"`
	DislikesCount int       `db:"dislikes_count"`
}

func ConnectDb(uri string) {
	if uri == "" {
		log.Fatal("You must set your 'POSTGRES_URI' environmental variable.")
	}
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}

	DB = db
}

func CreateViews() {
	pipeline := fmt.Sprintf(`
			SELECT file_id,
			   timestamp,
			   content_type,
			   coalesce(array_length(likes, 1), 0) AS likes_count,
			   coalesce(array_length(dislikes, 1), 0) AS dislikes_count
		FROM %s.posts;
		`, conf.Postgres.DBName)

	_, err := DB.Exec(fmt.Sprintf("CREATE OR REPLACE VIEW %s.rating AS %s", conf.Postgres.DBName, pipeline))
	if err != nil {
		panic(err)
	}

}

func DisconnectDb() {
	func() {
		if err := DB.Close(); err != nil {
			panic(err)
		}
	}()
}
