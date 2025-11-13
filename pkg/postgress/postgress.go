package postgress

import (
	"database/sql"
	"fmt"

	"github.com/celio001/prodify/config"
	_ "github.com/lib/pq"
)

type DbConfig struct {
	User     string
	Host     string
	Port     string
	Password string
	DbName   string
}

func NewInstance() (*sql.DB, error) {
	dbConfig := DbConfig{
		User:     config.GetString("USER_POST"),
		Host:     config.GetString("HOST_POST"),
		Port:     config.GetString("PORT_POST"),
		Password: config.GetString("PASSWORD_POST"),
		DbName:   config.GetString("DB_NAME_POST"),
	}

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DbName)

	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
