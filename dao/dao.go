package dao

import (
	"database/sql"
)

type Database struct {
	Conn *sql.DB
}
