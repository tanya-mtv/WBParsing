package db

import (
	"database/sql"
	"fmt"
	"log"
	"parsingWB/internal/config"

	_ "github.com/alexbrainman/odbc"
)

func NewMSSQLDB(cfg *config.ConfigMSSQL) *sql.DB {

	connstring := fmt.Sprintf("driver={%s};SERVER=%s,%d;UID=%s;PWD=%s;DATABASE=%s;TrustServerCertificate=yes", cfg.DriverName, cfg.Server, cfg.Port, cfg.User, cfg.Password, cfg.DSN)

	db, err := sql.Open(cfg.Driver, connstring)

	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}

	return db
}
