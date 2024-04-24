package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zhosyaaa/RoommateTap/internal/config"
	"sync"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func initializeDB(config config.Config) (*sql.DB, error) {
	var errInit error
	dbOnce.Do(func() {
		dbConnString := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
			config.Database.Sslmode,
		)
		fmt.Println(dbConnString)
		db, errInit = sql.Open("postgres", dbConnString)
		if errInit != nil {
			return
		}
		errPing := db.Ping()
		if errPing != nil {
			return
		}
	})
	return db, errInit
}
func GetDBInstance(cfg config.Config) (*sql.DB, error) {
	var errGetDB error
	if db == nil {
		var initErr error
		db, initErr = initializeDB(cfg)
		if initErr != nil {
			errGetDB = initErr
		}
	}
	return db, errGetDB
}
