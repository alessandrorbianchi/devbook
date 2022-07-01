package banco

import (
	"database/sql"
	"send/src/config"

	_ "github.com/go-sql-driver/mysql"
)

func Conectar() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.StringConexaoBanco)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
