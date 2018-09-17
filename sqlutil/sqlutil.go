package sqlutil

import (
	"database/sql"
	_ "github.com/lib/pq" //postgres driver
)

//Conn ...
func Conn(connStr string) *sql.DB {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("Connecting Error: " + err.Error())

	}

	err = conn.Ping()
	if err != nil {
		panic("Ping Error:" + err.Error())
	}

	return conn
}
