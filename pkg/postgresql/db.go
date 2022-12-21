package postgresql

import (
	"database/sql"
	"os"
	logger "project/log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var Conn *sql.DB


//var dsn string = "host=localhost port=5432 user=postgres password=nikon1337 dbname=users sslmode=disable timezone=UTC connect_timeout=5"
var dsn string = os.Getenv("DSN")
var count int64


func DatabaseInit() *sql.DB {
	for {
		var err error
		Conn, err = openDB(dsn)
		
		if err != nil {
			logger.Info("Potgres Not Yet Ready...")
			count++
		} else {
			logger.Info("Connected To Postgres!")
			queryCreateTable := `CREATE TABLE IF NOT EXISTS public.User (
				id serial primary key,
				name varchar(50) not null,
				email varchar(250) unique not null,
				password varchar(50) not null
			 );`
			stmt, err := Conn.Prepare(queryCreateTable)
			if err != nil {
				logger.Error("error when trying to prepare create table statement",err)
				return nil
			}
			defer stmt.Close()

			stmt.Exec()

			return Conn
		}

		if count > 10 {
			logger.Error("error connecting to DB",err)
			return nil
		}

		logger.Info("Backing off for two seconds..")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx",dsn)
	if err != nil {
		return nil,err
	}

	err = db.Ping()
	if err != nil {
		return nil,err
	}

	return db, nil
}