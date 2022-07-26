package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conectar() (*sql.DB, error) {

	stringConexao := "golang:golang@/devbook?charset=utf8"

	db, erro := sql.Open("mysql", stringConexao)

	if erro != nil {
		return nil, erro
	}

	if erro := db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil

}
